---
subcategory: "Microsoft.Network - Various Networking Services"
page_title: "dnsZones/CAA"
description: |-
  Manages a DNS CAA Record.
---

# Microsoft.Network/dnsZones/CAA - DNS CAA Record

This article demonstrates how to use `azapi` provider to manage the DNS CAA Record resource in Azure.

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

resource "azapi_resource" "dnsZone" {
  type                      = "Microsoft.Network/dnsZones@2018-05-01"
  parent_id                 = azapi_resource.resourceGroup.id
  name                      = "${var.resource_name}.com"
  location                  = "global"
  schema_validation_enabled = false
  response_export_values    = ["*"]
}

resource "azapi_resource" "CAA" {
  type      = "Microsoft.Network/dnsZones/CAA@2018-05-01"
  parent_id = azapi_resource.dnsZone.id
  name      = var.resource_name
  body = {
    properties = {
      TTL = 300
      caaRecords = [
        {
          flags = 1
          tag   = "issuewild"
          value = ";"
        },
        {
          flags = 0
          tag   = "iodef"
          value = "mailto:terraform@nonexist.tld"
        },
        {
          flags = 0
          tag   = "issue"
          value = "example.com"
        },
        {
          flags = 0
          tag   = "issue"
          value = "example.net"
        },
      ]
      metadata = {
      }
    }
  }
  schema_validation_enabled = false
  response_export_values    = ["*"]
}


```



## Arguments Reference

The following arguments are supported:

* `type` - (Required) The type of the resource. This should be set to `Microsoft.Network/dnsZones/CAA@api-version`. The available api-versions for this resource are: [`2017-09-01`, `2017-10-01`, `2018-03-01-preview`, `2018-05-01`, `2023-07-01-preview`].

* `parent_id` - (Required) The ID of the azure resource in which this resource is created. The allowed values are:  
  `/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Network/dnsZones/{resourceName}`

* `name` - (Required) Specifies the name of the azure resource. Changing this forces a new resource to be created.

* `body` - (Required) Specifies the configuration of the resource. More information about the arguments in `body` can be found in the [Microsoft documentation](https://learn.microsoft.com/en-us/azure/templates/Microsoft.Network/dnsZones/CAA?pivots=deployment-language-terraform).

For other arguments, please refer to the [azapi_resource](https://registry.terraform.io/providers/Azure/azapi/latest/docs/resources/resource) documentation.

## Import

 ```shell
 # Azure resource can be imported using the resource id, e.g.
 terraform import azapi_resource.example /subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Network/dnsZones/{resourceName}/CAA/{resourceName}
 
 # It also supports specifying API version by using the resource id with api-version as a query parameter, e.g.
 terraform import azapi_resource.example /subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Network/dnsZones/{resourceName}/CAA/{resourceName}?api-version=2023-07-01-preview
 ```
