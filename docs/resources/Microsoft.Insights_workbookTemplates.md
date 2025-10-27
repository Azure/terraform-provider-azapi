---
subcategory: "Microsoft.Insights - Azure Monitor"
page_title: "workbookTemplates"
description: |-
  Manages a Application Insights Workbook Template.
---

# Microsoft.Insights/workbookTemplates - Application Insights Workbook Template

This article demonstrates how to use `azapi` provider to manage the Application Insights Workbook Template resource in Azure.



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

resource "azapi_resource" "workbookTemplate" {
  type      = "Microsoft.Insights/workbookTemplates@2020-11-20"
  parent_id = azapi_resource.resourceGroup.id
  name      = var.resource_name
  location  = var.location
  body = {
    properties = {
      galleries = [
        {
          category     = "workbook"
          name         = "test"
          order        = 0
          resourceType = "Azure Monitor"
          type         = "workbook"
        },
      ]
      priority = 0
      templateData = {
        "$schema" = "https://github.com/Microsoft/Application-Insights-Workbooks/blob/master/schema/workbook.json"
        items = [
          {
            content = {
              json = "## New workbook\n---\n\nWelcome to your new workbook."
            }
            name = "text - 2"
            type = 1
          },
        ]
        styleSettings = {
        }
        version = "Notebook/1.0"
      }
    }

  }
  schema_validation_enabled = false
  response_export_values    = ["*"]
}


```



## Arguments Reference

The following arguments are supported:

* `type` - (Required) The type of the resource. This should be set to `Microsoft.Insights/workbookTemplates@api-version`. The available api-versions for this resource are: [`2019-10-17-preview`, `2020-11-20`].

* `parent_id` - (Required) The ID of the azure resource in which this resource is created. The allowed values are:  
  `/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}`

* `name` - (Required) Specifies the name of the azure resource. Changing this forces a new resource to be created.

* `body` - (Required) Specifies the configuration of the resource. More information about the arguments in `body` can be found in the [Microsoft documentation](https://learn.microsoft.com/en-us/azure/templates/Microsoft.Insights/workbookTemplates?pivots=deployment-language-terraform).

For other arguments, please refer to the [azapi_resource](https://registry.terraform.io/providers/Azure/azapi/latest/docs/resources/resource) documentation.

## Import

 ```shell
 # Azure resource can be imported using the resource id, e.g.
 terraform import azapi_resource.example /subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Insights/workbookTemplates/{resourceName}
 
 # It also supports specifying API version by using the resource id with api-version as a query parameter, e.g.
 terraform import azapi_resource.example /subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Insights/workbookTemplates/{resourceName}?api-version=2020-11-20
 ```
