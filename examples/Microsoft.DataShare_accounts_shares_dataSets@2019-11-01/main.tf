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
  default = "eastus"
}

data "azapi_client_config" "current" {}

# Grant the Data Share account's managed identity rights on the Storage Account
# Storage Blob Data Owner roleDefinitionId
locals {
  storage_blob_data_owner = "b7e6dc6d-f1e8-4753-8033-0f276bb0955b"
}

resource "azapi_resource" "roleAssignment" {
  type      = "Microsoft.Authorization/roleAssignments@2022-04-01"
  name      = uuidv5("url", "${azapi_resource.storageAccount.id}/roleAssignments/${azapi_resource.account.output.identity.principalId}")
  parent_id = azapi_resource.storageAccount.id
  body = {
    properties = {
      roleDefinitionId = "/subscriptions/${data.azapi_client_config.current.subscription_id}/providers/Microsoft.Authorization/roleDefinitions/${local.storage_blob_data_owner}"
      principalId      = azapi_resource.account.output.identity.principalId
      principalType    = "ServicePrincipal"
    }
  }
  depends_on = [azapi_resource.account, azapi_resource.storageAccount]
}


resource "azapi_resource" "resourceGroup" {
  type     = "Microsoft.Resources/resourceGroups@2020-06-01"
  name     = var.resource_name
  location = var.location
}

resource "azapi_resource" "account" {
  type      = "Microsoft.DataShare/accounts@2019-11-01"
  parent_id = azapi_resource.resourceGroup.id
  name      = var.resource_name
  location  = var.location
  identity {
    type = "SystemAssigned"
  }
}

resource "azapi_resource" "storageAccount" {
  type      = "Microsoft.Storage/storageAccounts@2023-05-01"
  parent_id = azapi_resource.resourceGroup.id
  name      = var.resource_name
  location  = var.location
  body = {
    kind = "StorageV2"
    properties = {
      accessTier                   = "Hot"
      allowBlobPublicAccess        = false
      allowCrossTenantReplication  = false
      allowSharedKeyAccess         = true
      defaultToOAuthAuthentication = true
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
      isLocalUserEnabled = false
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
  name      = var.resource_name
  body = {
    properties = {
      publicAccess = "None"
    }
  }
}

resource "azapi_resource" "share" {
  type      = "Microsoft.DataShare/accounts/shares@2019-11-01"
  parent_id = azapi_resource.account.id
  name      = var.resource_name
  body = {
    properties = {
      description = ""
      shareKind   = "CopyBased"
      terms       = ""
    }
  }
}

resource "azapi_resource" "dataSet" {
  type      = "Microsoft.DataShare/accounts/shares/dataSets@2019-11-01"
  parent_id = azapi_resource.share.id
  name      = var.resource_name
  body = {
    kind = "AdlsGen2File"
    properties = {
      filePath           = "myfile.txt"
      fileSystem         = var.resource_name
      resourceGroup      = azapi_resource.resourceGroup.name
      storageAccountName = var.resource_name
      subscriptionId     = data.azapi_client_config.current.subscription_id
    }
  }
}
