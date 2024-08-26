terraform {
  required_providers {
    azapi = {
      source = "Azure/azapi"
    }
  }
}

provider "azapi" {
  enable_hcl_output_for_data_source = true
}

provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "example" {
  name     = "example-rg"
  location = "west europe"
}

resource "azurerm_container_registry" "example" {
  name                = "example"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location
  sku                 = "Premium"
  admin_enabled       = false
}

data "azapi_resource" "example" {
  name      = "example"
  parent_id = azurerm_resource_group.example.id
  type      = "Microsoft.ContainerRegistry/registries@2020-11-01-preview"

  response_export_values = ["properties.loginServer", "properties.policies.quarantinePolicy.status"]
}

// it will output "registry1.azurecr.io"
output "login_server" {
  value = data.azapi_resource.example.output.properties.loginServer
}

// it will output "disabled"
output "quarantine_policy" {
  value = data.azapi_resource.example.output.properties.policies.quarantinePolicy.status
}
