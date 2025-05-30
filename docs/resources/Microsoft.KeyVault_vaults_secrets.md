---
subcategory: "Microsoft.KeyVault - Key Vault"
page_title: "vaults/secrets"
description: |-
  Manages a Key Vault Secrets.
---

# Microsoft.KeyVault/vaults/secrets - Key Vault Secrets

This article demonstrates how to use `azapi` provider to manage the Key Vault Secrets resource in Azure.

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

resource "azapi_resource" "vault" {
  type      = "Microsoft.KeyVault/vaults@2023-02-01"
  parent_id = azapi_resource.resourceGroup.id
  name      = var.resource_name
  location  = var.location
  body = {
    properties = {
      sku = {
        family = "A"
        name   = "standard"
      }
      accessPolicies   = []
      enableSoftDelete = true
      tenantId         = data.azurerm_client_config.current.tenant_id
    }
  }
  schema_validation_enabled = false
  response_export_values    = ["*"]
  lifecycle {
    ignore_changes = [body.properties.accessPolicies]
  }
}

data "azapi_resource_id" "secret" {
  type      = "Microsoft.KeyVault/vaults/secrets@2023-02-01"
  parent_id = azapi_resource.vault.id
  name      = var.resource_name
}

resource "azapi_resource_action" "put_secret" {
  type        = "Microsoft.KeyVault/vaults/secrets@2023-02-01"
  resource_id = data.azapi_resource_id.secret.id
  method      = "PUT"
  body = {
    properties = {
      value = "szechuan"
    }
  }
  response_export_values = ["*"]
}

```



## Arguments Reference

The following arguments are supported:

* `type` - (Required) The type of the resource. This should be set to `Microsoft.KeyVault/vaults/secrets@api-version`. The available api-versions for this resource are: [`2016-10-01`, `2018-02-14`, `2018-02-14-preview`, `2019-09-01`, `2020-04-01-preview`, `2021-04-01-preview`, `2021-06-01-preview`, `2021-10-01`, `2021-11-01-preview`, `2022-02-01-preview`, `2022-07-01`, `2022-11-01`, `2023-02-01`, `2023-07-01`, `2024-04-01-preview`, `2024-11-01`, `2024-12-01-preview`].

* `parent_id` - (Required) The ID of the azure resource in which this resource is created. The allowed values are:  
  `/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.KeyVault/vaults/{resourceName}`

* `name` - (Required) Specifies the name of the azure resource. Changing this forces a new resource to be created.

* `body` - (Required) Specifies the configuration of the resource. More information about the arguments in `body` can be found in the [Microsoft documentation](https://learn.microsoft.com/en-us/azure/templates/Microsoft.KeyVault/vaults/secrets?pivots=deployment-language-terraform).

For other arguments, please refer to the [azapi_resource](https://registry.terraform.io/providers/Azure/azapi/latest/docs/resources/resource) documentation.

## Import

 ```shell
 # Azure resource can be imported using the resource id, e.g.
 terraform import azapi_resource.example /subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.KeyVault/vaults/{resourceName}/secrets/{resourceName}
 
 # It also supports specifying API version by using the resource id with api-version as a query parameter, e.g.
 terraform import azapi_resource.example /subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.KeyVault/vaults/{resourceName}/secrets/{resourceName}?api-version=2024-12-01-preview
 ```
