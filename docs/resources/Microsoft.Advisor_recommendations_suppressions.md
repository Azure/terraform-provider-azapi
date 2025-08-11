---
subcategory: "Microsoft.Advisor - Azure Advisor"
page_title: "recommendations/suppressions"
description: |-
  Manages a Specifies a suppression for an Azure Advisor recommendation.
---

# Microsoft.Advisor/recommendations/suppressions - Specifies a suppression for an Azure Advisor recommendation

This article demonstrates how to use `azapi` provider to manage the Specifies a suppression for an Azure Advisor recommendation resource in Azure.

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

variable "resource_name" {
  type    = string
  default = "acctest0001"
}

variable "location" {
  type    = string
  default = "westus"
}

variable "recommendation_id" {
  type = string
}

data "azapi_client_config" "current" {}

resource "azapi_resource" "suppression" {
  type      = "Microsoft.Advisor/recommendations/suppressions@2023-01-01"
  parent_id = "${data.azapi_client_config.current.subscription_resource_id}/providers/Microsoft.Advisor/recommendations/${var.recommendation_id}"
  name      = var.resource_name
  body = {
    properties = {
      suppressionId = ""
      ttl           = "00:30:00"
    }
  }
}

```



## Arguments Reference

The following arguments are supported:

* `type` - (Required) The type of the resource. This should be set to `Microsoft.Advisor/recommendations/suppressions@api-version`. The available api-versions for this resource are: [`2016-07-12-preview`, `2017-03-31`, `2017-04-19`, `2020-01-01`, `2022-09-01`, `2022-10-01`, `2023-01-01`, `2023-09-01-preview`, `2024-11-18-preview`, `2025-01-01`, `2025-05-01-preview`].

* `parent_id` - (Required) The ID of the azure resource in which this resource is created. The allowed values are:  
  `{any azure resource id}`

* `name` - (Required) Specifies the name of the azure resource. Changing this forces a new resource to be created.

* `body` - (Required) Specifies the configuration of the resource. More information about the arguments in `body` can be found in the [Microsoft documentation](https://learn.microsoft.com/en-us/azure/templates/Microsoft.Advisor/recommendations/suppressions?pivots=deployment-language-terraform).

For other arguments, please refer to the [azapi_resource](https://registry.terraform.io/providers/Azure/azapi/latest/docs/resources/resource) documentation.

## Import

 ```shell
 # Azure resource can be imported using the resource id, e.g.
 terraform import azapi_resource.example {any azure resource id}/suppressions/{resourceName}
 
 # It also supports specifying API version by using the resource id with api-version as a query parameter, e.g.
 terraform import azapi_resource.example {any azure resource id}/suppressions/{resourceName}?api-version=2025-05-01-preview
 ```
