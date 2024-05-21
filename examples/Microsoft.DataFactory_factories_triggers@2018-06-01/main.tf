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

resource "azapi_resource" "factory" {
  type      = "Microsoft.DataFactory/factories@2018-06-01"
  parent_id = azapi_resource.resourceGroup.id
  name      = var.resource_name
  location  = var.location
  body = {
    properties = {
      publicNetworkAccess = "Enabled"
      repoConfiguration   = null
    }
  }
  schema_validation_enabled = false
  response_export_values    = ["*"]
}

resource "azapi_resource" "pipeline" {
  type      = "Microsoft.DataFactory/factories/pipelines@2018-06-01"
  parent_id = azapi_resource.factory.id
  name      = var.resource_name
  body = {
    properties = {
      annotations = [
      ]
      description = ""
      parameters = {
        test = {
          defaultValue = "testparameter"
          type         = "String"
        }
      }
      variables = {
      }
    }
  }
  schema_validation_enabled = false
  response_export_values    = ["*"]
}

resource "azapi_resource" "trigger" {
  type      = "Microsoft.DataFactory/factories/triggers@2018-06-01"
  parent_id = azapi_resource.factory.id
  name      = var.resource_name
  body = {
    properties = {
      description = ""
      pipeline = {
        parameters = {
        }
        pipelineReference = {
          referenceName = azapi_resource.pipeline.name
          type          = "PipelineReference"
        }
      }
      type = "TumblingWindowTrigger"
      typeProperties = {
        frequency      = "Minute"
        interval       = 15
        maxConcurrency = 50
        startTime      = "2022-09-21T00:00:00Z"
      }
    }
  }
  schema_validation_enabled = false
  response_export_values    = ["*"]
}

