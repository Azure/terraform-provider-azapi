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

