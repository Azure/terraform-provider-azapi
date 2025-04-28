---
subcategory: "Microsoft.AVS - Azure VMware Solution"
page_title: "privateClouds"
description: |-
  Manages a Azure VMware Solution Private Cloud.
---

# Microsoft.AVS/privateClouds - Azure VMware Solution Private Cloud

This article demonstrates how to use `azapi` provider to manage the Azure VMware Solution Private Cloud resource in Azure.

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
  default = "centralus"
}

resource "azapi_resource" "resourceGroup" {
  type     = "Microsoft.Resources/resourceGroups@2020-06-01"
  name     = var.resource_name
  location = var.location
}

resource "azapi_resource" "privateCloud" {
  type      = "Microsoft.AVS/privateClouds@2022-05-01"
  parent_id = azapi_resource.resourceGroup.id
  name      = var.resource_name
  location  = var.location
  body = {
    properties = {
      internet = "Disabled"
      managementCluster = {
        clusterSize = 3
      }
      networkBlock = "192.168.48.0/22"
    }
    sku = {
      name = "av36"
    }
  }
  schema_validation_enabled = false
  response_export_values    = ["*"]
}


```



## Arguments Reference

The following arguments are supported:

* `type` - (Required) The type of the resource. This should be set to `Microsoft.AVS/privateClouds@api-version`. The available api-versions for this resource are: [`2020-03-20`, `2020-07-17-preview`, `2021-01-01-preview`, `2021-06-01`, `2021-12-01`, `2022-05-01`, `2023-03-01`, `2023-09-01`].

* `parent_id` - (Required) The ID of the azure resource in which this resource is created. The allowed values are:  
  `/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}`

* `name` - (Required) Specifies the name of the azure resource. Changing this forces a new resource to be created.

* `body` - (Required) Specifies the configuration of the resource. More information about the arguments in `body` can be found in the [Microsoft documentation](https://learn.microsoft.com/en-us/azure/templates/Microsoft.AVS/privateClouds?pivots=deployment-language-terraform).

## Import

 ```shell
 # Azure resource can be imported using the resource id, e.g.
 terraform import azapi_resource.example /subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.AVS/privateClouds/{resourceName}
 
 # It also supports specifying API version by using the resource id with api-version as a query parameter, e.g.
 terraform import azapi_resource.example /subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.AVS/privateClouds/{resourceName}?api-version=2023-09-01
 ```
