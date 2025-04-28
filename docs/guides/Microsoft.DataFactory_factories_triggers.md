---
subcategory: "Microsoft.DataFactory - Data Factory"
page_title: "factories/triggers"
description: |-
  Manages a Trigger inside an Azure Data Factory.
---

# Microsoft.DataFactory/factories/triggers - Trigger inside an Azure Data Factory

This article demonstrates how to use `azapi` provider to manage the Trigger inside an Azure Data Factory resource in Azure.

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

resource "azapi_resource" "factory" {
  type      = "Microsoft.DataFactory/factories@2018-06-01"
  parent_id = azapi_resource.resourceGroup.id
  name      = var.resource_name
  location  = var.location
  body = {
    properties = {
      publicNetworkAccess = "Enabled"
      repoConfiguration   = null
    }
  }
  schema_validation_enabled = false
  response_export_values    = ["*"]
}

resource "azapi_resource" "pipeline" {
  type      = "Microsoft.DataFactory/factories/pipelines@2018-06-01"
  parent_id = azapi_resource.factory.id
  name      = var.resource_name
  body = {
    properties = {
      annotations = [
      ]
      description = ""
      parameters = {
        test = {
          defaultValue = "testparameter"
          type         = "String"
        }
      }
      variables = {
      }
    }
  }
  schema_validation_enabled = false
  response_export_values    = ["*"]
}

resource "azapi_resource" "trigger" {
  type      = "Microsoft.DataFactory/factories/triggers@2018-06-01"
  parent_id = azapi_resource.factory.id
  name      = var.resource_name
  body = {
    properties = {
      description = ""
      pipeline = {
        parameters = {
        }
        pipelineReference = {
          referenceName = azapi_resource.pipeline.name
          type          = "PipelineReference"
        }
      }
      type = "TumblingWindowTrigger"
      typeProperties = {
        frequency      = "Minute"
        interval       = 15
        maxConcurrency = 50
        startTime      = "2022-09-21T00:00:00Z"
      }
    }
  }
  schema_validation_enabled = false
  response_export_values    = ["*"]
}


```



## Arguments Reference

The following arguments are supported:

* `type` - (Required) The type of the resource. This should be set to `Microsoft.DataFactory/factories/triggers@api-version`. The available api-versions for this resource are: [`2017-09-01-preview`, `2018-06-01`].

* `parent_id` - (Required) The ID of the azure resource in which this resource is created. The allowed values are:  
  `/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.DataFactory/factories/{resourceName}`

* `name` - (Required) Specifies the name of the azure resource. Changing this forces a new resource to be created.

* `body` - (Required) Specifies the configuration of the resource. More information about the arguments in `body` can be found in the [Microsoft documentation](https://learn.microsoft.com/en-us/azure/templates/Microsoft.DataFactory/factories/triggers?pivots=deployment-language-terraform).

## Import

 ```shell
 # Azure resource can be imported using the resource id, e.g.
 terraform import azapi_resource.example /subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.DataFactory/factories/{resourceName}/triggers/{resourceName}
 
 # It also supports specifying API version by using the resource id with api-version as a query parameter, e.g.
 terraform import azapi_resource.example /subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.DataFactory/factories/{resourceName}/triggers/{resourceName}?api-version=2018-06-01
 ```
