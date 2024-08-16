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

resource "azapi_resource" "streamingJob" {
  type      = "Microsoft.StreamAnalytics/streamingJobs@2020-03-01"
  parent_id = azapi_resource.resourceGroup.id
  name      = var.resource_name
  location  = var.location
  body = {
    properties = {
      cluster = {
      }
      compatibilityLevel                 = "1.0"
      contentStoragePolicy               = "SystemAccount"
      dataLocale                         = "en-GB"
      eventsLateArrivalMaxDelayInSeconds = 60
      eventsOutOfOrderMaxDelayInSeconds  = 50
      eventsOutOfOrderPolicy             = "Adjust"
      jobType                            = "Cloud"
      outputErrorPolicy                  = "Drop"
      sku = {
        name = "Standard"
      }
      transformation = {
        name = "main"
        properties = {
          query          = "    SELECT *\n    INTO [YourOutputAlias]\n    FROM [YourInputAlias]\n"
          streamingUnits = 3
        }
      }
    }
  }
  schema_validation_enabled = false
  response_export_values    = ["*"]
}

resource "azapi_resource" "function" {
  type      = "Microsoft.StreamAnalytics/streamingJobs/functions@2020-03-01"
  parent_id = azapi_resource.streamingJob.id
  name      = var.resource_name
  body = {
    properties = {
      properties = {
        binding = {
          properties = {
            script = "function getRandomNumber(in) {\n  return in;\n}\n"
          }
          type = "Microsoft.StreamAnalytics/JavascriptUdf"
        }
        inputs = [
          {
            dataType                 = "bigint"
            isConfigurationParameter = false
          },
        ]
        output = {
          dataType = "bigint"
        }
      }
      type = "Scalar"
    }
  }
  schema_validation_enabled = false
  response_export_values    = ["*"]
}

