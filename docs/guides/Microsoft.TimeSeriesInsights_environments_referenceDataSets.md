---
subcategory: "Microsoft.TimeSeriesInsights - Azure Time Series Insights"
page_title: "environments/referenceDataSets"
description: |-
  Manages a Time Series Insights Environments Reference Data Sets.
---

# Microsoft.TimeSeriesInsights/environments/referenceDataSets - Time Series Insights Environments Reference Data Sets

This article demonstrates how to use `azapi` provider to manage the Time Series Insights Environments Reference Data Sets resource in Azure.

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

resource "azapi_resource" "environment" {
  type      = "Microsoft.TimeSeriesInsights/environments@2020-05-15"
  parent_id = azapi_resource.resourceGroup.id
  name      = var.resource_name
  location  = var.location
  body = {
    kind = "Gen1"
    properties = {
      dataRetentionTime            = "P30D"
      storageLimitExceededBehavior = "PurgeOldData"
    }
    sku = {
      capacity = 1
      name     = "S1"
    }
  }
  schema_validation_enabled = false
  response_export_values    = ["*"]
}

resource "azapi_resource" "referenceDataSet" {
  type      = "Microsoft.TimeSeriesInsights/environments/referenceDataSets@2020-05-15"
  parent_id = azapi_resource.environment.id
  name      = var.resource_name
  location  = var.location
  body = {
    properties = {
      dataStringComparisonBehavior = "Ordinal"
      keyProperties = [
        {
          name = "keyProperty1"
          type = "String"
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

* `type` - (Required) The type of the resource. This should be set to `Microsoft.TimeSeriesInsights/environments/referenceDataSets@api-version`. The available api-versions for this resource are: [`2017-02-28-preview`, `2017-11-15`, `2018-08-15-preview`, `2020-05-15`, `2021-03-31-preview`, `2021-06-30-preview`].

* `parent_id` - (Required) The ID of the azure resource in which this resource is created. The allowed values are:  
  `/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.TimeSeriesInsights/environments/{resourceName}`

* `name` - (Required) Specifies the name of the azure resource. Changing this forces a new resource to be created.

* `body` - (Required) Specifies the configuration of the resource. More information about the arguments in `body` can be found in the [Microsoft documentation](https://learn.microsoft.com/en-us/azure/templates/Microsoft.TimeSeriesInsights/environments/referenceDataSets?pivots=deployment-language-terraform).

## Import

 ```shell
 # Azure resource can be imported using the resource id, e.g.
 terraform import azapi_resource.example /subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.TimeSeriesInsights/environments/{resourceName}/referenceDataSets/{resourceName}
 
 # It also supports specifying API version by using the resource id with api-version as a query parameter, e.g.
 terraform import azapi_resource.example /subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.TimeSeriesInsights/environments/{resourceName}/referenceDataSets/{resourceName}?api-version=2021-06-30-preview
 ```
