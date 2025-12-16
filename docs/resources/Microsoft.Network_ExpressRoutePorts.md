---
subcategory: "Microsoft.Network - Various Networking Services"
page_title: "ExpressRoutePorts"
description: |-
  Manages a Express Route Port.
---

# Microsoft.Network/ExpressRoutePorts - Express Route Port

This article demonstrates how to use `azapi` provider to manage the Express Route Port resource in Azure.



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

resource "azapi_resource" "ExpressRoutePort" {
  type      = "Microsoft.Network/ExpressRoutePorts@2022-07-01"
  parent_id = azapi_resource.resourceGroup.id
  name      = var.resource_name
  location  = var.location
  body = {
    properties = {
      bandwidthInGbps = 10
      billingType     = "MeteredData"
      encapsulation   = "Dot1Q"
      peeringLocation = "Airtel-Chennai2-CLS"
    }
    tags = {
      ENV = "Test"
    }
  }
  schema_validation_enabled = false
  response_export_values    = ["*"]
}


```



## Arguments Reference

The following arguments are supported:

* `type` - (Required) The type of the resource. This should be set to `Microsoft.Network/ExpressRoutePorts@api-version`. The available api-versions for this resource are: [`2018-08-01`, `2018-10-01`, `2018-11-01`, `2018-12-01`, `2019-02-01`, `2019-04-01`, `2019-06-01`, `2019-07-01`, `2019-08-01`, `2019-09-01`, `2019-11-01`, `2019-12-01`, `2020-03-01`, `2020-04-01`, `2020-05-01`, `2020-06-01`, `2020-07-01`, `2020-08-01`, `2020-11-01`, `2021-02-01`, `2021-03-01`, `2021-05-01`, `2021-08-01`, `2022-01-01`, `2022-05-01`, `2022-07-01`, `2022-09-01`, `2022-11-01`, `2023-02-01`, `2023-04-01`, `2023-05-01`, `2023-06-01`, `2023-09-01`, `2023-11-01`, `2024-01-01`, `2024-03-01`, `2024-05-01`, `2024-07-01`, `2024-10-01`, `2025-01-01`, `2025-03-01`].

* `parent_id` - (Required) The ID of the azure resource in which this resource is created. The allowed values are:  
  `/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}`

* `name` - (Required) Specifies the name of the azure resource. Changing this forces a new resource to be created.

* `body` - (Required) Specifies the configuration of the resource. More information about the arguments in `body` can be found in the [Microsoft documentation](https://learn.microsoft.com/en-us/azure/templates/Microsoft.Network/ExpressRoutePorts?pivots=deployment-language-terraform).

For other arguments, please refer to the [azapi_resource](https://registry.terraform.io/providers/Azure/azapi/latest/docs/resources/resource) documentation.

## Import

 ```shell
 # Azure resource can be imported using the resource id, e.g.
 terraform import azapi_resource.example /subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Network/ExpressRoutePorts/{resourceName}
 
 # It also supports specifying API version by using the resource id with api-version as a query parameter, e.g.
 terraform import azapi_resource.example /subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Network/ExpressRoutePorts/{resourceName}?api-version=2025-03-01
 ```
