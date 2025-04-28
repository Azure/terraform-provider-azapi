---
subcategory: "Microsoft.Synapse - Azure Synapse Analytics"
page_title: "workspaces/bigDataPools"
description: |-
  Manages a Synapse Spark Pool.
---

# Microsoft.Synapse/workspaces/bigDataPools - Synapse Spark Pool

This article demonstrates how to use `azapi` provider to manage the Synapse Spark Pool resource in Azure.

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
      sqlAdministratorLogin         = "sqladminuser"
      sqlAdministratorLoginPassword = "H@Sh1CoR3!"
    }
  }
  schema_validation_enabled = false
  response_export_values    = ["*"]
}

resource "azapi_resource" "bigDataPool" {
  type      = "Microsoft.Synapse/workspaces/bigDataPools@2021-06-01-preview"
  parent_id = azapi_resource.workspace.id
  name      = var.resource_name
  location  = var.location
  body = {
    properties = {
      autoPause = {
        enabled = false
      }
      autoScale = {
        enabled = false
      }
      cacheSize             = 0
      defaultSparkLogFolder = "/logs"
      dynamicExecutorAllocation = {
        enabled      = false
        maxExecutors = 0
        minExecutors = 0
      }
      isComputeIsolationEnabled   = false
      nodeCount                   = 3
      nodeSize                    = "Small"
      nodeSizeFamily              = "MemoryOptimized"
      sessionLevelPackagesEnabled = false
      sparkEventsFolder           = "/events"
      sparkVersion                = "2.4"
    }
  }
  schema_validation_enabled = false
  response_export_values    = ["*"]
}


```



## Arguments Reference

The following arguments are supported:

* `type` - (Required) The type of the resource. This should be set to `Microsoft.Synapse/workspaces/bigDataPools@api-version`. The available api-versions for this resource are: [`2019-06-01-preview`, `2020-12-01`, `2021-03-01`, `2021-04-01-preview`, `2021-05-01`, `2021-06-01`, `2021-06-01-preview`].

* `parent_id` - (Required) The ID of the azure resource in which this resource is created. The allowed values are:  
  `/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Synapse/workspaces/{resourceName}`

* `name` - (Required) Specifies the name of the azure resource. Changing this forces a new resource to be created.

* `body` - (Required) Specifies the configuration of the resource. More information about the arguments in `body` can be found in the [Microsoft documentation](https://learn.microsoft.com/en-us/azure/templates/Microsoft.Synapse/workspaces/bigDataPools?pivots=deployment-language-terraform).

## Import

 ```shell
 # Azure resource can be imported using the resource id, e.g.
 terraform import azapi_resource.example /subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Synapse/workspaces/{resourceName}/bigDataPools/{resourceName}
 
 # It also supports specifying API version by using the resource id with api-version as a query parameter, e.g.
 terraform import azapi_resource.example /subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Synapse/workspaces/{resourceName}/bigDataPools/{resourceName}?api-version=2021-06-01-preview
 ```
