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

resource "azapi_resource" "resourceGroup" {
  type     = "Microsoft.Resources/resourceGroups@2020-06-01"
  name     = var.resource_name
  location = var.location
}

resource "azapi_resource" "managedHSM" {
  type      = "Microsoft.KeyVault/managedHSMs@2021-10-01"
  parent_id = azapi_resource.resourceGroup.id
  name      = "kvHsm230630033342437496"
  location  = var.location
  body = {
    properties = {
      createMode            = "default"
      enablePurgeProtection = false
      enableSoftDelete      = true
      initialAdminObjectIds = [
        data.azurerm_client_config.current.object_id,
      ]
      publicNetworkAccess       = "Enabled"
      softDeleteRetentionInDays = 90
      tenantId                  = data.azurerm_client_config.current.tenant_id
    }
    sku = {
      family = "B"
      name   = "Standard_B1"
    }
  }
  schema_validation_enabled = false
  response_export_values    = ["*"]
}

