---
subcategory: "Microsoft.TimeSeriesInsights - Azure Time Series Insights"
page_title: "environments/eventSources"
description: |-
  Manages a Time Series Insights Environments Event Sources.
---

# Microsoft.TimeSeriesInsights/environments/eventSources - Time Series Insights Environments Event Sources

This article demonstrates how to use `azapi` provider to manage the Time Series Insights Environments Event Sources resource in Azure.



## Example Usage

### default

```hcl
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


```



## Arguments Reference

The following arguments are supported:

* `type` - (Required) The type of the resource. This should be set to `Microsoft.TimeSeriesInsights/environments/eventSources@api-version`. The available api-versions for this resource are: [`2017-02-28-preview`, `2017-11-15`, `2018-08-15-preview`, `2020-05-15`, `2021-03-31-preview`, `2021-06-30-preview`].

* `parent_id` - (Required) The ID of the azure resource in which this resource is created. The allowed values are:  
  `/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.TimeSeriesInsights/environments/{resourceName}`

* `name` - (Required) Specifies the name of the azure resource. Changing this forces a new resource to be created.

* `body` - (Required) Specifies the configuration of the resource. More information about the arguments in `body` can be found in the [Microsoft documentation](https://learn.microsoft.com/en-us/azure/templates/Microsoft.TimeSeriesInsights/environments/eventSources?pivots=deployment-language-terraform).

For other arguments, please refer to the [azapi_resource](https://registry.terraform.io/providers/Azure/azapi/latest/docs/resources/resource) documentation.

## Import

 ```shell
 # Azure resource can be imported using the resource id, e.g.
 terraform import azapi_resource.example /subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.TimeSeriesInsights/environments/{resourceName}/eventSources/{resourceName}
 
 # It also supports specifying API version by using the resource id with api-version as a query parameter, e.g.
 terraform import azapi_resource.example /subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.TimeSeriesInsights/environments/{resourceName}/eventSources/{resourceName}?api-version=2021-06-30-preview
 ```
