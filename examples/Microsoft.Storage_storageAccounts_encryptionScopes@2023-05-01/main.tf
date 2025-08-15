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

resource "azapi_resource" "vault" {
  type      = "Microsoft.KeyVault/vaults@2023-02-01"
  parent_id = azapi_resource.resourceGroup.id
  name      = "${var.resource_name}-kv"
  location  = var.location
  body = {
    properties = {
      accessPolicies = [
        {
          objectId = data.azapi_client_config.current.object_id
          permissions = {
            certificates = []
            keys         = ["Get", "Create", "Delete", "List", "Restore", "Recover", "UnwrapKey", "WrapKey", "Purge", "Encrypt", "Decrypt", "Sign", "Verify", "GetRotationPolicy"]
            secrets      = []
            storage      = []
          }
          tenantId = data.azapi_client_config.current.tenant_id
        },
        {
          objectId = azapi_resource.storageAccount.identity[0].principal_id
          permissions = {
            certificates = []
            keys         = ["Get", "UnwrapKey", "WrapKey"]
            secrets      = []
            storage      = []
          }
          tenantId = data.azapi_client_config.current.tenant_id
        }
      ]
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
  depends_on = [azapi_resource.storageAccount]
}

resource "azapi_resource" "storageAccount" {
  type      = "Microsoft.Storage/storageAccounts@2023-05-01"
  parent_id = azapi_resource.resourceGroup.id
  name      = "${var.resource_name}sa"
  location  = var.location
  identity {
    type         = "SystemAssigned"
    identity_ids = []
  }
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

# Commenting out first access policy due to provider state tracking issues
# resource "azapi_resource" "accessPolicy" {
#   type      = "Microsoft.KeyVault/vaults/accessPolicies@2023-02-01"
#   parent_id = azapi_resource.vault.id
#   name      = "add"
#   body = {
#     properties = {
#       accessPolicies = [{
#         objectId = data.azapi_client_config.current.object_id
#         permissions = {
#           certificates = []
#           keys         = ["Get", "Create", "Delete", "List", "Restore", "Recover", "UnwrapKey", "WrapKey", "Purge", "Encrypt", "Decrypt", "Sign", "Verify", "GetRotationPolicy"]
#           secrets      = []
#           storage      = []
#         }
#         tenantId = data.azapi_client_config.current.tenant_id
#       }]
#     }
#   }
# }

# Commenting out storage account access policy due to provider state tracking issues
# resource "azapi_resource" "accessPolicy_1" {
#   type      = "Microsoft.KeyVault/vaults/accessPolicies@2023-02-01"
#   parent_id = azapi_resource.vault.id
#   name      = "add"
#   body = {
#     properties = {
#       accessPolicies = [{
#         objectId = azapi_resource.storageAccount.identity[0].principal_id
#         permissions = {
#           certificates = []
#           keys         = ["Get", "UnwrapKey", "WrapKey"]
#           secrets      = []
#           storage      = []
#         }
#         tenantId = data.azapi_client_config.current.tenant_id
#       }]
#     }
#   }
# }

resource "azapi_resource" "key" {
  type      = "Microsoft.KeyVault/vaults/keys@2023-02-01"
  parent_id = azapi_resource.vault.id
  name      = "${var.resource_name}-key"
  body = {
    properties = {
      kty     = "RSA"
      keySize = 2048
      keyOps  = ["encrypt", "decrypt", "sign", "verify", "wrapKey", "unwrapKey"]
    }
  }
}

resource "azapi_resource" "encryptionScope" {
  type      = "Microsoft.Storage/storageAccounts/encryptionScopes@2023-05-01"
  parent_id = azapi_resource.storageAccount.id
  name      = "${var.resource_name}-scope"
  body = {
    properties = {
      keyVaultProperties = {
        keyUri = azapi_resource.key.output.properties.keyUriWithVersion
      }
      source = "Microsoft.KeyVault"
      state  = "Enabled"
    }
  }
  depends_on = [azapi_resource.vault]
}

