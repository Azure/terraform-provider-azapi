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

* `environment` - (Optional) The Cloud Environment which should be used. Possible values are `public`, `usgovernment` and `china`. Defaults to `public`. This can also be sourced from the `ARM_ENVIRONMENT` Environment Variable.

* `subscription_id` - (Optional) The Subscription ID which should be used. This can also be sourced from the `ARM_SUBSCRIPTION_ID` Environment Variable.

* `tenant_id` - (Optional) The Tenant ID should be used. This can also be sourced from the `ARM_TENANT_ID` Environment Variable.

---

It's possible to configure the behaviour of certain resources using the following properties: 

* `default_tags` - (Optional) A mapping of tags which should be assigned to the azure resource as default tags. `tags` in each resource block can override the `default_tags`.

* `default_location` - (Optional) The default Azure Region where the azure resource should exist. `location` in each resource block can override the `default_location`. Changing this forces new resources to be created.
---

When authenticating as a Service Principal using a Client Certificate, the following fields can be set:

* `client_certificate_path` - (Optional) The path to the Client Certificate associated with the Service Principal which should be used. This can also be sourced from the `ARM_CLIENT_CERTIFICATE_PATH` Environment Variable.

More information on [how to configure a Service Principal using a Client Certificate can be found in this guide](guides/service_principal_client_certificate.html).

---

When authenticating as a Service Principal using a Client Secret, the following fields can be set:

* `client_secret` - (Optional) The Client Secret which should be used. This can also be sourced from the `ARM_CLIENT_SECRET` Environment Variable.

More information on [how to configure a Service Principal using a Client Secret can be found in this guide](guides/service_principal_client_secret.html).
