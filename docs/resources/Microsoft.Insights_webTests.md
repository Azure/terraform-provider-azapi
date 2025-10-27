---
subcategory: "Microsoft.Insights - Azure Monitor"
page_title: "webTests"
description: |-
  Manages a Application Insights Standard WebTest.
---

# Microsoft.Insights/webTests - Application Insights Standard WebTest

This article demonstrates how to use `azapi` provider to manage the Application Insights Standard WebTest resource in Azure.



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

resource "azapi_resource" "webTest" {
  type      = "Microsoft.Insights/webTests@2022-06-15"
  parent_id = azapi_resource.resourceGroup.id
  name      = var.resource_name
  location  = var.location
  body = {
    kind = "standard"
    properties = {
      Description = ""
      Enabled     = false
      Frequency   = 300
      Kind        = "standard"
      Locations = [
        {
          Id = "us-tx-sn1-azr"
        },
      ]
      Name = var.resource_name
      Request = {
        FollowRedirects = false
        Headers = [
          {
            key   = "x-header"
            value = "testheader"
          },
          {
            key   = "x-header-2"
            value = "testheader2"
          },
        ]
        HttpVerb               = "GET"
        ParseDependentRequests = false
        RequestUrl             = "http://microsoft.com"
      }
      RetryEnabled       = false
      SyntheticMonitorId = var.resource_name
      Timeout            = 30
      ValidationRules = {
        ExpectedHttpStatusCode = 200
        SSLCheck               = false
      }
    }
    tags = {
      "hidden-link:${azapi_resource.component.id}" = "Resource"
    }
  }
  schema_validation_enabled = false
  response_export_values    = ["*"]
}


```



## Arguments Reference

The following arguments are supported:

* `type` - (Required) The type of the resource. This should be set to `Microsoft.Insights/webTests@api-version`. The available api-versions for this resource are: [`2015-05-01`, `2018-05-01-preview`, `2020-10-05-preview`, `2022-06-15`].

* `parent_id` - (Required) The ID of the azure resource in which this resource is created. The allowed values are:  
  `/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}`

* `name` - (Required) Specifies the name of the azure resource. Changing this forces a new resource to be created.

* `body` - (Required) Specifies the configuration of the resource. More information about the arguments in `body` can be found in the [Microsoft documentation](https://learn.microsoft.com/en-us/azure/templates/Microsoft.Insights/webTests?pivots=deployment-language-terraform).

For other arguments, please refer to the [azapi_resource](https://registry.terraform.io/providers/Azure/azapi/latest/docs/resources/resource) documentation.

## Import

 ```shell
 # Azure resource can be imported using the resource id, e.g.
 terraform import azapi_resource.example /subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Insights/webTests/{resourceName}
 
 # It also supports specifying API version by using the resource id with api-version as a query parameter, e.g.
 terraform import azapi_resource.example /subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Insights/webTests/{resourceName}?api-version=2022-06-15
 ```
