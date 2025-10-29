---
subcategory: "Microsoft.Synapse - Azure Synapse Analytics"
page_title: "workspaces/sqlPools/securityAlertPolicies"
description: |-
  Manages a Security Alert Policy for a Synapse SQL Pool.
---

# Microsoft.Synapse/workspaces/sqlPools/securityAlertPolicies - Security Alert Policy for a Synapse SQL Pool

This article demonstrates how to use `azapi` provider to manage the Security Alert Policy for a Synapse SQL Pool resource in Azure.



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

resource "azapi_resource" "storageAccount" {
  type      = "Microsoft.Storage/storageAccounts@2023-05-01"
  parent_id = azapi_resource.resourceGroup.id
  name      = "${var.resource_name}st"
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
  parent_id = "${azapi_resource.storageAccount.id}/blobServices/default"
  name      = "${var.resource_name}fs"
  body = {
    properties = {
      publicAccess = "None"
    }
  }
}

resource "azapi_resource" "workspace" {
  type      = "Microsoft.Synapse/workspaces@2021-06-01"
  parent_id = azapi_resource.resourceGroup.id
  name      = "${var.resource_name}ws"
  location  = var.location
  identity {
    type         = "SystemAssigned"
    identity_ids = []
  }
  body = {
    properties = {
      azureADOnlyAuthentication = false
      defaultDataLakeStorage = {
        accountUrl = azapi_resource.storageAccount.output.properties.primaryEndpoints.dfs
        filesystem = azapi_resource.filesystem.name
      }
      managedResourceGroupName         = ""
      managedVirtualNetwork            = ""
      publicNetworkAccess              = "Enabled"
      sqlAdministratorLogin            = "sqladminuser"
      sqlAdministratorLoginPassword    = "H@Sh1CoR3!"
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

resource "azapi_resource_action" "storageKeys" {
  type                   = "Microsoft.Storage/storageAccounts@2023-05-01"
  resource_id            = azapi_resource.storageAccount.id
  action                 = "listKeys"
  response_export_values = ["keys"]
}

resource "azapi_resource" "securityAlertPolicy" {
  type      = "Microsoft.Synapse/workspaces/sqlPools/securityAlertPolicies@2021-06-01"
  parent_id = azapi_resource.sqlPool.id
  name      = "default"
  body = {
    properties = {
      disabledAlerts          = ["Data_Exfiltration", "Sql_Injection"]
      retentionDays           = 20
      state                   = "Enabled"
      storageAccountAccessKey = azapi_resource_action.storageKeys.output.keys[0].value
      storageEndpoint         = azapi_resource.storageAccount.output.properties.primaryEndpoints.blob
    }
  }
}


```



## Arguments Reference

The following arguments are supported:

* `type` - (Required) The type of the resource. This should be set to `Microsoft.Synapse/workspaces/sqlPools/securityAlertPolicies@api-version`. The available api-versions for this resource are: [`2019-06-01-preview`, `2020-12-01`, `2021-03-01`, `2021-04-01-preview`, `2021-05-01`, `2021-06-01`, `2021-06-01-preview`].

* `parent_id` - (Required) The ID of the azure resource in which this resource is created. The allowed values are:  
  `/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Synapse/workspaces/{resourceName}/sqlPools/{resourceName}`

* `name` - (Required) Specifies the name of the azure resource. Changing this forces a new resource to be created.

* `body` - (Required) Specifies the configuration of the resource. More information about the arguments in `body` can be found in the [Microsoft documentation](https://learn.microsoft.com/en-us/azure/templates/Microsoft.Synapse/workspaces/sqlPools/securityAlertPolicies?pivots=deployment-language-terraform).

For other arguments, please refer to the [azapi_resource](https://registry.terraform.io/providers/Azure/azapi/latest/docs/resources/resource) documentation.

## Import

 ```shell
 # Azure resource can be imported using the resource id, e.g.
 terraform import azapi_resource.example /subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Synapse/workspaces/{resourceName}/sqlPools/{resourceName}/securityAlertPolicies/{resourceName}
 
 # It also supports specifying API version by using the resource id with api-version as a query parameter, e.g.
 terraform import azapi_resource.example /subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Synapse/workspaces/{resourceName}/sqlPools/{resourceName}/securityAlertPolicies/{resourceName}?api-version=2021-06-01-preview
 ```
