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
  default = "eastus2"
}

resource "azapi_resource" "resourceGroup" {
  type     = "Microsoft.Resources/resourceGroups@2020-06-01"
  name     = var.resource_name
  location = var.location
}

resource "azapi_resource" "netAppAccount" {
  type      = "Microsoft.NetApp/netAppAccounts@2022-05-01"
  parent_id = azapi_resource.resourceGroup.id
  name      = var.resource_name
  location  = var.location
  body = {
    properties = {
      activeDirectories = [
      ]
    }
  }
  schema_validation_enabled = false
  response_export_values    = ["*"]
}

resource "azapi_resource" "snapshotPolicy" {
  type      = "Microsoft.NetApp/netAppAccounts/snapshotPolicies@2022-05-01"
  parent_id = azapi_resource.netAppAccount.id
  name      = var.resource_name
  location  = var.location
  body = {
    properties = {
      dailySchedule = {
        hour            = 22
        minute          = 15
        snapshotsToKeep = 1
      }
      enabled = true
      hourlySchedule = {
        minute          = 15
        snapshotsToKeep = 1
      }
      monthlySchedule = {
        daysOfMonth     = "30,15,1"
        hour            = 5
        minute          = 0
        snapshotsToKeep = 1
      }
      weeklySchedule = {
        day             = "Monday,Friday"
        hour            = 23
        minute          = 0
        snapshotsToKeep = 1
      }
    }
  }
  schema_validation_enabled = false
  response_export_values    = ["*"]
}

