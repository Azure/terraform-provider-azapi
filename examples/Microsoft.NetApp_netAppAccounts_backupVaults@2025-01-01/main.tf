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
  default = "westus"
}

resource "azapi_resource" "resourceGroup" {
  type     = "Microsoft.Resources/resourceGroups@2020-06-01"
  name     = var.resource_name
  location = var.location
  tags = {
    SkipNRMSNSG = "true"
  }
}

resource "azapi_resource" "netAppAccount" {
  type      = "Microsoft.NetApp/netAppAccounts@2025-01-01"
  parent_id = azapi_resource.resourceGroup.id
  name      = var.resource_name
  location  = var.location
  body = {
    properties = {}
  }
  tags = {
    CreatedOnDate    = "2023-08-17T08:01:00Z"
    SkipASMAzSecPack = "true"
  }
}

resource "azapi_resource" "backupVault" {
  type      = "Microsoft.NetApp/netAppAccounts/backupVaults@2025-01-01"
  parent_id = azapi_resource.netAppAccount.id
  name      = "${var.resource_name}-backupvault"
  location  = var.location
  tags = {
    testTag = "testTagValue"
  }
}
