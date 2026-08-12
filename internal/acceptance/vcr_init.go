package acceptance

import (
	"os"

	"github.com/Azure/terraform-provider-azapi/internal/vcr"
)

// During replay every request is served from a cassette using a fake
// credential, so we pin the canonical, non-secret identifiers and disable the
// live authentication methods before any test runs. This mirrors the sibling
// terraform-provider-azurerm behaviour of forcing placeholder subscriptions so
// that resource IDs (and ImportStep) line up with the sanitized cassette, while
// additionally supplying the fake service-principal values azapi's azcore
// pipeline expects to construct a credential.
func init() {
	if !vcr.IsReplaying() {
		return
	}

	overrides := map[string]string{
		"TF_ACC":                        "1",
		"ARM_CLIENT_ID":                 vcr.CanonicalClientID,
		"ARM_CLIENT_SECRET":             vcr.CanonicalClientSecret,
		"ARM_TENANT_ID":                 vcr.CanonicalTenantID,
		"ARM_SUBSCRIPTION_ID":           vcr.CanonicalSubscriptionID,
		"ARM_USE_CLI":                   "false",
		"ARM_USE_MSI":                   "false",
		"ARM_USE_OIDC":                  "false",
		"ARM_USE_AKS_WORKLOAD_IDENTITY": "false",
	}
	for k, v := range overrides {
		_ = os.Setenv(k, v)
	}
}
