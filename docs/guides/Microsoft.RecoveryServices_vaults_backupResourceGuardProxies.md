---
subcategory: "Microsoft.RecoveryServices - Azure Site Recovery"
page_title: "vaults/backupResourceGuardProxies"
description: |-
  Manages a association of a Resource Guard and Recovery Services Vault.
---

# Microsoft.RecoveryServices/vaults/backupResourceGuardProxies - association of a Resource Guard and Recovery Services Vault

This article demonstrates how to use `azapi` provider to manage the association of a Resource Guard and Recovery Services Vault resource in Azure.

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

resource "azapi_resource" "resourceGuard" {
  type      = "Microsoft.DataProtection/resourceGuards@2022-04-01"
  parent_id = azapi_resource.resourceGroup.id
  name      = var.resource_name
  location  = var.location
  body = {
    properties = {
      vaultCriticalOperationExclusionList = [
      ]
    }
  }
  schema_validation_enabled = false
  response_export_values    = ["*"]
}

resource "azapi_resource" "backupResourceGuardProxy" {
  type      = "Microsoft.RecoveryServices/vaults/backupResourceGuardProxies@2023-02-01"
  parent_id = azapi_resource.vault.id
  name      = var.resource_name
  body = {
    properties = {
      resourceGuardResourceId = azapi_resource.resourceGuard.id
    }
    type = "Microsoft.RecoveryServices/vaults/backupResourceGuardProxies"
  }
  schema_validation_enabled = false
  response_export_values    = ["*"]
}

```



## Arguments Reference

The following arguments are supported:

* `type` - (Required) The type of the resource. This should be set to `Microsoft.RecoveryServices/vaults/backupResourceGuardProxies@api-version`. The available api-versions for this resource are: [`2021-02-01-preview`, `2021-07-01`, `2021-08-01`, `2021-10-01`, `2021-12-01`, `2022-01-01`, `2022-02-01`, `2022-03-01`, `2022-04-01`, `2022-06-01-preview`, `2022-09-01-preview`, `2022-09-30-preview`, `2022-10-01`, `2023-01-01`, `2023-02-01`, `2023-04-01`, `2023-06-01`, `2023-08-01`, `2024-01-01`, `2024-02-01`, `2024-04-01`, `2024-04-30-preview`, `2024-07-30-preview`, `2024-10-01`, `2024-11-01-preview`, `2025-01-01`, `2025-02-01`, `2025-02-28-preview`].

* `parent_id` - (Required) The ID of the azure resource in which this resource is created. The allowed values are:  
  `/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.RecoveryServices/vaults/{resourceName}`

* `name` - (Required) Specifies the name of the azure resource. Changing this forces a new resource to be created.

* `body` - (Required) Specifies the configuration of the resource. More information about the arguments in `body` can be found in the [Microsoft documentation](https://learn.microsoft.com/en-us/azure/templates/Microsoft.RecoveryServices/vaults/backupResourceGuardProxies?pivots=deployment-language-terraform).

For other arguments, please refer to the [azapi_resource](https://registry.terraform.io/providers/Azure/azapi/latest/docs/resources/resource) documentation.

## Import

 ```shell
 # Azure resource can be imported using the resource id, e.g.
 terraform import azapi_resource.example /subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.RecoveryServices/vaults/{resourceName}/backupResourceGuardProxies/{resourceName}
 
 # It also supports specifying API version by using the resource id with api-version as a query parameter, e.g.
 terraform import azapi_resource.example /subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.RecoveryServices/vaults/{resourceName}/backupResourceGuardProxies/{resourceName}?api-version=2025-02-28-preview
 ```
