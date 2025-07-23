---
subcategory: "Microsoft.Synapse - Azure Synapse Analytics"
page_title: "workspaces/sqlPools/workloadGroups"
description: |-
  Manages a Synapse SQL Pool Workload Group.
---

# Microsoft.Synapse/workspaces/sqlPools/workloadGroups - Synapse SQL Pool Workload Group

This article demonstrates how to use `azapi` provider to manage the Synapse SQL Pool Workload Group resource in Azure.

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

variable "sql_administrator_login" {
  type        = string
  description = "The SQL administrator login for the Synapse workspace"
}

variable "sql_administrator_login_password" {
  type        = string
  description = "The SQL administrator login password for the Synapse workspace"
  sensitive   = true
}

resource "azapi_resource" "resourceGroup" {
  type     = "Microsoft.Resources/resourceGroups@2020-06-01"
  name     = var.resource_name
  location = var.location
}

resource "azapi_resource" "storageAccount" {
  type      = "Microsoft.Storage/storageAccounts@2021-09-01"
  parent_id = azapi_resource.resourceGroup.id
  name      = var.resource_name
  location  = var.location
  body = {
    kind = "StorageV2"
    properties = {
    }
    sku = {
      name = "Standard_LRS"
    }
  }
  schema_validation_enabled = false
  response_export_values    = ["*"]
}

data "azapi_resource_action" "listKeys" {
  type                   = "Microsoft.Storage/storageAccounts@2022-09-01"
  resource_id            = azapi_resource.storageAccount.id
  action                 = "listKeys"
  response_export_values = ["*"]
}

data "azapi_resource" "blobService" {
  type      = "Microsoft.Storage/storageAccounts/blobServices@2022-09-01"
  parent_id = azapi_resource.storageAccount.id
  name      = "default"
}

resource "azapi_resource" "container" {
  type      = "Microsoft.Storage/storageAccounts/blobServices/containers@2022-09-01"
  name      = var.resource_name
  parent_id = data.azapi_resource.blobService.id
  body = {
    properties = {
      metadata = {
        key = "value"
      }
    }
  }
  response_export_values = ["*"]
}

resource "azapi_resource" "workspace" {
  type      = "Microsoft.Synapse/workspaces@2021-06-01"
  parent_id = azapi_resource.resourceGroup.id
  name      = var.resource_name
  location  = var.location
  identity {
    type         = "SystemAssigned"
    identity_ids = []
  }
  body = {
    properties = {
      defaultDataLakeStorage = {
        accountUrl = azapi_resource.storageAccount.output.properties.primaryEndpoints.dfs
        filesystem = azapi_resource.container.name
      }
      managedVirtualNetwork         = ""
      publicNetworkAccess           = "Enabled"
      sqlAdministratorLogin         = var.sql_administrator_login
      sqlAdministratorLoginPassword = var.sql_administrator_login_password
    }
  }
  schema_validation_enabled = false
  response_export_values    = ["*"]
}

resource "azapi_resource" "sqlPool" {
  type      = "Microsoft.Synapse/workspaces/sqlPools@2021-06-01"
  parent_id = azapi_resource.workspace.id
  name      = var.resource_name
  location  = var.location
  body = {
    properties = {
      createMode = "Default"
    }
    sku = {
      name = "DW100c"
    }
  }
  schema_validation_enabled = false
  response_export_values    = ["*"]
}

resource "azapi_resource" "workloadGroup" {
  type      = "Microsoft.Synapse/workspaces/sqlPools/workloadGroups@2021-06-01"
  parent_id = azapi_resource.sqlPool.id
  name      = var.resource_name
  body = {
    properties = {
      importance                   = "normal"
      maxResourcePercent           = 100
      maxResourcePercentPerRequest = 3
      minResourcePercent           = 0
      minResourcePercentPerRequest = 3
    }
  }
  schema_validation_enabled = false
  response_export_values    = ["*"]
}


```



## Arguments Reference

The following arguments are supported:

* `type` - (Required) The type of the resource. This should be set to `Microsoft.Synapse/workspaces/sqlPools/workloadGroups@api-version`. The available api-versions for this resource are: [`2019-06-01-preview`, `2020-12-01`, `2021-03-01`, `2021-04-01-preview`, `2021-05-01`, `2021-06-01`, `2021-06-01-preview`].

* `parent_id` - (Required) The ID of the azure resource in which this resource is created. The allowed values are:  
  `/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Synapse/workspaces/{resourceName}/sqlPools/{resourceName}`

* `name` - (Required) Specifies the name of the azure resource. Changing this forces a new resource to be created.

* `body` - (Required) Specifies the configuration of the resource. More information about the arguments in `body` can be found in the [Microsoft documentation](https://learn.microsoft.com/en-us/azure/templates/Microsoft.Synapse/workspaces/sqlPools/workloadGroups?pivots=deployment-language-terraform).

For other arguments, please refer to the [azapi_resource](https://registry.terraform.io/providers/Azure/azapi/latest/docs/resources/resource) documentation.

## Import

 ```shell
 # Azure resource can be imported using the resource id, e.g.
 terraform import azapi_resource.example /subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Synapse/workspaces/{resourceName}/sqlPools/{resourceName}/workloadGroups/{resourceName}
 
 # It also supports specifying API version by using the resource id with api-version as a query parameter, e.g.
 terraform import azapi_resource.example /subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Synapse/workspaces/{resourceName}/sqlPools/{resourceName}/workloadGroups/{resourceName}?api-version=2021-06-01-preview
 ```
