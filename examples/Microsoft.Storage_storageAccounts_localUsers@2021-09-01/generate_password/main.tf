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
  default = "acctest0001dfdg"
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
      accessTier                   = "Hot"
      allowBlobPublicAccess        = true
      allowCrossTenantReplication  = true
      allowSharedKeyAccess         = true
      defaultToOAuthAuthentication = false
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
      isHnsEnabled      = false
      isNfsV3Enabled    = false
      isSftpEnabled     = false
      minimumTlsVersion = "TLS1_2"
      networkAcls = {
        defaultAction = "Allow"
      }
      publicNetworkAccess      = "Enabled"
      supportsHttpsTrafficOnly = true
    }
    sku = {
      name = "Standard_LRS"
    }
  }
  schema_validation_enabled = false
  response_export_values    = ["*"]
}

resource "azapi_resource" "localUser" {
  type      = "Microsoft.Storage/storageAccounts/localUsers@2021-09-01"
  parent_id = azapi_resource.storageAccount.id
  name      = var.resource_name
  body = {
    properties = {
      hasSshPassword = true,
      homeDirectory  = "containername/"
      hasSharedKey   = true,
      hasSshKey      = false,
      permissionScopes = [{
        permissions  = "cwl",
        service      = "blob",
        resourceName = "containername"
      }]
    }
  }
}

resource "azapi_resource_action" "localUser" {
  type        = "Microsoft.Storage/storageAccounts/localUsers@2022-05-01"
  resource_id = azapi_resource.localUser.id
  action      = "regeneratePassword"
  body = {
    username = "TestUserName"
  }
  response_export_values = ["sshPassword"]
}
