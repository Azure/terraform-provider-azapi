---
layout: "azapi"
page_title: "Provider: Azure API"
description: |-
  The AzAPI Provider is used to interact with the many resources supported by Azure Resource Manager through its APIs.

---

# AzAPI Provider

The AzAPI provider is a very thin layer on top of the Azure ARM REST APIs. This provider compliments the [AzureRM provider](https://registry.terraform.io/providers/hashicorp/azurerm/latest/docs) by enabling the management of Azure resources that are not yet or may never be supported in the AzureRM provider such as private/public preview services and features. 

Documentation regarding the [Data Sources](/docs/configuration/data-sources.html) and [Resources](/docs/configuration/resources.html) supported by the AzAPI Provider can be found in the navigation to the left.

Interested in the provider's latest features, or want to make sure you're up to date? Check out the [changelog](https://github.com/Azure/terraform-provider-azapi/blob/main/CHANGELOG.md) for version information and release notes.

Strongly recommended to install [AzApi VSCode Extension](https://marketplace.visualstudio.com/items?itemName=azapi-vscode.azapi), it provides a rich authoring experience to help you use the AzApi provider.

Also, there is a rich library of [examples](https://github.com/Azure/terraform-provider-azapi/tree/main/examples) to help you get started.

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
    azapi = {
      source = "azure/azapi"
    }
  }
}

provider "azapi" {
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

* `client_certificate_password` - (Optional) The password associated with the Client Certificate. This can also be sourced from the `ARM_CLIENT_CERTIFICATE_PASSWORD` Environment Variable.

* `client_certificate_path` - (Optional) The path to the Client Certificate associated with the Service Principal which should be used. This can also be sourced from the `ARM_CLIENT_CERTIFICATE_PATH` Environment Variable.

More information on [how to configure a Service Principal using a Client Certificate can be found in this guide](guides/service_principal_client_certificate.html).

---

When authenticating as a Service Principal using a Client Secret, the following fields can be set:

* `client_secret` - (Optional) The Client Secret which should be used. This can also be sourced from the `ARM_CLIENT_SECRET` Environment Variable.

More information on [how to configure a Service Principal using a Client Secret can be found in this guide](guides/service_principal_client_secret.html).

---

For some advanced scenarios, such as where more granular permissions are necessary - the following properties can be set:

* `disable_correlation_request_id` - (Optional) Disable sending the `x-ms-correlation-request-id` header. This can also be sourced from the `ARM_DISABLE_CORRELATION_REQUEST_ID` environment variable. Defaults to `false`.

* `disable_terraform_partner_id` - (Optional) Disable sending the Terraform Partner ID if a custom `partner_id` isn't specified, which allows Microsoft to better understand the usage of Terraform. The Partner ID does not give HashiCorp any direct access to usage information. This can also be sourced from the `ARM_DISABLE_TERRAFORM_PARTNER_ID` environment variable. Defaults to `false`.

* `partner_id` - (Optional) A GUID/UUID that is [registered](https://docs.microsoft.com/azure/marketplace/azure-partner-customer-usage-attribution#register-guids-and-offers) with Microsoft to facilitate partner resource usage attribution. This can also be sourced from the `ARM_PARTNER_ID` Environment Variable.

* `skip_provider_registration` - (Optional) Should the Provider skip registering the Resource Providers it supports? This can also be sourced from the `ARM_SKIP_PROVIDER_REGISTRATION` Environment Variable. Defaults to `false`.

-> By default, Terraform will attempt to register the Resource Providers that the provisioning resources belong to. If you're running in an environment with restricted permissions, or wish to manage Resource Provider Registration outside of Terraform you may wish to disable this flag; however, please note that the error messages returned from Azure may be confusing as a result (example: `API version 2019-01-01 was not found for Microsoft.Foo`).
