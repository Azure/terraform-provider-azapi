---
subcategory: "Microsoft.Insights - Azure Monitor"
page_title: "workbooks"
description: |-
  Manages a Azure Workbook.
---

# Microsoft.Insights/workbooks - Azure Workbook

This article demonstrates how to use `azapi` provider to manage the Azure Workbook resource in Azure.

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

resource "azapi_resource" "workbook" {
  type      = "Microsoft.Insights/workbooks@2022-04-01"
  parent_id = azapi_resource.resourceGroup.id
  name      = "be1ad266-d329-4454-b693-8287e4d3b35d"
  location  = var.location
  body = {
    kind = "shared"
    properties = {
      category       = "workbook"
      displayName    = "acctest-amw-230630032616547405"
      serializedData = "{\"fallbackResourceIds\":[\"Azure Monitor\"],\"isLocked\":false,\"items\":[{\"content\":{\"json\":\"Test2022\"},\"name\":\"text - 0\",\"type\":1}],\"version\":\"Notebook/1.0\"}"
      sourceId       = "azure monitor"
    }
  }
  schema_validation_enabled = false
  response_export_values    = ["*"]
}


```



## Arguments Reference

The following arguments are supported:

* `type` - (Required) The type of the resource. This should be set to `Microsoft.Insights/workbooks@api-version`. The available api-versions for this resource are: [`2015-05-01`, `2018-06-17-preview`, `2020-10-20`, `2021-03-08`, `2021-08-01`, `2022-04-01`, `2023-06-01`].

* `parent_id` - (Required) The ID of the azure resource in which this resource is created. The allowed values are:  
  `/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}`

* `name` - (Required) Specifies the name of the azure resource. Changing this forces a new resource to be created.

* `body` - (Required) Specifies the configuration of the resource. More information about the arguments in `body` can be found in the [Microsoft documentation](https://learn.microsoft.com/en-us/azure/templates/Microsoft.Insights/workbooks?pivots=deployment-language-terraform).

For other arguments, please refer to the [azapi_resource](https://registry.terraform.io/providers/Azure/azapi/latest/docs/resources/resource) documentation.

## Import

 ```shell
 # Azure resource can be imported using the resource id, e.g.
 terraform import azapi_resource.example /subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Insights/workbooks/{resourceName}
 
 # It also supports specifying API version by using the resource id with api-version as a query parameter, e.g.
 terraform import azapi_resource.example /subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Insights/workbooks/{resourceName}?api-version=2023-06-01
 ```
