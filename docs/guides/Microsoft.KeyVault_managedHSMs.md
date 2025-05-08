---
subcategory: "Microsoft.KeyVault - Key Vault"
page_title: "managedHSMs"
description: |-
  Manages a Key Vault Managed Hardware Security Module.
---

# Microsoft.KeyVault/managedHSMs - Key Vault Managed Hardware Security Module

This article demonstrates how to use `azapi` provider to manage the Key Vault Managed Hardware Security Module resource in Azure.

## Example Usage

### default

```hcl
terraform {
  required_providers {
    azapi = {
      source = "Azure/azapi"
    }
    azurerm = {
      source = "hashicorp/azurerm"
    }
  }
}

provider "azurerm" {
  features {
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

data "azurerm_client_config" "current" {
}

resource "azapi_resource" "resourceGroup" {
  type     = "Microsoft.Resources/resourceGroups@2020-06-01"
  name     = var.resource_name
  location = var.location
}

resource "azapi_resource" "managedHSM" {
  type      = "Microsoft.KeyVault/managedHSMs@2021-10-01"
  parent_id = azapi_resource.resourceGroup.id
  name      = "kvHsm230630033342437496"
  location  = var.location
  body = {
    properties = {
      createMode            = "default"
      enablePurgeProtection = false
      enableSoftDelete      = true
      initialAdminObjectIds = [
        data.azurerm_client_config.current.object_id,
      ]
      publicNetworkAccess       = "Enabled"
      softDeleteRetentionInDays = 90
      tenantId                  = data.azurerm_client_config.current.tenant_id
    }
    sku = {
      family = "B"
      name   = "Standard_B1"
    }
  }
  schema_validation_enabled = false
  response_export_values    = ["*"]
}


```



## Arguments Reference

The following arguments are supported:

* `type` - (Required) The type of the resource. This should be set to `Microsoft.KeyVault/managedHSMs@api-version`. The available api-versions for this resource are: [`2020-04-01-preview`, `2021-04-01-preview`, `2021-06-01-preview`, `2021-10-01`, `2021-11-01-preview`, `2022-02-01-preview`, `2022-07-01`, `2022-11-01`, `2023-02-01`, `2023-07-01`, `2024-04-01-preview`, `2024-11-01`, `2024-12-01-preview`].

* `parent_id` - (Required) The ID of the azure resource in which this resource is created. The allowed values are:  
  `/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}`

* `name` - (Required) Specifies the name of the azure resource. Changing this forces a new resource to be created.

* `body` - (Required) Specifies the configuration of the resource. More information about the arguments in `body` can be found in the [Microsoft documentation](https://learn.microsoft.com/en-us/azure/templates/Microsoft.KeyVault/managedHSMs?pivots=deployment-language-terraform).

## Import

 ```shell
 # Azure resource can be imported using the resource id, e.g.
 terraform import azapi_resource.example /subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.KeyVault/managedHSMs/{resourceName}
 
 # It also supports specifying API version by using the resource id with api-version as a query parameter, e.g.
 terraform import azapi_resource.example /subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.KeyVault/managedHSMs/{resourceName}?api-version=2024-12-01-preview
 ```
