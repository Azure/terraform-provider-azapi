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

data "azapi_client_config" "current" {}

locals {
  key_vault_crypto_officer_role_id = "/subscriptions/${data.azapi_client_config.current.subscription_id}/providers/Microsoft.Authorization/roleDefinitions/14b46e9e-c2b7-41b4-b07b-48a6ebf60603"
}

resource "azapi_resource" "resourceGroup" {
  type     = "Microsoft.Resources/resourceGroups@2020-06-01"
  name     = var.resource_name
  location = var.location
}

resource "azapi_resource" "vault" {
  type      = "Microsoft.KeyVault/vaults@2026-02-01"
  parent_id = azapi_resource.resourceGroup.id
  name      = var.resource_name
  location  = var.location
  body = {
    properties = {
      sku = {
        family = "A"
        name   = "standard"
      }
      accessPolicies          = []
      enableRbacAuthorization = true
      enableSoftDelete        = true
      enablePurgeProtection   = true
      tenantId                = data.azapi_client_config.current.tenant_id
    }
  }
  response_export_values = ["*"]
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

data "azapi_resource_id" "key" {
  type      = "Microsoft.KeyVault/vaults/keys@2026-02-01"
  parent_id = azapi_resource.vault.id
  name      = var.resource_name
}

resource "azapi_resource_action" "put_key" {
  type        = "Microsoft.KeyVault/vaults/keys@2026-02-01"
  resource_id = data.azapi_resource_id.key.id
  method      = "PUT"
  body = {
    properties = {
      keySize = 2048
      kty     = "RSA"
      keyOps  = ["encrypt", "decrypt", "sign", "verify", "wrapKey", "unwrapKey"]
    }
  }
  response_export_values = ["*"]
  retry = {
    error_message_regex  = ["Forbidden", "Unauthorized", "authorization"]
    interval_seconds     = 10
    max_interval_seconds = 60
  }
  depends_on = [azapi_resource.keyVaultCryptoOfficer]
}
