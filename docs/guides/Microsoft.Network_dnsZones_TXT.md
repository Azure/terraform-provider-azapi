---
subcategory: "Microsoft.Network - Various Networking Services"
page_title: "dnsZones/TXT"
description: |-
  Manages a DNS TXT Record.
---

# Microsoft.Network/dnsZones/TXT - DNS TXT Record

This article demonstrates how to use `azapi` provider to manage the DNS TXT Record resource in Azure.

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

resource "azapi_resource" "TXT" {
  type      = "Microsoft.Network/dnsZones/TXT@2018-05-01"
  parent_id = azapi_resource.dnsZone.id
  name      = var.resource_name
  body = {
    properties = {
      TTL = 300
      TXTRecords = [
        {
          value = [
            "Quick brown fox",
          ]
        },
        {
          value = [
            "A long text......A long text......A long text......A long text......A long text......A long text......A long text......A long text......A long text......A long text......A long text......A long text......A long text......A long text......A long text.....",
            ".A long text......A long text......A long text......A long text......A long text......A long text......A long text......A long text......A long text......A long text......A long text......A long text......A long text......A long text......A long text....",
            "..A long text......A long text......A long text......A long text......A long text......A long text......A long text......A long text......A long text......A long text......A long text......A long text......",
          ]
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

* `type` - (Required) The type of the resource. This should be set to `Microsoft.Network/dnsZones/TXT@api-version`. The available api-versions for this resource are: [`2015-05-04-preview`, `2016-04-01`, `2017-09-01`, `2017-10-01`, `2018-03-01-preview`, `2018-05-01`, `2023-07-01-preview`].

* `parent_id` - (Required) The ID of the azure resource in which this resource is created. The allowed values are:  
  `/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Network/dnsZones/{resourceName}`

* `name` - (Required) Specifies the name of the azure resource. Changing this forces a new resource to be created.

* `body` - (Required) Specifies the configuration of the resource. More information about the arguments in `body` can be found in the [Microsoft documentation](https://learn.microsoft.com/en-us/azure/templates/Microsoft.Network/dnsZones/TXT?pivots=deployment-language-terraform).

## Import

 ```shell
 # Azure resource can be imported using the resource id, e.g.
 terraform import azapi_resource.example /subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Network/dnsZones/{resourceName}/TXT/{resourceName}
 
 # It also supports specifying API version by using the resource id with api-version as a query parameter, e.g.
 terraform import azapi_resource.example /subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Network/dnsZones/{resourceName}/TXT/{resourceName}?api-version=2023-07-01-preview
 ```
