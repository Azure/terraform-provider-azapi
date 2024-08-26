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

resource "azapi_resource" "storageAccount" {
  type      = "Microsoft.Storage/storageAccounts@2021-09-01"
  parent_id = azapi_resource.resourceGroup.id
  name      = var.resource_name
  location  = var.location
  body = {
    kind = "StorageV2"
    properties = {
      accessTier                   = "Hot"
      allowBlobPublicAccess        = true
      allowCrossTenantReplication  = true
      allowSharedKeyAccess         = true
      defaultToOAuthAuthentication = false
      encryption = {
        keySource = "Microsoft.Storage"
        services = {
          queue = {
            keyType = "Service"
          }
          table = {
            keyType = "Service"
          }
        }
      }
      isHnsEnabled      = false
      isNfsV3Enabled    = false
      isSftpEnabled     = false
      minimumTlsVersion = "TLS1_2"
      networkAcls = {
        defaultAction = "Allow"
      }
      publicNetworkAccess      = "Enabled"
      supportsHttpsTrafficOnly = true
    }
    sku = {
      name = "Standard_LRS"
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
      name     = "B1"
    }
    tags = {
      purpose = "testing"
    }
  }
  schema_validation_enabled = false
  response_export_values    = ["*"]
}

data "azapi_resource_action" "listKeys" {
  type                   = "Microsoft.Storage/storageAccounts@2021-09-01"
  resource_id            = azapi_resource.storageAccount.id
  action                 = "listKeys"
  response_export_values = ["*"]
}

resource "azapi_resource" "environment" {
  type      = "Microsoft.TimeSeriesInsights/environments@2020-05-15"
  parent_id = azapi_resource.resourceGroup.id
  name      = var.resource_name
  location  = var.location
  body = {
    kind = "Gen2"
    properties = {
      storageConfiguration = {
        accountName   = azapi_resource.storageAccount.name
        managementKey = data.azapi_resource_action.listKeys.output.keys[0].value
      }
      timeSeriesIdProperties = [
        {
          name = "id"
          type = "String"
        },
      ]
    }
    sku = {
      capacity = 1
      name     = "L1"
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

resource "azapi_resource" "eventSource" {
  type      = "Microsoft.TimeSeriesInsights/environments/eventSources@2020-05-15"
  parent_id = azapi_resource.environment.id
  name      = var.resource_name
  location  = var.location
  body = {
    kind = "Microsoft.IoTHub"
    properties = {
      consumerGroupName     = "test"
      eventSourceResourceId = azapi_resource.IotHub.id
      iotHubName            = azapi_resource.IotHub.name
      keyName               = "iothubowner"
      sharedAccessKey       = data.azapi_resource_action.listkeys.output.value[0].primaryKey
      timestampPropertyName = ""
    }
  }
  schema_validation_enabled = false
  response_export_values    = ["*"]
}

