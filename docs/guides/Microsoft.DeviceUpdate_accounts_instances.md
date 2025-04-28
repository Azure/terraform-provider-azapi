---
subcategory: "Microsoft.DeviceUpdate - Device Update for IoT Hub"
page_title: "accounts/instances"
description: |-
  Manages a IoT Hub Device Update Instance.
---

# Microsoft.DeviceUpdate/accounts/instances - IoT Hub Device Update Instance

This article demonstrates how to use `azapi` provider to manage the IoT Hub Device Update Instance resource in Azure.

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

resource "azapi_resource" "account" {
  type      = "Microsoft.DeviceUpdate/accounts@2022-10-01"
  parent_id = azapi_resource.resourceGroup.id
  name      = var.resource_name
  location  = var.location
  body = {
    properties = {
      publicNetworkAccess = "Enabled"
      sku                 = "Standard"
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

resource "azapi_resource" "instance" {
  type      = "Microsoft.DeviceUpdate/accounts/instances@2022-10-01"
  parent_id = azapi_resource.account.id
  name      = var.resource_name
  location  = var.location
  body = {
    properties = {
      accountName       = azapi_resource.account.name
      enableDiagnostics = false
      iotHubs = [
        {
          resourceId = azapi_resource.IotHub.id
        },
      ]
    }

  }
  schema_validation_enabled = false
  response_export_values    = ["*"]
}


```



## Arguments Reference

The following arguments are supported:

* `type` - (Required) The type of the resource. This should be set to `Microsoft.DeviceUpdate/accounts/instances@api-version`. The available api-versions for this resource are: [`2020-03-01-preview`, `2022-04-01-preview`, `2022-10-01`, `2022-12-01-preview`, `2023-07-01`].

* `parent_id` - (Required) The ID of the azure resource in which this resource is created. The allowed values are:  
  `/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.DeviceUpdate/accounts/{resourceName}`

* `name` - (Required) Specifies the name of the azure resource. Changing this forces a new resource to be created.

* `body` - (Required) Specifies the configuration of the resource. More information about the arguments in `body` can be found in the [Microsoft documentation](https://learn.microsoft.com/en-us/azure/templates/Microsoft.DeviceUpdate/accounts/instances?pivots=deployment-language-terraform).

## Import

 ```shell
 # Azure resource can be imported using the resource id, e.g.
 terraform import azapi_resource.example /subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.DeviceUpdate/accounts/{resourceName}/instances/{resourceName}
 
 # It also supports specifying API version by using the resource id with api-version as a query parameter, e.g.
 terraform import azapi_resource.example /subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.DeviceUpdate/accounts/{resourceName}/instances/{resourceName}?api-version=2023-07-01
 ```
