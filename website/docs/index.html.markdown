---
layout: "azurerm-restapi"
page_title: "Provider: Azure"
description: |-
  The AzureRM Rest API Provider is used to interact with the many resources supported by Azure Resource Manager (also known as AzureRM) through its APIs.

---

# AzureRM Rest API Provider

The AzureRM Rest API Provider can be used to configure infrastructure in [Microsoft Azure](https://azure.microsoft.com/en-us/) using the Azure Resource Manager API's. Documentation regarding the [Data Sources](/docs/configuration/data-sources.html) and [Resources](/docs/configuration/resources.html) supported by the Azure Provider can be found in the navigation to the left.

To learn the basics of Terraform using this provider, follow the
hands-on [get started tutorials](https://learn.hashicorp.com/tutorials/terraform/infrastructure-as-code?in=terraform/azure-get-started) on HashiCorp's Learn platform.

Interested in the provider's latest features, or want to make sure you're up to date? Check out the [changelog](https://github.com/Azure/terraform-provider-azurerm-generic/blob/develop/CHANGELOG.md) for version information and release notes.

## Authenticating to Azure

Terraform supports a number of different methods for authenticating to Azure:

* [Authenticating to Azure using the Azure CLI](guides/azure_cli.html)
* [Authenticating to Azure using Managed Service Identity](guides/managed_service_identity.html)
* [Authenticating to Azure using a Service Principal and a Client Certificate](guides/service_principal_client_certificate.html)
* [Authenticating to Azure using a Service Principal and a Client Secret](guides/service_principal_client_secret.html)

---

We recommend using either a Service Principal or Managed Service Identity when running Terraform non-interactively (such as when running Terraform in a CI server) - and authenticating using the Azure CLI when running Terraform locally.

## Example Usage

```hcl
# We strongly recommend using the required_providers block to set the
# Azure Provider source and version being used
terraform {
  required_providers {
    azurerm-restapi = {
      source = "Azure/azurerm-restapi"
    }
  }
}

provider "azurerm-restapi" {
}

```

## Argument Reference

The following arguments are supported:

* `client_id` - (Optional) The Client ID which should be used. This can also be sourced from the `ARM_CLIENT_ID` Environment Variable.

* `environment` - (Optional) The Cloud Environment which should be used. Possible values are `public`, `usgovernment`, `german`, and `china`. Defaults to `public`. This can also be sourced from the `ARM_ENVIRONMENT` Environment Variable.

* `subscription_id` - (Optional) The Subscription ID which should be used. This can also be sourced from the `ARM_SUBSCRIPTION_ID` Environment Variable.

* `tenant_id` - (Optional) The Tenant ID should be used. This can also be sourced from the `ARM_TENANT_ID` Environment Variable.

---

It's possible to configure the behaviour of certain resources using the following properties: 

* `schema_validation_enabled` - (Optional) Whether enabled the validation on `type` and `body` with embedded schema. Defaults to `true`.

---

When authenticating as a Service Principal using a Client Certificate, the following fields can be set:

* `client_certificate_password` - (Optional) The password associated with the Client Certificate. This can also be sourced from the `ARM_CLIENT_CERTIFICATE_PASSWORD` Environment Variable.

* `client_certificate_path` - (Optional) The path to the Client Certificate associated with the Service Principal which should be used. This can also be sourced from the `ARM_CLIENT_CERTIFICATE_PATH` Environment Variable.

More information on [how to configure a Service Principal using a Client Certificate can be found in this guide](guides/service_principal_client_certificate.html).

---

When authenticating as a Service Principal using a Client Secret, the following fields can be set:

* `client_secret` - (Optional) The Client Secret which should be used. This can also be sourced from the `ARM_CLIENT_SECRET` Environment Variable.

More information on [how to configure a Service Principal using a Client Secret can be found in this guide](guides/service_principal_client_secret.html).

---

When authenticating using Managed Service Identity, the following fields can be set:

* `msi_endpoint` - (Optional) The path to a custom endpoint for Managed Service Identity - in most circumstances, this should be detected automatically. This can also, be sourced from the `ARM_MSI_ENDPOINT` Environment Variable.

* `use_msi` - (Optional) Should Managed Service Identity be used for Authentication? This can also be sourced from the `ARM_USE_MSI` Environment Variable. Defaults to `false`.

More information on [how to configure a Service Principal using Managed Service Identity can be found in this guide](guides/managed_service_identity.html).

---

For some advanced scenarios, such as where more granular permissions are necessary - the following properties can be set:

* `metadata_host` - (Optional) The Hostname of the Azure Metadata Service (for example `management.azure.com`), used to obtain the Cloud Environment when using a Custom Azure Environment. This can also be sourced from the `ARM_METADATA_HOSTNAME` Environment Variable.

~> **Note:** `environment` must be set to the requested environment name in the list of available environments held in the `metadata_host`.
