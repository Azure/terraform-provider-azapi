# Terraform Provider for Azure Resource Manager Rest API

The AzAPI provider is a very thin layer on top of the [Azure ARM REST APIs](https://learn.microsoft.com/rest/api/azure). Use this new provider to authenticate to and manage Azure resources and functionality using the Azure Resource Manager APIs directly.

This provider complements the AzureRM provider by enabling the management of Azure resources that are not yet or may never be supported in the AzureRM provider such as private/public preview services and features.

## Get started with AzAPI

* [AzAPI Provider Overview](https://docs.microsoft.com/en-us/azure/developer/terraform/overview-azapi-provider)

* [Learn how to use azapi_resource](https://docs.microsoft.com/en-us/azure/developer/terraform/get-started-azapi-resource)

* [Learn how to use azapi_update_resource](https://docs.microsoft.com/en-us/azure/developer/terraform/get-started-azapi-update-resource)

* VS Code users: install the [Microsoft Terraform extension](https://marketplace.visualstudio.com/items?itemName=ms-azuretools.vscode-azureterraform) for authoring assistance (the former AzAPI VS Code extension has been deprecated and its features were merged here).

Also, there is a rich library of [examples](https://github.com/Azure/terraform-provider-azapi/tree/main/examples) to help you get started.

## Usage Example

The following example shows how to use `azapi_resource` to manage machine learning compute resource.

```hcl
terraform {
  required_providers {
    azapi = {
      source  = "Azure/azapi"
    }
  }
}

provider "azapi" {
  # More information on the authentication methods supported by
  # the AzAPI Provider can be found here:
  # https://registry.terraform.io/providers/Azure/azapi/latest/docs

  # subscription_id = "..."
  # client_id       = "..."
  # client_secret   = "..."
  # tenant_id       = "..."
}

// azurerm provider
provider "azurerm" {
  features {}
}

// create a resource group, it's recommended to use azapi provider with azurerm provider
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

// create a automation account
resource "azapi_resource" "automationAccount" {
  type      = "Microsoft.Automation/automationAccounts@2021-06-22"
  name      = "myAccount"
  parent_id = azurerm_resource_group.example.id

  location = azurerm_resource_group.example.location
  body = {
    properties = {
      disableLocalAuth    = true
      publicNetworkAccess = false
      sku = {
        name = "Basic"
      }
    }
  }
}

```

Further [usage documentation is available on the Terraform website](https://registry.terraform.io/providers/Azure/azapi/latest/docs).

## Contributing

Please see [CONTRIBUTING.md](CONTRIBUTING.md) for details on how to contribute to this project.

## Credits

We wish to thank HashiCorp for the use of some MPLv2-licensed code from their open source project [terraform-provider-azurerm](https://github.com/hashicorp/terraform-provider-azurerm).
