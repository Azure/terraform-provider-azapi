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

resource "azapi_resource" "serverGroupsv2" {
  type      = "Microsoft.DBforPostgreSQL/serverGroupsv2@2022-11-08"
  parent_id = azapi_resource.resourceGroup.id
  name      = var.resource_name
  location  = var.location
  body = {
    properties = {
      administratorLoginPassword      = "H@Sh1CoR3!"
      coordinatorEnablePublicIpAccess = true
      coordinatorServerEdition        = "GeneralPurpose"
      coordinatorStorageQuotaInMb     = 131072
      coordinatorVCores               = 2
      enableHa                        = false
      nodeCount                       = 0
      nodeEnablePublicIpAccess        = false
      nodeServerEdition               = "MemoryOptimized"
    }
  }
  schema_validation_enabled = false
  response_export_values    = ["*"]
}

