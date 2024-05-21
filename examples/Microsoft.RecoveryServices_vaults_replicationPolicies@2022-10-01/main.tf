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
  type     = "Microsoft.Resources/resourceGroups@2020-06-01"
  name     = var.resource_name
  location = var.location
}

resource "azapi_resource" "vault" {
  type      = "Microsoft.RecoveryServices/vaults@2022-10-01"
  parent_id = azapi_resource.resourceGroup.id
  name      = var.resource_name
  location  = var.location
  body = {
    properties = {
      publicNetworkAccess = "Enabled"
    }
    sku = {
      name = "Standard"
    }
  }
  schema_validation_enabled = false
  response_export_values    = ["*"]
}

resource "azapi_resource" "replicationPolicy" {
  type      = "Microsoft.RecoveryServices/vaults/replicationPolicies@2022-10-01"
  parent_id = azapi_resource.vault.id
  name      = var.resource_name
  body = {
    properties = {
      providerSpecificInput = {
        appConsistentFrequencyInMinutes   = 240
        crashConsistentFrequencyInMinutes = 10
        enableMultiVmSync                 = "True"
        instanceType                      = "InMageRcm"
        recoveryPointHistoryInMinutes     = 1440
      }
    }
  }
  schema_validation_enabled = false
  response_export_values    = ["*"]
}

