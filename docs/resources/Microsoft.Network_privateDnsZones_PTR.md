---
subcategory: "Microsoft.Network - Various Networking Services"
page_title: "privateDnsZones/PTR"
description: |-
  Manages a Private DNS PTR Record.
---

# Microsoft.Network/privateDnsZones/PTR - Private DNS PTR Record

This article demonstrates how to use `azapi` provider to manage the Private DNS PTR Record resource in Azure.



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

resource "azapi_resource" "privateDnsZone" {
  type                      = "Microsoft.Network/privateDnsZones@2018-09-01"
  parent_id                 = azapi_resource.resourceGroup.id
  name                      = "230630033756174960.0.10.in-addr.arpa"
  location                  = "global"
  schema_validation_enabled = false
  response_export_values    = ["*"]
}

resource "azapi_resource" "PTR" {
  type      = "Microsoft.Network/privateDnsZones/PTR@2018-09-01"
  parent_id = azapi_resource.privateDnsZone.id
  name      = var.resource_name
  body = {
    properties = {
      metadata = {
      }
      ptrRecords = [
        {
          ptrdname = "test2.contoso.com"
        },
        {
          ptrdname = "test.contoso.com"
        },
      ]
      ttl = 300
    }
  }
  schema_validation_enabled = false
  response_export_values    = ["*"]
}


```



## Arguments Reference

The following arguments are supported:

* `type` - (Required) The type of the resource. This should be set to `Microsoft.Network/privateDnsZones/PTR@api-version`. The available api-versions for this resource are: [`2018-09-01`, `2020-01-01`, `2020-06-01`, `2024-06-01`].

* `parent_id` - (Required) The ID of the azure resource in which this resource is created. The allowed values are:  
  `/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Network/privateDnsZones/{resourceName}`

* `name` - (Required) Specifies the name of the azure resource. Changing this forces a new resource to be created.

* `body` - (Required) Specifies the configuration of the resource. More information about the arguments in `body` can be found in the [Microsoft documentation](https://learn.microsoft.com/en-us/azure/templates/Microsoft.Network/privateDnsZones/PTR?pivots=deployment-language-terraform).

For other arguments, please refer to the [azapi_resource](https://registry.terraform.io/providers/Azure/azapi/latest/docs/resources/resource) documentation.

## Import

 ```shell
 # Azure resource can be imported using the resource id, e.g.
 terraform import azapi_resource.example /subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Network/privateDnsZones/{resourceName}/PTR/{resourceName}
 
 # It also supports specifying API version by using the resource id with api-version as a query parameter, e.g.
 terraform import azapi_resource.example /subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Network/privateDnsZones/{resourceName}/PTR/{resourceName}?api-version=2024-06-01
 ```
