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

