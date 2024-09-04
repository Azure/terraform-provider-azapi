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
          query          = "   SELECT *\n   INTO [YourOutputAlias]\n   FROM [YourInputAlias]\n"
          streamingUnits = 3
        }
      }
    }
  }
  schema_validation_enabled = false
  response_export_values    = ["*"]
}

resource "azapi_resource" "IotHub" {
  type      = "Microsoft.Devices/IotHubs@2022-04-30-preview"
  parent_id = azapi_resource.resourceGroup.id
  name      = var.resource_name
  location  = var.location
  body = {
    properties = {
      cloudToDevice = {
      }
      enableFileUploadNotifications = false
      messagingEndpoints = {
      }
      routing = {
        fallbackRoute = {
          condition = "true"
          endpointNames = [
            "events",
          ]
          isEnabled = true
          source    = "DeviceMessages"
        }
      }
      storageEndpoints = {
      }
    }
    sku = {
      capacity = 1
      name     = "S1"
    }
  }
  schema_validation_enabled = false
  response_export_values    = ["*"]
}

data "azapi_resource_action" "listkeys" {
  type                   = "Microsoft.Devices/IotHubs@2022-04-30-preview"
  resource_id            = azapi_resource.IotHub.id
  action                 = "listkeys"
  response_export_values = ["*"]
}

resource "azapi_resource" "input" {
  type      = "Microsoft.StreamAnalytics/streamingJobs/inputs@2020-03-01"
  parent_id = azapi_resource.streamingJob.id
  name      = var.resource_name
  body = {
    properties = {
      datasource = {
        properties = {
          consumerGroupName      = "$Default"
          endpoint               = "messages/events"
          iotHubNamespace        = azapi_resource.IotHub.name
          sharedAccessPolicyKey  = data.azapi_resource_action.listkeys.output.value[0].primaryKey
          sharedAccessPolicyName = "iothubowner"
        }
        type = "Microsoft.Devices/IotHubs"
      }
      serialization = {
        properties = {}
        type       = "Avro"
      }
      type = "Stream"
    }
  }
  schema_validation_enabled = false
  response_export_values    = ["*"]
}

