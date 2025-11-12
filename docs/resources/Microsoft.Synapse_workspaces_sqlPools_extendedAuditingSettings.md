---
subcategory: "Microsoft.Synapse - Azure Synapse Analytics"
page_title: "workspaces/sqlPools/extendedAuditingSettings"
description: |-
  Manages a Synapse SQL Pool Extended Auditing Policy.
---

# Microsoft.Synapse/workspaces/sqlPools/extendedAuditingSettings - Synapse SQL Pool Extended Auditing Policy

This article demonstrates how to use `azapi` provider to manage the Synapse SQL Pool Extended Auditing Policy resource in Azure.



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

variable "sql_administrator_login_password" {
  type        = string
  sensitive   = true
  description = "The SQL administrator login password for the Synapse workspace"
}

resource "azapi_resource" "resourceGroup" {
  type     = "Microsoft.Resources/resourceGroups@2020-06-01"
  name     = var.resource_name
  location = var.location
}

resource "azapi_resource" "storageAccount" {
  type      = "Microsoft.Storage/storageAccounts@2023-05-01"
  parent_id = azapi_resource.resourceGroup.id
  name      = "${var.resource_name}st1"
  location  = var.location
  body = {
    kind = "BlobStorage"
    properties = {
      accessTier                   = "Hot"
      allowBlobPublicAccess        = true
      allowCrossTenantReplication  = false
      allowSharedKeyAccess         = true
      defaultToOAuthAuthentication = false
      dnsEndpointType              = "Standard"
      encryption = {
        keySource = "Microsoft.Storage"
        services = {
          queue = {
            keyType = "Service"
          }
          table = {
            keyType = "Service"
          }
        }
      }
      isHnsEnabled       = false
      isLocalUserEnabled = true
      isNfsV3Enabled     = false
      isSftpEnabled      = false
      minimumTlsVersion  = "TLS1_2"
      networkAcls = {
        bypass              = "AzureServices"
        defaultAction       = "Allow"
        ipRules             = []
        resourceAccessRules = []
        virtualNetworkRules = []
      }
      publicNetworkAccess      = "Enabled"
      supportsHttpsTrafficOnly = true
    }
    sku = {
      name = "Standard_LRS"
    }
  }
}

resource "azapi_resource" "storageAccount_1" {
  type      = "Microsoft.Storage/storageAccounts@2023-05-01"
  parent_id = azapi_resource.resourceGroup.id
  name      = "${var.resource_name}st2"
  location  = var.location
  body = {
    kind = "StorageV2"
    properties = {
      accessTier                   = "Hot"
      allowBlobPublicAccess        = true
      allowCrossTenantReplication  = false
      allowSharedKeyAccess         = true
      defaultToOAuthAuthentication = false
      dnsEndpointType              = "Standard"
      encryption = {
        keySource = "Microsoft.Storage"
        services = {
          queue = {
            keyType = "Service"
          }
          table = {
            keyType = "Service"
          }
        }
      }
      isHnsEnabled       = true
      isLocalUserEnabled = true
      isNfsV3Enabled     = false
      isSftpEnabled      = false
      minimumTlsVersion  = "TLS1_2"
      networkAcls = {
        bypass              = "AzureServices"
        defaultAction       = "Allow"
        ipRules             = []
        resourceAccessRules = []
        virtualNetworkRules = []
      }
      publicNetworkAccess      = "Enabled"
      supportsHttpsTrafficOnly = true
    }
    sku = {
      name = "Standard_LRS"
    }
  }
}

resource "azapi_resource" "filesystem" {
  type      = "Microsoft.Storage/storageAccounts/blobServices/containers@2023-05-01"
  parent_id = "${azapi_resource.storageAccount_1.id}/blobServices/default"
  name      = "${var.resource_name}fs"
  body = {
    properties = {}
  }
}

resource "azapi_resource" "workspace" {
  type      = "Microsoft.Synapse/workspaces@2021-06-01"
  parent_id = azapi_resource.resourceGroup.id
  name      = "${var.resource_name}sw"
  location  = var.location
  identity {
    type         = "SystemAssigned"
    identity_ids = []
  }
  body = {
    properties = {
      azureADOnlyAuthentication = false
      defaultDataLakeStorage = {
        accountUrl = "https://${azapi_resource.storageAccount_1.name}.dfs.core.windows.net"
        filesystem = azapi_resource.filesystem.name
      }
      managedResourceGroupName         = ""
      managedVirtualNetwork            = ""
      publicNetworkAccess              = "Enabled"
      sqlAdministratorLogin            = "sqladminuser"
      sqlAdministratorLoginPassword    = var.sql_administrator_login_password
      workspaceRepositoryConfiguration = {}
    }
  }
  depends_on = [azapi_resource.filesystem]
}

resource "azapi_resource" "sqlPool" {
  type      = "Microsoft.Synapse/workspaces/sqlPools@2021-06-01"
  parent_id = azapi_resource.workspace.id
  name      = "${var.resource_name}sp"
  location  = var.location
  body = {
    properties = {
      collation          = ""
      createMode         = "Default"
      storageAccountType = "GRS"
    }
    sku = {
      name = "DW100c"
    }
  }
}

resource "azapi_resource_action" "storageAccountKeys" {
  type                   = "Microsoft.Storage/storageAccounts@2023-05-01"
  resource_id            = azapi_resource.storageAccount.id
  action                 = "listKeys"
  response_export_values = ["keys"]
}

resource "azapi_resource" "extendedAuditingSetting" {
  type      = "Microsoft.Synapse/workspaces/sqlPools/extendedAuditingSettings@2021-06-01"
  parent_id = azapi_resource.sqlPool.id
  name      = "default"
  body = {
    properties = {
      isAzureMonitorTargetEnabled = true
      isStorageSecondaryKeyInUse  = false
      retentionDays               = 0
      state                       = "Enabled"
      storageAccountAccessKey     = azapi_resource_action.storageAccountKeys.output.keys[0].value
      storageEndpoint             = "https://${azapi_resource.storageAccount.name}.blob.core.windows.net/"
    }
  }
}


```



## Arguments Reference

The following arguments are supported:

* `type` - (Required) The type of the resource. This should be set to `Microsoft.Synapse/workspaces/sqlPools/extendedAuditingSettings@api-version`. The available api-versions for this resource are: [`2019-06-01-preview`, `2020-12-01`, `2021-03-01`, `2021-04-01-preview`, `2021-05-01`, `2021-06-01`, `2021-06-01-preview`].

* `parent_id` - (Required) The ID of the azure resource in which this resource is created. The allowed values are:  
  `/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Synapse/workspaces/{resourceName}/sqlPools/{resourceName}`

* `name` - (Required) Specifies the name of the azure resource. Changing this forces a new resource to be created.

* `body` - (Required) Specifies the configuration of the resource. More information about the arguments in `body` can be found in the [Microsoft documentation](https://learn.microsoft.com/en-us/azure/templates/Microsoft.Synapse/workspaces/sqlPools/extendedAuditingSettings?pivots=deployment-language-terraform).

For other arguments, please refer to the [azapi_resource](https://registry.terraform.io/providers/Azure/azapi/latest/docs/resources/resource) documentation.

## Import

 ```shell
 # Azure resource can be imported using the resource id, e.g.
 terraform import azapi_resource.example /subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Synapse/workspaces/{resourceName}/sqlPools/{resourceName}/extendedAuditingSettings/{resourceName}
 
 # It also supports specifying API version by using the resource id with api-version as a query parameter, e.g.
 terraform import azapi_resource.example /subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Synapse/workspaces/{resourceName}/sqlPools/{resourceName}/extendedAuditingSettings/{resourceName}?api-version=2021-06-01-preview
 ```
