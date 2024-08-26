terraform {
  required_providers {
    azapi = {
      source = "Azure/azapi"
    }
    azurerm = {
      source = "hashicorp/azurerm"
    }
  }
}

provider "azurerm" {
  features {}
}

provider "azapi" {
}

data "azurerm_client_config" "current" {
}

data "azurerm_management_group" "tenant_root" {
  name = data.azurerm_client_config.current.tenant_id
}

data "azurerm_management_group" "default" {
  name = "my-default-managementgroup"
}

// Note:
// If the hierarchy settings object already exists,
// you can manage the properties without having to run `terraform import`.
// Just change the resource from "azapi_resource" to "azapi_update_resource"
resource "azapi_resource" "hierarchy_settings" {
  type      = "Microsoft.Management/managementGroups/settings@2021-04-01"
  name      = "default"
  parent_id = data.azurerm_management_group.tenant_root.id
  body = {
    properties = {
      defaultManagementGroup               = data.azurerm_management_group.default.id
      requireAuthorizationForGroupCreation = true
    }
  }
}
