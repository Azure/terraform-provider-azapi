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
  features {
  }
}

provider "azapi" {
  skip_provider_registration = false
}

variable "resource_name" {
  type    = string
  default = "acctest0001"
}

variable "location" {
  type    = string
  default = "westeurope"
}

data "azurerm_client_config" "current" {
}

data "azapi_resource_action" "roleDefinitions" {
  type                   = "Microsoft.Authorization@2018-01-01-preview"
  resource_id            = "/providers/Microsoft.Authorization"
  action                 = "roleDefinitions"
  method                 = "GET"
  response_export_values = ["*"]
}

resource "azapi_resource" "resourceGroup" {
  type     = "Microsoft.Resources/resourceGroups@2020-06-01"
  name     = var.resource_name
  location = var.location
}

resource "azapi_resource" "applicationDefinition" {
  type      = "Microsoft.Solutions/applicationDefinitions@2021-07-01"
  parent_id = azapi_resource.resourceGroup.id
  name      = var.resource_name
  location  = var.location
  body = {
    properties = {
      authorizations = [
        {
          principalId      = data.azurerm_client_config.current.object_id
          roleDefinitionId = data.azapi_resource_action.roleDefinitions.output.value[0].name
        },
      ]
      description    = "Test Managed App Definition"
      displayName    = "TestManagedAppDefinition"
      isEnabled      = true
      lockLevel      = "ReadOnly"
      packageFileUri = "https://github.com/Azure/azure-managedapp-samples/raw/master/Managed Application Sample Packages/201-managed-storage-account/managedstorage.zip"
    }
  }
  schema_validation_enabled = false
  response_export_values    = ["*"]
}

