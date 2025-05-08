---
subcategory: "Microsoft.DevTestLab - Azure Lab Services"
page_title: "labs/schedules"
description: |-
  Manages a automated startup and shutdown schedules for Azure Dev Test Lab.
---

# Microsoft.DevTestLab/labs/schedules - automated startup and shutdown schedules for Azure Dev Test Lab

This article demonstrates how to use `azapi` provider to manage the automated startup and shutdown schedules for Azure Dev Test Lab resource in Azure.

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

resource "azapi_resource" "lab" {
  type      = "Microsoft.DevTestLab/labs@2018-09-15"
  parent_id = azapi_resource.resourceGroup.id
  name      = var.resource_name
  location  = var.location
  body = {
    properties = {
      labStorageType = "Premium"
    }
  }
  schema_validation_enabled = false
  response_export_values    = ["*"]
}

resource "azapi_resource" "schedule" {
  type      = "Microsoft.DevTestLab/labs/schedules@2018-09-15"
  parent_id = azapi_resource.lab.id
  name      = "LabVmsShutdown"
  location  = var.location
  body = {
    properties = {
      dailyRecurrence = {
        time = "0100"
      }
      notificationSettings = {
        status        = "Disabled"
        timeInMinutes = 0
        webhookUrl    = ""
      }
      status     = "Disabled"
      taskType   = "LabVmsShutdownTask"
      timeZoneId = "India Standard Time"
    }
    tags = {
      environment = "Production"
    }
  }
  schema_validation_enabled = false
  response_export_values    = ["*"]
}


```



## Arguments Reference

The following arguments are supported:

* `type` - (Required) The type of the resource. This should be set to `Microsoft.DevTestLab/labs/schedules@api-version`. The available api-versions for this resource are: [`2015-05-21-preview`, `2016-05-15`, `2018-09-15`].

* `parent_id` - (Required) The ID of the azure resource in which this resource is created. The allowed values are:  
  `/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.DevTestLab/labs/{resourceName}`

* `name` - (Required) Specifies the name of the azure resource. Changing this forces a new resource to be created.

* `body` - (Required) Specifies the configuration of the resource. More information about the arguments in `body` can be found in the [Microsoft documentation](https://learn.microsoft.com/en-us/azure/templates/Microsoft.DevTestLab/labs/schedules?pivots=deployment-language-terraform).

## Import

 ```shell
 # Azure resource can be imported using the resource id, e.g.
 terraform import azapi_resource.example /subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.DevTestLab/labs/{resourceName}/schedules/{resourceName}
 
 # It also supports specifying API version by using the resource id with api-version as a query parameter, e.g.
 terraform import azapi_resource.example /subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.DevTestLab/labs/{resourceName}/schedules/{resourceName}?api-version=2018-09-15
 ```
