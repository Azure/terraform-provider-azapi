---
subcategory: "Microsoft.RecoveryServices - Azure Site Recovery"
page_title: "vaults/replicationFabrics/replicationProtectionContainers"
description: |-
  Manages a site recovery services protection container on Azure.
---

# Microsoft.RecoveryServices/vaults/replicationFabrics/replicationProtectionContainers - site recovery services protection container on Azure

This article demonstrates how to use `azapi` provider to manage the site recovery services protection container on Azure resource in Azure.



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

resource "azapi_resource" "vault" {
  type      = "Microsoft.RecoveryServices/vaults@2022-10-01"
  parent_id = azapi_resource.resourceGroup.id
  name      = var.resource_name
  location  = var.location
  body = {
    properties = {
      publicNetworkAccess = "Enabled"
    }
    sku = {
      name = "Standard"
    }
  }
  schema_validation_enabled = false
  response_export_values    = ["*"]
}

resource "azapi_resource" "replicationFabric" {
  type      = "Microsoft.RecoveryServices/vaults/replicationFabrics@2022-10-01"
  parent_id = azapi_resource.vault.id
  name      = var.resource_name
  body = {
    properties = {
      customDetails = {
        instanceType = "Azure"
        location     = var.location
      }
    }
  }
  schema_validation_enabled = false
  response_export_values    = ["*"]
}

resource "azapi_resource" "replicationProtectionContainer" {
  type      = "Microsoft.RecoveryServices/vaults/replicationFabrics/replicationProtectionContainers@2022-10-01"
  parent_id = azapi_resource.replicationFabric.id
  name      = var.resource_name
  body = {
    properties = {
    }
  }
  schema_validation_enabled = false
  response_export_values    = ["*"]
}


```



## Arguments Reference

The following arguments are supported:

* `type` - (Required) The type of the resource. This should be set to `Microsoft.RecoveryServices/vaults/replicationFabrics/replicationProtectionContainers@api-version`. The available api-versions for this resource are: [`2016-08-10`, `2018-01-10`, `2018-07-10`, `2021-02-10`, `2021-03-01`, `2021-04-01`, `2021-06-01`, `2021-07-01`, `2021-08-01`, `2021-10-01`, `2021-11-01`, `2021-12-01`, `2022-01-01`, `2022-02-01`, `2022-03-01`, `2022-04-01`, `2022-05-01`, `2022-08-01`, `2022-09-10`, `2022-10-01`, `2023-01-01`, `2023-02-01`, `2023-04-01`, `2023-06-01`, `2023-08-01`, `2024-01-01`, `2024-02-01`, `2024-04-01`, `2024-10-01`, `2025-01-01`, `2025-02-01`].

* `parent_id` - (Required) The ID of the azure resource in which this resource is created. The allowed values are:  
  `/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.RecoveryServices/vaults/{resourceName}/replicationFabrics/{resourceName}`

* `name` - (Required) Specifies the name of the azure resource. Changing this forces a new resource to be created.

* `body` - (Required) Specifies the configuration of the resource. More information about the arguments in `body` can be found in the [Microsoft documentation](https://learn.microsoft.com/en-us/azure/templates/Microsoft.RecoveryServices/vaults/replicationFabrics/replicationProtectionContainers?pivots=deployment-language-terraform).

For other arguments, please refer to the [azapi_resource](https://registry.terraform.io/providers/Azure/azapi/latest/docs/resources/resource) documentation.

## Import

 ```shell
 # Azure resource can be imported using the resource id, e.g.
 terraform import azapi_resource.example /subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.RecoveryServices/vaults/{resourceName}/replicationFabrics/{resourceName}/replicationProtectionContainers/{resourceName}
 
 # It also supports specifying API version by using the resource id with api-version as a query parameter, e.g.
 terraform import azapi_resource.example /subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.RecoveryServices/vaults/{resourceName}/replicationFabrics/{resourceName}/replicationProtectionContainers/{resourceName}?api-version=2025-02-01
 ```
