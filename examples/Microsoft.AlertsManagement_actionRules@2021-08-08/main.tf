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

resource "azapi_resource" "actionRule" {
  type      = "Microsoft.AlertsManagement/actionRules@2021-08-08"
  parent_id = azapi_resource.resourceGroup.id
  name      = var.resource_name
  location  = "global"
  body = {
    properties = {
      actions = [
        {
          actionType = "RemoveAllActionGroups"
        },
      ]
      description = ""
      enabled     = true
      scopes = [
        azapi_resource.resourceGroup.id,
      ]
    }

  }
  schema_validation_enabled = false
  response_export_values    = ["*"]
}

