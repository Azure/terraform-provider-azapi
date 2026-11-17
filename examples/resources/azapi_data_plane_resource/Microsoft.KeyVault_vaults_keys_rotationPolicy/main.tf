terraform {
  required_providers {
    azapi = {
      source = "Azure/azapi"
    }
  }
}

provider "azapi" {
}

variable "resource_name" {
  type    = string
  default = "acctest0001"
}

variable "location" {
  type    = string
  default = "westeurope"
}

data "azapi_client_config" "current" {}

locals {
  key_vault_crypto_officer_role_id = "/subscriptions/${data.azapi_client_config.current.subscription_id}/providers/Microsoft.Authorization/roleDefinitions/14b46e9e-c2b7-41b4-b07b-48a6ebf60603"
}

resource "azapi_resource" "resourceGroup" {
  type     = "Microsoft.Resources/resourceGroups@2021-04-01"
  name     = var.resource_name
  location = var.location
}

resource "azapi_resource" "vault" {
  type      = "Microsoft.KeyVault/vaults@2026-02-01"
  parent_id = azapi_resource.resourceGroup.id
  name      = var.resource_name
  location  = azapi_resource.resourceGroup.location
  body = {
    properties = {
      sku = {
        family = "A"
        name   = "standard"
      }
      tenantId                  = data.azapi_client_config.current.tenant_id
      enableRbacAuthorization   = true
      softDeleteRetentionInDays = 7
      accessPolicies            = []
    }
  }
  response_export_values = {
    vaultUri = "properties.vaultUri"
  }
}

resource "azapi_resource" "keyVaultCryptoOfficer" {
  type      = "Microsoft.Authorization/roleAssignments@2022-04-01"
  parent_id = azapi_resource.vault.id
  name      = uuidv5("url", "${azapi_resource.vault.id}/${data.azapi_client_config.current.object_id}/${local.key_vault_crypto_officer_role_id}")
  body = {
    properties = {
      principalId      = data.azapi_client_config.current.object_id
      roleDefinitionId = local.key_vault_crypto_officer_role_id
    }
  }
}

resource "azapi_data_plane_resource" "key" {
  type      = "Microsoft.KeyVault/vaults/keys@2025-07-01"
  parent_id = trimsuffix(trimprefix(azapi_resource.vault.output.vaultUri, "https://"), "/")
  name      = var.resource_name
  body = {
    kty      = "RSA"
    key_size = 2048
    key_ops  = ["wrapKey", "unwrapKey"]
  }

  retry = {
    error_message_regex  = ["Forbidden", "Unauthorized", "authorization"]
    interval_seconds     = 10
    max_interval_seconds = 60
  }

  depends_on = [
    azapi_resource.keyVaultCryptoOfficer
  ]
}

# Key Vault always returns the Rotate action before the Notify action, and adds a Notify action
# 30 days before expiry when there is none, so list both in that order to avoid a diff.
# Destroying this resource sends a policy with no lifetime actions, so Key Vault stops rotating the key.
resource "azapi_data_plane_resource" "rotationPolicy" {
  type      = "Microsoft.KeyVault/vaults/keys/rotationPolicy@2025-07-01"
  parent_id = azapi_data_plane_resource.key.id
  body = {
    lifetimeActions = [
      {
        trigger = {
          timeAfterCreate = "P90D"
        }
        action = {
          type = "Rotate"
        }
      },
      {
        trigger = {
          timeBeforeExpiry = "P30D"
        }
        action = {
          type = "Notify"
        }
      },
    ]
    attributes = {
      expiryTime = "P1Y"
    }
  }
}
