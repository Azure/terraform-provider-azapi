---
subcategory: "Microsoft.Network - Various Networking Services"
page_title: "privateDnsZones/AAAA"
description: |-
  Manages a Private DNS AAAA Record.
---

# Microsoft.Network/privateDnsZones/AAAA - Private DNS AAAA Record

This article demonstrates how to use `azapi` provider to manage the Private DNS AAAA Record resource in Azure.

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
  name                      = "${var.resource_name}.com"
  location                  = "global"
  schema_validation_enabled = false
  response_export_values    = ["*"]
}

resource "azapi_resource" "AAAA" {
  type      = "Microsoft.Network/privateDnsZones/AAAA@2018-09-01"
  parent_id = azapi_resource.privateDnsZone.id
  name      = var.resource_name
  body = {
    properties = {
      aaaaRecords = [
        {
          ipv6Address = "fd5d:70bc:930e:d008:0000:0000:0000:7334"
        },
        {
          ipv6Address = "fd5d:70bc:930e:d008::7335"
        },
      ]
      metadata = {
      }
      ttl = 300
    }
  }
  schema_validation_enabled = false
  response_export_values    = ["*"]
}


```



## Arguments Reference

The following arguments are supported:

* `type` - (Required) The type of the resource. This should be set to `Microsoft.Network/privateDnsZones/AAAA@api-version`. The available api-versions for this resource are: [`2018-09-01`, `2020-01-01`, `2020-06-01`, `2024-06-01`].

* `parent_id` - (Required) The ID of the azure resource in which this resource is created. The allowed values are:  
  `/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Network/privateDnsZones/{resourceName}`

* `name` - (Required) Specifies the name of the azure resource. Changing this forces a new resource to be created.

* `body` - (Required) Specifies the configuration of the resource. More information about the arguments in `body` can be found in the [Microsoft documentation](https://learn.microsoft.com/en-us/azure/templates/Microsoft.Network/privateDnsZones/AAAA?pivots=deployment-language-terraform).

For other arguments, please refer to the [azapi_resource](https://registry.terraform.io/providers/Azure/azapi/latest/docs/resources/resource) documentation.

## Import

 ```shell
 # Azure resource can be imported using the resource id, e.g.
 terraform import azapi_resource.example /subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Network/privateDnsZones/{resourceName}/AAAA/{resourceName}
 
 # It also supports specifying API version by using the resource id with api-version as a query parameter, e.g.
 terraform import azapi_resource.example /subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Network/privateDnsZones/{resourceName}/AAAA/{resourceName}?api-version=2024-06-01
 ```
