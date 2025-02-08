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

locals {
  data_collection_log_table_name = "DataCollectionLog_CL"
  data_collection_log_columns = [
    {
      "name" : "RawData",
      "type" : "string"
    },
    {
      "name" : "FilePath",
      "type" : "string"
    },
    {
      "name" : "TimeGenerated",
      "type" : "datetime"
    }
  ]
}

resource "azapi_resource" "resourceGroup" {
  type     = "Microsoft.Resources/resourceGroups@2020-06-01"
  name     = var.resource_name
  location = var.location
}

resource "azapi_resource" "workspace" {
  type      = "Microsoft.OperationalInsights/workspaces@2022-10-01"
  parent_id = azapi_resource.resourceGroup.id
  name      = var.resource_name
  location  = var.location
  body = {
    properties = {
      features = {
        disableLocalAuth                            = false
        enableLogAccessUsingOnlyResourcePermissions = true
      }
      publicNetworkAccessForIngestion = "Enabled"
      publicNetworkAccessForQuery     = "Enabled"
      retentionInDays                 = 30
      sku = {
        name = "PerGB2018"
      }
      workspaceCapping = {
        dailyQuotaGb = -1
      }
    }
  }
  schema_validation_enabled = false
  response_export_values    = ["*"]
}

resource "azapi_resource" "table" {
  type      = "Microsoft.OperationalInsights/workspaces/tables@2022-10-01"
  parent_id = azapi_resource.workspace.id
  name      = local.data_collection_log_table_name
  body = {
    properties = {
      schema = {
        name    = local.data_collection_log_table_name
        columns = local.data_collection_log_columns
      }
    }
  }
  schema_validation_enabled = false
  response_export_values    = ["*"]
}
