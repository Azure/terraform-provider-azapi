---
subcategory: "Microsoft.ElasticSan - Elastic SAN"
page_title: "elasticSans/volumeGroups/volumes"
description: |-
  Manages a Elastic SAN Volume resource.
---

# Microsoft.ElasticSan/elasticSans/volumeGroups/volumes - Elastic SAN Volume resource

This article demonstrates how to use `azapi` provider to manage the Elastic SAN Volume resource resource in Azure.



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
  default = "westus"
}

resource "azapi_resource" "resourceGroup" {
  type     = "Microsoft.Resources/resourceGroups@2020-06-01"
  name     = var.resource_name
  location = var.location
}

resource "azapi_resource" "elasticSan" {
  type      = "Microsoft.ElasticSan/elasticSans@2023-01-01"
  parent_id = azapi_resource.resourceGroup.id
  name      = "${var.resource_name}-es"
  location  = var.location
  body = {
    properties = {
      baseSizeTiB             = 1
      extendedCapacitySizeTiB = 0
      sku = {
        name = "Premium_LRS"
        tier = "Premium"
      }
    }
  }
}

resource "azapi_resource" "volumeGroup" {
  type      = "Microsoft.ElasticSan/elasticSans/volumeGroups@2023-01-01"
  parent_id = azapi_resource.elasticSan.id
  name      = "${var.resource_name}-vg"
  body = {
    properties = {
      encryption = "EncryptionAtRestWithPlatformKey"
      networkAcls = {
        virtualNetworkRules = []
      }
      protocolType = "Iscsi"
    }
  }
}

resource "azapi_resource" "volume" {
  type      = "Microsoft.ElasticSan/elasticSans/volumeGroups/volumes@2023-01-01"
  parent_id = azapi_resource.volumeGroup.id
  name      = "${var.resource_name}-v"
  body = {
    properties = {
      sizeGiB = 1
    }
  }
}

```



## Arguments Reference

The following arguments are supported:

* `type` - (Required) The type of the resource. This should be set to `Microsoft.ElasticSan/elasticSans/volumeGroups/volumes@api-version`. The available api-versions for this resource are: [`2021-11-20-preview`, `2022-12-01-preview`, `2023-01-01`, `2024-05-01`, `2024-06-01-preview`, `2024-07-01-preview`, `2025-09-01`].

* `parent_id` - (Required) The ID of the azure resource in which this resource is created. The allowed values are:  
  `/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.ElasticSan/elasticSans/{resourceName}/volumeGroups/{resourceName}`

* `name` - (Required) Specifies the name of the azure resource. Changing this forces a new resource to be created.

* `body` - (Required) Specifies the configuration of the resource. More information about the arguments in `body` can be found in the [Microsoft documentation](https://learn.microsoft.com/en-us/azure/templates/Microsoft.ElasticSan/elasticSans/volumeGroups/volumes?pivots=deployment-language-terraform).

For other arguments, please refer to the [azapi_resource](https://registry.terraform.io/providers/Azure/azapi/latest/docs/resources/resource) documentation.

## Import

 ```shell
 # Azure resource can be imported using the resource id, e.g.
 terraform import azapi_resource.example /subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.ElasticSan/elasticSans/{resourceName}/volumeGroups/{resourceName}/volumes/{resourceName}
 
 # It also supports specifying API version by using the resource id with api-version as a query parameter, e.g.
 terraform import azapi_resource.example /subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.ElasticSan/elasticSans/{resourceName}/volumeGroups/{resourceName}/volumes/{resourceName}?api-version=2025-09-01
 ```
