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

data "azapi_client_config" "current" {}

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

resource "azapi_resource" "component" {
  type      = "Microsoft.Insights/components@2020-02-02"
  parent_id = azapi_resource.resourceGroup.id
  name      = "${var.resource_name}-ai"
  location  = var.location
  body = {
    kind = "web"
    properties = {
      Application_Type                = "web"
      DisableIpMasking                = false
      DisableLocalAuth                = false
      ForceCustomerStorageForProfiler = false
      RetentionInDays                 = 90
      SamplingPercentage              = 100
      publicNetworkAccessForIngestion = "Enabled"
      publicNetworkAccessForQuery     = "Enabled"
    }
  }
}

resource "azapi_resource" "vault" {
  type      = "Microsoft.KeyVault/vaults@2023-02-01"
  parent_id = azapi_resource.resourceGroup.id
  name      = "${var.resource_name}vault"
  location  = var.location
  body = {
    properties = {
      accessPolicies               = []
      createMode                   = "default"
      enablePurgeProtection        = true
      enableRbacAuthorization      = false
      enableSoftDelete             = true
      enabledForDeployment         = false
      enabledForDiskEncryption     = false
      enabledForTemplateDeployment = false
      publicNetworkAccess          = "Enabled"
      sku = {
        family = "A"
        name   = "standard"
      }
      tenantId = data.azapi_client_config.current.tenant_id
    }
  }
}

resource "azapi_resource" "storageAccount" {
  type      = "Microsoft.Storage/storageAccounts@2023-05-01"
  parent_id = azapi_resource.resourceGroup.id
  name      = "${lower(substr(var.resource_name, 0, 16))}acc"
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

resource "azapi_resource" "container" {
  type      = "Microsoft.Storage/storageAccounts/blobServices/containers@2023-05-01"
  parent_id = "${azapi_resource.storageAccount.id}/blobServices/default"
  name      = "datacontainer"
  body = {
    properties = {
      publicAccess = "None"
    }
  }
}

data "azapi_resource_action" "storage_keys" {
  type                   = "Microsoft.Storage/storageAccounts@2023-05-01"
  resource_id            = azapi_resource.storageAccount.id
  action                 = "listKeys"
  method                 = "POST"
  response_export_values = ["*"]
}

resource "azapi_resource" "workspace" {
  type      = "Microsoft.MachineLearningServices/workspaces@2024-04-01"
  parent_id = azapi_resource.resourceGroup.id
  name      = "${var.resource_name}-mlw"
  location  = var.location
  identity {
    type         = "SystemAssigned"
    identity_ids = []
  }
  body = {
    kind = "Default"
    properties = {
      applicationInsights = azapi_resource.component.id
      keyVault            = azapi_resource.vault.id
      publicNetworkAccess = "Enabled"
      storageAccount      = azapi_resource.storageAccount.id
      v1LegacyMode        = false
    }
    sku = {
      name = "Basic"
      tier = "Basic"
    }
  }
}

resource "azapi_resource" "dataStore" {
  type      = "Microsoft.MachineLearningServices/workspaces/dataStores@2024-04-01"
  parent_id = azapi_resource.workspace.id
  name      = replace("${var.resource_name}_ds", "-", "_")
  body = {
    properties = {
      accountName   = azapi_resource.storageAccount.name
      containerName = azapi_resource.container.name
      credentials = {
        credentialsType = "AccountKey"
        secrets = {
          key         = base64encode(data.azapi_resource_action.storage_keys.output.keys[0].value)
          secretsType = "AccountKey"
        }
      }
      datastoreType                 = "AzureBlob"
      description                   = ""
      endpoint                      = "core.windows.net"
      serviceDataAccessAuthIdentity = "None"
      tags                          = null
    }
  }
  depends_on = [azapi_resource.container]
}
