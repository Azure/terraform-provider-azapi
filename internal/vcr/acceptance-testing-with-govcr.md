# Acceptance testing with go-vcr

## Summary

[`go-vcr`](https://github.com/dnaeon/go-vcr) records the HTTP traffic of an acceptance test once and replays it later, so the test suite can run with **no network access and no Azure credentials**. Replaying is deterministic, which makes VCR tests fast and reliable in CI. This guide describes how the plumbing works in `terraform-provider-azapi`, how to run it, and where to look if you need to change something.

The harness lives in [`internal/vcr`](.) and is intentionally structured to match the sibling [`terraform-provider-azurerm`](https://github.com/hashicorp/terraform-provider-azurerm) implementation so the two providers share a mental model.

## How to run tests

The framework is controlled by the `TC_TEST_VIA_VCR` environment variable. To prevent accidentally mixing real and mock data, intent must be explicit.

| `TC_TEST_VIA_VCR`     | Mode        | Behavior |
| --------------------- | ----------- | -------- |
| `record`              | record      | Talks to the real Azure APIs using your environment (`ARM_SUBSCRIPTION_ID`, etc.) and writes a sanitized cassette. |
| `replay` / `true`     | replay      | Skips the network entirely and serves every request from the cassette using a fake credential. |
| unset / anything else | passthrough | VCR is disabled; `VcrResourceTest` tests are skipped and live tests are unaffected. |

Convenience make targets wrap the two active modes:

```sh
# record a single test live against Azure (needs credentials + quota)
make vcrrecord TESTARGS='-run TestVcrAccGenericResource_basic'

# replay every VCR test from its cassette (no network, no credentials)
make vcrreplay
```

## Recording new tests: `TC_TEST_VIA_VCR=record`

The framework talks to the real Azure APIs using your environment variables. It records exactly what happens to a `.yaml` cassette under `internal/services/testdata/cassettes/<TestName>.yaml`, scrubbing your real identifiers and secrets **after** the HTTP traffic completes and immediately before the file is written to disk.

## Replaying tests: `TC_TEST_VIA_VCR=replay`

The acceptance framework skips the network entirely. At start-up it globally overrides the authentication environment variables (`ARM_SUBSCRIPTION_ID`, `ARM_TENANT_ID`, `ARM_CLIENT_ID`, `ARM_CLIENT_SECRET`) to their canonical placeholders and disables the live login methods (`ARM_USE_CLI`/`ARM_USE_MSI`/`ARM_USE_OIDC`/`ARM_USE_AKS_WORKLOAD_IDENTITY`). The provider then natively builds its request URLs and applied state from those placeholders, which is what makes them line up with the redacted cassettes — including during `ImportStep`. See `internal/acceptance/vcr_init.go`.

Because azapi authenticates through the `azure-sdk-for-go` `azcore` pipeline (which would otherwise fetch a real token), replay swaps in a **fake credential** that mints a syntactically valid JWT with canonical claims without contacting Azure. This is wired in `internal/clients/client.go` and implemented by `FakeCredential()` in `recorder.go`.

## Passthrough (default fallback)

If the variable is unset or unrecognised the framework falls back to passthrough: VCR is ignored and traffic goes straight to Azure without recording or replaying. `VcrResourceTest` skips itself in this mode.

## What actually happens: redaction, audit & the matcher

The core lives in [`internal/vcr/recorder.go`](recorder.go), with redaction and the safety audit in [`sanitize_interaction.go`](sanitize_interaction.go) and [`audit_interaction.go`](audit_interaction.go).

1. **Per-test recorders.** `GetRecorder(testName, subscriptionID)` builds (and caches, keyed by test name) one recorder per test. `RegisterTestT`/`CurrentTestName` (in `test_context.go`) map the running goroutine to its test so helpers such as `BuildTestClient()` can find the right recorder. `StopRecorder(testName)` flushes and releases it. Because recorders are per test, VCR tests run in parallel.

2. **Canonical placeholders.** Real subscription and tenant IDs are rewritten to fixed placeholders (`00000000-…`, `11111111-…`). To redact a new sensitive identifier, register it with `AddReplacement`/`AddSecret` in `GetRecorder`.

3. **Pre-match subscription scrubbing.** go-vcr's default matcher compares the full request against the cassette. During `ImportStep` (or other edge cases) Terraform may send a request carrying a real subscription ID that would never match an already-redacted cassette. `subscriptionMatcher` solves this by normalising subscription IDs on **both** the incoming request and the recorded request before comparing method + URL.

4. **The BeforeSave hook.** Scrubbing happens only once the test finishes, via go-vcr's `BeforeSaveHook` (`(*sanitizer).BeforeSave`). It rewrites all known secret/identifier values across request and response URLs, bodies, and headers, then writes the clean `.yaml`. Doing this offline at save-time keeps `go-azure-sdk` long-running-operation polling intact during recording.

5. **Fail-closed audit.** Because redaction is a denylist, recording also runs an audit (`audit_interaction.go`) as a safety net: after redaction every interaction is re-scanned for high-signal secret shapes (JWTs, PEM private keys, SAS `sig=` tokens, connection-string keys/passwords, secret JSON fields, high-entropy key blobs). If anything survives, the recording **fails and no cassette is written**, so an unrecognised secret cannot be committed by accident.

6. **Deterministic "random" data.** Replay needs predictable values so requests match. `RandInt(name)`/`RandString(name, n)` derive stable values from the test name via `fnv`, and VCR tests use fixed, name-derived locations (`LocationPrimary`/`Secondary`/`Ternary`). Changing these formulas invalidates existing cassettes.

## Registering runtime secrets

The BeforeSave redaction is a denylist, so a value your test reads back from Azure at runtime (e.g. a storage-account key) must be registered so it is scrubbed:

```go
key := /* value returned by Azure */
vcr.RegisterSecret(key) // no-op unless VCR is recording
```

If a legitimate recording trips the audit, either add a redaction rule in `sanitize_interaction.go` or register the value as above.

## Writing a VCR test

Use `VcrResourceTest` instead of `ResourceTest`; it skips itself when VCR is disabled. See `TestVcrAccGenericResource_basic` in `internal/services/azapi_resource_test.go` for an example. Always review a freshly recorded cassette before committing it — the sanitizer and audit run automatically, but a quick human check is cheap insurance. Cassettes should be version controlled.
