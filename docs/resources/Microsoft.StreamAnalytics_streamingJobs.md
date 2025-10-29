---
subcategory: "Microsoft.StreamAnalytics - Azure Stream Analytics"
page_title: "streamingJobs"
description: |-
  Manages a Stream Analytics Job.
---

# Microsoft.StreamAnalytics/streamingJobs - Stream Analytics Job

This article demonstrates how to use `azapi` provider to manage the Stream Analytics Job resource in Azure.



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

resource "azapi_resource" "streamingJob" {
  type      = "Microsoft.StreamAnalytics/streamingJobs@2020-03-01"
  parent_id = azapi_resource.resourceGroup.id
  name      = var.resource_name
  location  = var.location
  body = {
    properties = {
      cluster = {
      }
      compatibilityLevel                 = "1.0"
      contentStoragePolicy               = "SystemAccount"
      dataLocale                         = "en-GB"
      eventsLateArrivalMaxDelayInSeconds = 60
      eventsOutOfOrderMaxDelayInSeconds  = 50
      eventsOutOfOrderPolicy             = "Adjust"
      jobType                            = "Cloud"
      outputErrorPolicy                  = "Drop"
      sku = {
        name = "Standard"
      }
      transformation = {
        name = "main"
        properties = {
          query          = "   SELECT *\n   INTO [YourOutputAlias]\n   FROM [YourInputAlias]\n"
          streamingUnits = 3
        }
      }
    }
  }
  schema_validation_enabled = false
  response_export_values    = ["*"]
}


```



## Arguments Reference

The following arguments are supported:

* `type` - (Required) The type of the resource. This should be set to `Microsoft.StreamAnalytics/streamingJobs@api-version`. The available api-versions for this resource are: [`2016-03-01`, `2017-04-01-preview`, `2020-03-01`, `2021-10-01-preview`].

* `parent_id` - (Required) The ID of the azure resource in which this resource is created. The allowed values are:  
  `/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}`

* `name` - (Required) Specifies the name of the azure resource. Changing this forces a new resource to be created.

* `body` - (Required) Specifies the configuration of the resource. More information about the arguments in `body` can be found in the [Microsoft documentation](https://learn.microsoft.com/en-us/azure/templates/Microsoft.StreamAnalytics/streamingJobs?pivots=deployment-language-terraform).

For other arguments, please refer to the [azapi_resource](https://registry.terraform.io/providers/Azure/azapi/latest/docs/resources/resource) documentation.

## Import

 ```shell
 # Azure resource can be imported using the resource id, e.g.
 terraform import azapi_resource.example /subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.StreamAnalytics/streamingJobs/{resourceName}
 
 # It also supports specifying API version by using the resource id with api-version as a query parameter, e.g.
 terraform import azapi_resource.example /subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.StreamAnalytics/streamingJobs/{resourceName}?api-version=2021-10-01-preview
 ```
