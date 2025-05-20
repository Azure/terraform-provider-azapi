---
subcategory: "Microsoft.Devices - Azure IoT Hub, Azure IoT Hub Device Provisioning Service"
page_title: "provisioningServices"
description: |-
  Manages a IoT Device Provisioning Service.
---

# Microsoft.Devices/provisioningServices - IoT Device Provisioning Service

This article demonstrates how to use `azapi` provider to manage the IoT Device Provisioning Service resource in Azure.

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

resource "azapi_resource" "provisioningService" {
  type      = "Microsoft.Devices/provisioningServices@2022-02-05"
  parent_id = azapi_resource.resourceGroup.id
  name      = var.resource_name
  location  = var.location
  body = {
    properties = {
      allocationPolicy    = "Hashed"
      enableDataResidency = false
      iotHubs = [
      ]
      publicNetworkAccess = "Enabled"
    }
    sku = {
      capacity = 1
      name     = "S1"
    }
  }
  schema_validation_enabled = false
  response_export_values    = ["*"]
}


```



## Arguments Reference

The following arguments are supported:

* `type` - (Required) The type of the resource. This should be set to `Microsoft.Devices/provisioningServices@api-version`. The available api-versions for this resource are: [`2017-08-21-preview`, `2017-11-15`, `2018-01-22`, `2020-01-01`, `2020-03-01`, `2020-09-01-preview`, `2021-10-15`, `2022-02-05`, `2022-12-12`, `2023-03-01-preview`, `2025-02-01-preview`].

* `parent_id` - (Required) The ID of the azure resource in which this resource is created. The allowed values are:  
  `/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}`

* `name` - (Required) Specifies the name of the azure resource. Changing this forces a new resource to be created.

* `body` - (Required) Specifies the configuration of the resource. More information about the arguments in `body` can be found in the [Microsoft documentation](https://learn.microsoft.com/en-us/azure/templates/Microsoft.Devices/provisioningServices?pivots=deployment-language-terraform).

For other arguments, please refer to the [azapi_resource](https://registry.terraform.io/providers/Azure/azapi/latest/docs/resources/resource) documentation.

## Import

 ```shell
 # Azure resource can be imported using the resource id, e.g.
 terraform import azapi_resource.example /subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Devices/provisioningServices/{resourceName}
 
 # It also supports specifying API version by using the resource id with api-version as a query parameter, e.g.
 terraform import azapi_resource.example /subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Devices/provisioningServices/{resourceName}?api-version=2025-02-01-preview
 ```
