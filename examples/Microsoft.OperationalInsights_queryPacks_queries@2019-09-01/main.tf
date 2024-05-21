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

resource "azapi_resource" "queryPack" {
  type      = "Microsoft.OperationalInsights/queryPacks@2019-09-01"
  parent_id = azapi_resource.resourceGroup.id
  name      = var.resource_name
  location  = var.location
  body = {
    properties = {
    }

  }
  schema_validation_enabled = false
  response_export_values    = ["*"]
}

resource "azapi_resource" "query" {
  type      = "Microsoft.OperationalInsights/queryPacks/queries@2019-09-01"
  parent_id = azapi_resource.queryPack.id
  name      = "aca50e92-d3e6-8f7d-1f70-2ec7adc1a926"
  body = {
    properties = {
      body        = "    let newExceptionsTimeRange = 1d;\n    let timeRangeToCheckBefore = 7d;\n    exceptions\n    | where timestamp < ago(timeRangeToCheckBefore)\n    | summarize count() by problemId\n    | join kind= rightanti (\n        exceptions\n        | where timestamp >= ago(newExceptionsTimeRange)\n        | extend stack = tostring(details[0].rawStack)\n        | summarize count(), dcount(user_AuthenticatedId), min(timestamp), max(timestamp), any(stack) by problemId\n    ) on problemId\n    | order by count_ desc\n"
      displayName = "Exceptions - New in the last 24 hours"
      related = {
      }
    }
  }
  schema_validation_enabled = false
  response_export_values    = ["*"]
}

