---
subcategory: "Microsoft.Cdn - Content Delivery Network"
page_title: "profiles/originGroups"
description: |-
  Manages a Front Door (standard/premium) Origin Group.
---

# Microsoft.Cdn/profiles/originGroups - Front Door (standard/premium) Origin Group

This article demonstrates how to use `azapi` provider to manage the Front Door (standard/premium) Origin Group resource in Azure.

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

resource "azapi_resource" "profile" {
  type      = "Microsoft.Cdn/profiles@2021-06-01"
  parent_id = azapi_resource.resourceGroup.id
  name      = var.resource_name
  location  = "global"
  body = {
    properties = {
      originResponseTimeoutSeconds = 120
    }
    sku = {
      name = "Standard_AzureFrontDoor"
    }
  }
  schema_validation_enabled = false
  response_export_values    = ["*"]
}

resource "azapi_resource" "originGroup" {
  type      = "Microsoft.Cdn/profiles/originGroups@2021-06-01"
  parent_id = azapi_resource.profile.id
  name      = var.resource_name
  body = {
    properties = {
      loadBalancingSettings = {
        additionalLatencyInMilliseconds = 0
        sampleSize                      = 16
        successfulSamplesRequired       = 3
      }
      sessionAffinityState                                  = "Enabled"
      trafficRestorationTimeToHealedOrNewEndpointsInMinutes = 10
    }
  }
  schema_validation_enabled = false
  response_export_values    = ["*"]
}


```



## Arguments Reference

The following arguments are supported:

* `type` - (Required) The type of the resource. This should be set to `Microsoft.Cdn/profiles/originGroups@api-version`. The available api-versions for this resource are: [`2020-09-01`, `2021-06-01`, `2022-05-01-preview`, `2022-11-01-preview`, `2023-05-01`, `2023-07-01-preview`, `2024-02-01`, `2024-05-01-preview`, `2024-06-01-preview`, `2024-09-01`, `2025-01-01-preview`, `2025-04-15`].

* `parent_id` - (Required) The ID of the azure resource in which this resource is created. The allowed values are:  
  `/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Cdn/profiles/{resourceName}`

* `name` - (Required) Specifies the name of the azure resource. Changing this forces a new resource to be created.

* `body` - (Required) Specifies the configuration of the resource. More information about the arguments in `body` can be found in the [Microsoft documentation](https://learn.microsoft.com/en-us/azure/templates/Microsoft.Cdn/profiles/originGroups?pivots=deployment-language-terraform).

For other arguments, please refer to the [azapi_resource](https://registry.terraform.io/providers/Azure/azapi/latest/docs/resources/resource) documentation.

## Import

 ```shell
 # Azure resource can be imported using the resource id, e.g.
 terraform import azapi_resource.example /subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Cdn/profiles/{resourceName}/originGroups/{resourceName}
 
 # It also supports specifying API version by using the resource id with api-version as a query parameter, e.g.
 terraform import azapi_resource.example /subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Cdn/profiles/{resourceName}/originGroups/{resourceName}?api-version=2025-04-15
 ```
