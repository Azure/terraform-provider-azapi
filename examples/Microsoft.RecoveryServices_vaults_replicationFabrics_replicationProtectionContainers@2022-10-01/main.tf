terraform {
  required_providers {
    azapi = {
      source = "Azure/azapi"
    }
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

resource "azapi_resource" "resourceGroup" {
  type                      = "Microsoft.Resources/resourceGroups@2020-06-01"
  name                      = var.resource_name
  location                  = var.location
}

resource "azapi_resource" "vault" {
  type      = "Microsoft.RecoveryServices/vaults@2022-10-01"
  parent_id = azapi_resource.resourceGroup.id
  name      = var.resource_name
  location  = var.location
  body = jsonencode({
    properties = {
      publicNetworkAccess = "Enabled"
    }
    sku = {
      name = "Standard"
    }
  })
  schema_validation_enabled = false
  response_export_values    = ["*"]
}

resource "azapi_resource" "replicationFabric" {
  type      = "Microsoft.RecoveryServices/vaults/replicationFabrics@2022-10-01"
  parent_id = azapi_resource.vault.id
  name      = var.resource_name
  body = jsonencode({
    properties = {
      customDetails = {
        instanceType = "Azure"
        location     = var.location
      }
    }
  })
  schema_validation_enabled = false
  response_export_values    = ["*"]
}

resource "azapi_resource" "replicationProtectionContainer" {
  type      = "Microsoft.RecoveryServices/vaults/replicationFabrics/replicationProtectionContainers@2022-10-01"
  parent_id = azapi_resource.replicationFabric.id
  name      = var.resource_name
  body = jsonencode({
    properties = {
    }
  })
  schema_validation_enabled = false
  response_export_values    = ["*"]
}

