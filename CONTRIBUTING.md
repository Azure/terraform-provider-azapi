## Developer Requirements

* [Terraform](https://www.terraform.io/downloads.html) see [.terraform-version](.terraform-version) for the required version
* [Go](https://golang.org/doc/install) see [go.mod](go.mod) for the required version

### On Windows

If you're on Windows you'll also need:

* [Git Bash for Windows](https://git-scm.com/download/win)
* [Make for Windows](http://gnuwin32.sourceforge.net/packages/make.htm)

For *GNU32 Make*, make sure its bin path is added to PATH environment variable.*

For *Git Bash for Windows*, at the step of "Adjusting your PATH environment", please choose "Use Git and optional Unix tools from Windows Command Prompt".*

Or install via [Chocolatey](https://chocolatey.org/install) (`Git Bash for Windows` must be installed per steps above)

```powershell
choco install make golang terraform -y
refreshenv
```

You must run `Developing the Provider` commands in `bash` because `sh` scrips are invoked as part of these.

## Developing the Provider

If you wish to work on the provider, you'll first need [Go](http://www.golang.org) installed on your machine. You'll also need to correctly setup a [GOPATH](http://golang.org/doc/code.html#GOPATH), as well as adding `$GOPATH/bin` to your `$PATH`.

First clone the repository to: `$GOPATH/src/github.com/Azure/terraform-provider-azapi`

```sh
mkdir -p $GOPATH/src/github.com/Azure; cd $GOPATH/src/github.com/Azure
git clone git@github.com:Azure/terraform-provider-azapi
cd $GOPATH/src/github.com/Azure/terraform-provider-azapi
```

Once inside the provider directory, you can run `make tools` to install the dependent tooling required to compile the provider.

At this point you can compile the provider by running `make build`, which will build the provider and put the provider binary in the `$GOPATH/bin` directory.

```sh
$ make build
...
$ $GOPATH/bin/terraform-provider-azapi
...
```

You can also cross-compile if necessary:

```sh
GOOS=windows GOARCH=amd64 make build
```

In order to run the `Unit Tests` for the provider, you can run:

```sh
make test
```

The majority of tests in the provider are `Acceptance Tests` - which provisions real resources in Azure. It's possible to run the entire acceptance test suite by running `make testacc` - however it's likely you'll want to run a subset, which you can do using a prefix, by running:

```sh
make acctests TESTARGS='-run=<nameOfTheTest>' TESTTIMEOUT='60m'
```

* `<nameOfTheTest>` should be self-explanatory as it is the name of the test you want to run. An example could be `TestAccGenericResource_basic`. Since `-run` can be used with regular expressions you can use it to specify multiple tests like in `TestAccGenericResource_` to run all tests that match that expression

The following Environment Variables must be set in your shell prior to running acceptance tests:

* `ARM_CLIENT_ID`
* `ARM_CLIENT_SECRET`
* `ARM_READER_CLIENT_ID`
* `ARM_READER_CLIENT_SECRET`
* `ARM_SUBSCRIPTION_ID`
* `ARM_TENANT_ID`
* `ARM_ENVIRONMENT`: set the value to `public` for Azure Public Cloud
* `ARM_TEST_LOCATION`
* `ARM_TEST_LOCATION_ALT`
* `ARM_TEST_LOCATION_ALT2`

To setup `ARM_CLIENT_ID` and `ARM_CLIENT_SECRET`, with Azure CLI logged in to your tenant and subscription, run the following command to setup a service principal with `Contributor` role:

```sh
az ad sp create-for-rbac \
  --name "myname-azapi-acctest" \
  --role "Contributor" \
  --scopes "/subscriptions/00000000-0000-0000-0000-000000000000"
```

The resulting json will contain the following output:

```json
{
  "appId": "00000000-0000-0000-0000-000000000000",
  "displayName": "myname-azapi-acctest",
  "password": "<redacted>",
  "tenant": "00000000-0000-0000-0000-000000000000"
}
```

Use the `appId` for `ARM_CLIENT_ID` and the `password` for `ARM_CLIENT_SECRET`.

The `ARM_READER_CLIENT_ID` and `ARM_READER_CLIENT_SECRET` follow the same process but with the `Reader` role instead of `Contributor`.

**Note:** Acceptance tests create real resources in Azure which often cost money to run.

## Generating Documentation

We use [tfplugindocs](https://github.com/hashicorp/terraform-plugin-docs) to automatically generate documentation for the provider.
Please ensure that the `MarkdownDescription` field is set in the schema for each resource and data source.

To generate the documentation run either:

```sh
make docs
```

or...

```sh
go generate ./...
```

### Templates

Each resource is documented using a template. The template is located in the `templates` directory. The template is a markdown file with placeholders that are replaced with the actual values from the schema. There is a general template for all resources/data sources, and an optional specific template for each resource/data source where customization is required.

### Guides

Guides should be stored in the `templates/guides` directory. They will be inclided in the documentation and copied to the `docs` directory by the `tfplugindocs` tool.

### Examples

The `examples/resources` and `examples/data-sources` directory contains examples for each resource and data source. The examples are used to generate the documentation for each resource and data source. The examples are written in HCL and must be called `resource.tf` or `data-source.tf`. These are then embedded into the documentation and are used to generate the `Example` section.

## Raising PR

1. Provide a clear summary of the problem in the PR description.
1. Keep scope small, for example: a single bug fix, behaviour update of a single property. Split multiple problems into separate PRs using [stacked PR feature](https://docs.github.com/en/pull-requests/how-tos/stacked-pull-requests) if necessary. Do not raise a PR that changes multiple resources / data sources, or updates the behaviour of multiple unrelated properties.
1. After feedbacks are addressed, explicitly request a review again to indicate readiness.

---

## Developer: Using the locally compiled Azure Provider binary

When using Terraform 0.14 and later, after successfully compiling the Azure Provider, you must [instruct Terraform to use your locally compiled provider binary](https://www.terraform.io/docs/commands/cli-config.html#development-overrides-for-provider-developers) instead of the official binary from the Terraform Registry.

For example, add the following to `~/.terraformrc` for a provider binary located in `/home/developer/go/bin`:

```
provider_installation {

  # Use /home/developer/go/bin as an overridden package directory
  # for the Azure/azapi provider. This disables the version and checksum
  # verifications for this provider and forces Terraform to look for the
  # azapi provider plugin in the given directory.
  dev_overrides {
    "Azure/azapi" = "/home/developer/go/bin"
  }

  # For all other providers, install them directly from their origin provider
  # registries as normal. If you omit this, Terraform will _only_ use
  # the dev_overrides block, and so no other providers will be available.
  direct {}
}
```

## Developer: Using the `skip_on` struct field tag

The `skip_on` struct field tag is used to skip the external API call when only attributes that affect the internal state are modified, e.g. retry configuration. The `skip_on` struct field tag is used to skip the external API call when only attributes that affect the internal state are modified, e.g. retry configuration. The `skip_on` struct field tag is a comma-separated list of operations that must be met in order to skip the field.

The provider will compare the state with the plan, and check for changes. If the only fields to me modified are those with the `skip_on` struct field tag set to the supplied operation, e.g. `update`, the provider will skip the external API call.

The following operations are supported:

* `update` - Skip the external API call when the operation is an update.

## Additional topics

Refer to [contributing/README.md](contributing/README.md) for additional contribution topics.