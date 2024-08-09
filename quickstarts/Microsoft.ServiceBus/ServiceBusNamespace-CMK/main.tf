terraform {
  required_providers {
    azapi = {
      source = "Azure/azapi"
    }
  }
}

provider "azapi" {
}

provider "azurerm" {
  features {}
}

data "azurerm_client_config" "current" {}

resource "azurerm_resource_group" "test" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_user_assigned_identity" "test" {
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  name                = "example-identity"
}

resource "azapi_resource" "test" {
  name      = "example-namespace"
  type      = "Microsoft.ServiceBus/namespaces@2021-06-01-preview"
  parent_id = azurerm_resource_group.test.id
  location  = azurerm_resource_group.test.location

  identity {
    type         = "SystemAssigned, UserAssigned"
    identity_ids = [azurerm_user_assigned_identity.test.id]
  }

  body = {
    sku = {
      name = "Premium"
    }
  }
}

resource "azurerm_key_vault" "test" {
  name                     = "example-vault"
  location                 = azurerm_resource_group.test.location
  resource_group_name      = azurerm_resource_group.test.name
  tenant_id                = data.azurerm_client_config.current.tenant_id
  sku_name                 = "standard"
  soft_delete_enabled      = true
  purge_protection_enabled = true

  access_policy {
    tenant_id = azapi_resource.test.identity.0.tenant_id
    object_id = azapi_resource.test.identity.0.principal_id
    key_permissions = [
      "Get", "Create", "List", "Restore", "Recover", "UnwrapKey", "WrapKey", "Purge", "Encrypt", "Decrypt", "Sign", "Verify"
    ]
    secret_permissions = [
      "Get",
    ]
  }

  access_policy {
    tenant_id = data.azurerm_client_config.current.tenant_id
    object_id = data.azurerm_client_config.current.object_id
    key_permissions = [
      "Get", "Create", "Delete", "List", "Restore", "Recover", "UnwrapKey", "WrapKey", "Purge", "Encrypt", "Decrypt", "Sign", "Verify"
    ]
    secret_permissions = [
      "Get",
    ]
  }

  access_policy {
    tenant_id = azurerm_user_assigned_identity.test.tenant_id
    object_id = azurerm_user_assigned_identity.test.principal_id
    key_permissions = [
      "Get", "Create", "Delete", "List", "Restore", "Recover", "UnwrapKey", "WrapKey", "Purge", "Encrypt", "Decrypt", "Sign", "Verify"
    ]
    secret_permissions = [
      "Get",
    ]
  }
}

resource "azurerm_key_vault_key" "test" {
  name         = "example-sb-key"
  key_vault_id = azurerm_key_vault.test.id
  key_type     = "RSA"
  key_size     = 2048
  key_opts     = ["decrypt", "encrypt", "sign", "unwrapKey", "verify", "wrapKey"]
}

// patch resource used to enable CMK on servicebus namespace
resource "azapi_update_resource" "test" {
  resource_id = azapi_resource.test.resource_id
  type        = "Microsoft.ServiceBus/namespaces@2021-06-01-preview"

  body = {
    properties = {
      encryption = {
        keySource = "Microsoft.KeyVault"
        keyVaultProperties = [
          {
            identity = {
              userAssignedIdentity = azurerm_user_assigned_identity.test.id
            }
            keyName     = azurerm_key_vault_key.test.name
            keyVaultUri = azurerm_key_vault.test.vault_uri
            keyVersion  = azurerm_key_vault_key.test.version
          }
        ]
      }
    }
  }
}
