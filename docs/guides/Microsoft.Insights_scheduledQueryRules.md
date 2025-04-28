---
subcategory: "Microsoft.Insights - Azure Monitor"
page_title: "scheduledQueryRules"
description: |-
  Manages a AlertingAction Scheduled Query Rules resource within Azure Monitor.
---

# Microsoft.Insights/scheduledQueryRules - AlertingAction Scheduled Query Rules resource within Azure Monitor

This article demonstrates how to use `azapi` provider to manage the AlertingAction Scheduled Query Rules resource within Azure Monitor resource in Azure.

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

resource "azapi_resource" "component" {
  type      = "Microsoft.Insights/components@2020-02-02"
  parent_id = azapi_resource.resourceGroup.id
  name      = var.resource_name
  location  = var.location
  body = {
    kind = "web"
    properties = {
      Application_Type                = "web"
      DisableIpMasking                = false
      DisableLocalAuth                = false
      ForceCustomerStorageForProfiler = false
      RetentionInDays                 = 90
      SamplingPercentage              = 100
      publicNetworkAccessForIngestion = "Enabled"
      publicNetworkAccessForQuery     = "Enabled"
    }
  }
  schema_validation_enabled = false
  response_export_values    = ["*"]
}

resource "azapi_resource" "scheduledQueryRule" {
  type      = "Microsoft.Insights/scheduledQueryRules@2021-08-01"
  parent_id = azapi_resource.resourceGroup.id
  name      = var.resource_name
  location  = var.location
  body = {
    kind = "LogAlert"
    properties = {
      autoMitigate                          = false
      checkWorkspaceAlertsStorageConfigured = false
      criteria = {
        allOf = [
          {
            dimensions      = null
            operator        = "Equal"
            query           = " requests\n| summarize CountByCountry=count() by client_CountryOrRegion\n"
            threshold       = 5
            timeAggregation = "Count"
          },
        ]
      }
      enabled             = true
      evaluationFrequency = "PT5M"
      scopes = [
        azapi_resource.component.id,
      ]
      severity            = 3
      skipQueryValidation = false
      targetResourceTypes = null
      windowSize          = "PT5M"
    }

  }
  schema_validation_enabled = false
  response_export_values    = ["*"]
}


```



## Arguments Reference

The following arguments are supported:

* `type` - (Required) The type of the resource. This should be set to `Microsoft.Insights/scheduledQueryRules@api-version`. The available api-versions for this resource are: [`2018-04-16`, `2020-05-01-preview`, `2021-02-01-preview`, `2021-08-01`, `2022-06-15`, `2022-08-01-preview`, `2023-03-15-preview`, `2023-12-01`, `2024-01-01-preview`].

* `parent_id` - (Required) The ID of the azure resource in which this resource is created. The allowed values are:  
  `/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}`

* `name` - (Required) Specifies the name of the azure resource. Changing this forces a new resource to be created.

* `body` - (Required) Specifies the configuration of the resource. More information about the arguments in `body` can be found in the [Microsoft documentation](https://learn.microsoft.com/en-us/azure/templates/Microsoft.Insights/scheduledQueryRules?pivots=deployment-language-terraform).

## Import

 ```shell
 # Azure resource can be imported using the resource id, e.g.
 terraform import azapi_resource.example /subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Insights/scheduledQueryRules/{resourceName}
 
 # It also supports specifying API version by using the resource id with api-version as a query parameter, e.g.
 terraform import azapi_resource.example /subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Insights/scheduledQueryRules/{resourceName}?api-version=2024-01-01-preview
 ```
