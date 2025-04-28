---
subcategory: "Microsoft.Maintenance - Azure Maintenance"
page_title: "maintenanceConfigurations"
description: |-
  Manages a Maintenance Configuration.
---

# Microsoft.Maintenance/maintenanceConfigurations - Maintenance Configuration

This article demonstrates how to use `azapi` provider to manage the Maintenance Configuration resource in Azure.

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

resource "azapi_resource" "maintenanceConfiguration" {
  type      = "Microsoft.Maintenance/maintenanceConfigurations@2022-07-01-preview"
  parent_id = azapi_resource.resourceGroup.id
  name      = var.resource_name
  location  = var.location
  body = {
    properties = {
      extensionProperties = {
      }
      maintenanceScope = "SQLDB"
      namespace        = "Microsoft.Maintenance"
      visibility       = "Custom"
    }
  }
  schema_validation_enabled = false
  response_export_values    = ["*"]
}


```



## Arguments Reference

The following arguments are supported:

* `type` - (Required) The type of the resource. This should be set to `Microsoft.Maintenance/maintenanceConfigurations@api-version`. The available api-versions for this resource are: [`2018-06-01-preview`, `2020-04-01`, `2020-07-01-preview`, `2021-04-01-preview`, `2021-05-01`, `2021-09-01-preview`, `2022-07-01-preview`, `2022-11-01-preview`, `2023-04-01`, `2023-09-01-preview`, `2023-10-01-preview`].

* `parent_id` - (Required) The ID of the azure resource in which this resource is created. The allowed values are:  
  `/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}`

* `name` - (Required) Specifies the name of the azure resource. Changing this forces a new resource to be created.

* `body` - (Required) Specifies the configuration of the resource. More information about the arguments in `body` can be found in the [Microsoft documentation](https://learn.microsoft.com/en-us/azure/templates/Microsoft.Maintenance/maintenanceConfigurations?pivots=deployment-language-terraform).

## Import

 ```shell
 # Azure resource can be imported using the resource id, e.g.
 terraform import azapi_resource.example /subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Maintenance/maintenanceConfigurations/{resourceName}
 
 # It also supports specifying API version by using the resource id with api-version as a query parameter, e.g.
 terraform import azapi_resource.example /subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Maintenance/maintenanceConfigurations/{resourceName}?api-version=2023-10-01-preview
 ```
