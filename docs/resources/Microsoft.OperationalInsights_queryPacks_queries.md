---
subcategory: "Microsoft.OperationalInsights - Azure Monitor"
page_title: "queryPacks/queries"
description: |-
  Manages a Log Analytics Query Pack Query.
---

# Microsoft.OperationalInsights/queryPacks/queries - Log Analytics Query Pack Query

This article demonstrates how to use `azapi` provider to manage the Log Analytics Query Pack Query resource in Azure.



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

resource "azapi_resource" "queryPack" {
  type      = "Microsoft.OperationalInsights/queryPacks@2019-09-01"
  parent_id = azapi_resource.resourceGroup.id
  name      = var.resource_name
  location  = var.location
  body = {
    properties = {
    }

  }
  schema_validation_enabled = false
  response_export_values    = ["*"]
}

resource "azapi_resource" "query" {
  type      = "Microsoft.OperationalInsights/queryPacks/queries@2019-09-01"
  parent_id = azapi_resource.queryPack.id
  name      = "aca50e92-d3e6-8f7d-1f70-2ec7adc1a926"
  body = {
    properties = {
      body        = "    let newExceptionsTimeRange = 1d;\n    let timeRangeToCheckBefore = 7d;\n    exceptions\n    | where timestamp < ago(timeRangeToCheckBefore)\n    | summarize count() by problemId\n    | join kind= rightanti (\n        exceptions\n        | where timestamp >= ago(newExceptionsTimeRange)\n        | extend stack = tostring(details[0].rawStack)\n        | summarize count(), dcount(user_AuthenticatedId), min(timestamp), max(timestamp), any(stack) by problemId\n    ) on problemId\n    | order by count_ desc\n"
      displayName = "Exceptions - New in the last 24 hours"
      related = {
      }
    }
  }
  schema_validation_enabled = false
  response_export_values    = ["*"]
}


```



## Arguments Reference

The following arguments are supported:

* `type` - (Required) The type of the resource. This should be set to `Microsoft.OperationalInsights/queryPacks/queries@api-version`. The available api-versions for this resource are: [`2019-09-01`, `2019-09-01-preview`, `2023-09-01`, `2025-02-01`].

* `parent_id` - (Required) The ID of the azure resource in which this resource is created. The allowed values are:  
  `/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.OperationalInsights/queryPacks/{resourceName}`

* `name` - (Required) Specifies the name of the azure resource. Changing this forces a new resource to be created.

* `body` - (Required) Specifies the configuration of the resource. More information about the arguments in `body` can be found in the [Microsoft documentation](https://learn.microsoft.com/en-us/azure/templates/Microsoft.OperationalInsights/queryPacks/queries?pivots=deployment-language-terraform).

For other arguments, please refer to the [azapi_resource](https://registry.terraform.io/providers/Azure/azapi/latest/docs/resources/resource) documentation.

## Import

 ```shell
 # Azure resource can be imported using the resource id, e.g.
 terraform import azapi_resource.example /subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.OperationalInsights/queryPacks/{resourceName}/queries/{resourceName}
 
 # It also supports specifying API version by using the resource id with api-version as a query parameter, e.g.
 terraform import azapi_resource.example /subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.OperationalInsights/queryPacks/{resourceName}/queries/{resourceName}?api-version=2025-02-01
 ```
