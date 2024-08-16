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

resource "azapi_resource" "lab" {
  type      = "Microsoft.DevTestLab/labs@2018-09-15"
  parent_id = azapi_resource.resourceGroup.id
  name      = var.resource_name
  location  = var.location
  body = {
    properties = {
      labStorageType = "Premium"
    }
  }
  schema_validation_enabled = false
  response_export_values    = ["*"]
}

resource "azapi_resource" "schedule" {
  type      = "Microsoft.DevTestLab/labs/schedules@2018-09-15"
  parent_id = azapi_resource.lab.id
  name      = "LabVmsShutdown"
  location  = var.location
  body = {
    properties = {
      dailyRecurrence = {
        time = "0100"
      }
      notificationSettings = {
        status        = "Disabled"
        timeInMinutes = 0
        webhookUrl    = ""
      }
      status     = "Disabled"
      taskType   = "LabVmsShutdownTask"
      timeZoneId = "India Standard Time"
    }
    tags = {
      environment = "Production"
    }
  }
  schema_validation_enabled = false
  response_export_values    = ["*"]
}

