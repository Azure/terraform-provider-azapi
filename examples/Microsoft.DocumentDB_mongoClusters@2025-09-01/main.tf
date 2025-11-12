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
  default = "westus3"
}

variable "mongo_admin_password" {
  type        = string
  description = "The administrator password for the MongoDB cluster"
  sensitive   = true
}

resource "azapi_resource" "resourceGroup" {
  type     = "Microsoft.Resources/resourceGroups@2020-06-01"
  name     = var.resource_name
  location = var.location
}

resource "azapi_resource" "userAssignedIdentity" {
  type                      = "Microsoft.ManagedIdentity/userAssignedIdentities@2023-01-31"
  parent_id                 = azapi_resource.resourceGroup.id
  name                      = var.resource_name
  location                  = var.location
  schema_validation_enabled = false
  response_export_values    = ["*"]
}

data "azapi_client_config" "current" {}

resource "azapi_resource" "vault" {
  type      = "Microsoft.KeyVault/vaults@2023-02-01"
  parent_id = azapi_resource.resourceGroup.id
  name      = "${var.resource_name}-kv"
  location  = var.location
  body = {
    properties = {
      createMode                   = "default"
      enablePurgeProtection        = true
      enableSoftDelete             = true
      enableRbacAuthorization      = true
      enabledForDeployment         = true
      enabledForDiskEncryption     = true
      enabledForTemplateDeployment = true
      publicNetworkAccess          = "Enabled"
      accessPolicies               = []
      sku = {
        family = "A"
        name   = "standard"
      }
      tenantId = data.azapi_client_config.current.tenant_id
    }
  }
}

data "azapi_resource_list" "kvCertificatesOfficerRoleDefinition" {
  type      = "Microsoft.Authorization/roleDefinitions@2022-04-01"
  parent_id = azapi_resource.vault.id
  response_export_values = {
    definition_id = "value[?properties.roleName == 'Key Vault Crypto Officer'].id | [0]"
  }
}

resource "azapi_resource" "kvRoleAssignmentTf" {
  type      = "Microsoft.Authorization/roleAssignments@2022-04-01"
  parent_id = azapi_resource.vault.id
  name      = uuid()
  body = {
    properties = {
      principalId      = data.azapi_client_config.current.object_id
      roleDefinitionId = data.azapi_resource_list.kvCertificatesOfficerRoleDefinition.output.definition_id
    }
  }
  lifecycle {
    ignore_changes = [name]
  }
}

data "azapi_resource_list" "kvCertificatesUserRoleDefinition" {
  type      = "Microsoft.Authorization/roleDefinitions@2022-04-01"
  parent_id = azapi_resource.vault.id
  response_export_values = {
    definition_id = "value[?properties.roleName == 'Key Vault Crypto Service Encryption User'].id | [0]"
  }
}

resource "azapi_resource" "kvRoleAssignmentIdentity" {
  type      = "Microsoft.Authorization/roleAssignments@2022-04-01"
  parent_id = azapi_resource.vault.id
  name      = uuid()
  body = {
    properties = {
      principalId      = azapi_resource.userAssignedIdentity.output.properties.principalId
      roleDefinitionId = data.azapi_resource_list.kvCertificatesUserRoleDefinition.output.definition_id
    }
  }
  lifecycle {
    ignore_changes = [name]
  }
}

data "azapi_resource_id" "key" {
  type      = "Microsoft.KeyVault/vaults/keys@2023-02-01"
  parent_id = azapi_resource.vault.id
  name      = var.resource_name
}

resource "azapi_resource_action" "key" {
  type        = "Microsoft.KeyVault/vaults/keys@2023-02-01"
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
  depends_on = [
    azapi_resource.kvRoleAssignmentTf,
    azapi_resource.kvRoleAssignmentIdentity,
  ]
}

resource "azapi_resource" "mongoCluster" {
  type      = "Microsoft.DocumentDB/mongoClusters@2025-09-01"
  parent_id = azapi_resource.resourceGroup.id
  name      = var.resource_name
  location  = var.location
  identity {
    type         = "UserAssigned"
    identity_ids = [azapi_resource.userAssignedIdentity.id]
  }
  body = {
    properties = {
      administrator = {
        userName = "mongoAdmin"
        password = var.mongo_admin_password
      }
      storage = {
        sizeGb = 32
      }
      compute = {
        tier = "M40"
      }
      sharding = {
        shardCount = 1
      }
      highAvailability = {
        targetMode = "Disabled"
      }
      encryption = {
        customerManagedKeyEncryption = {
          keyEncryptionKeyIdentity = {
            identityType                   = "UserAssignedIdentity"
            userAssignedIdentityResourceId = azapi_resource.userAssignedIdentity.id
          }
          keyEncryptionKeyUrl = azapi_resource_action.key.output.properties.keyUri
        }
      }
    }
  }
  schema_validation_enabled = false
  response_export_values    = ["*"]
}
