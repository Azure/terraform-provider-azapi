terraform {
  required_providers {
    azapi = {
      source = "Azure/azapi"
    }
    time = {
      source = "hashicorp/time"
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

variable "replica_location" {
  type    = string
  default = "centralus"
}

variable "mongo_admin_username" {
  type    = string
  default = "mongoAdmin"
}

variable "mongo_admin_password" {
  type        = string
  description = "The administrator password for the MongoDB cluster"
  sensitive   = true
}

variable "mongo_restore_admin_password" {
  type        = string
  description = "The administrator password for the restored MongoDB cluster"
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

# replica key vault

resource "azapi_resource" "vault_replica" {
  type      = "Microsoft.KeyVault/vaults@2023-02-01"
  parent_id = azapi_resource.resourceGroup.id
  name      = "${var.resource_name}-kv-replica"
  location  = var.replica_location
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

data "azapi_resource_list" "kvCertificatesOfficerRoleDefinition_replica" {
  type      = "Microsoft.Authorization/roleDefinitions@2022-04-01"
  parent_id = azapi_resource.vault_replica.id
  response_export_values = {
    definition_id = "value[?properties.roleName == 'Key Vault Crypto Officer'].id | [0]"
  }
}

resource "azapi_resource" "kvRoleAssignmentTf_replica" {
  type      = "Microsoft.Authorization/roleAssignments@2022-04-01"
  parent_id = azapi_resource.vault_replica.id
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

data "azapi_resource_list" "kvCertificatesUserRoleDefinition_replica" {
  type      = "Microsoft.Authorization/roleDefinitions@2022-04-01"
  parent_id = azapi_resource.vault_replica.id
  response_export_values = {
    definition_id = "value[?properties.roleName == 'Key Vault Crypto Service Encryption User'].id | [0]"
  }
}

resource "azapi_resource" "kvRoleAssignmentIdentity_replica" {
  type      = "Microsoft.Authorization/roleAssignments@2022-04-01"
  parent_id = azapi_resource.vault_replica.id
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

data "azapi_resource_id" "key_replica" {
  type      = "Microsoft.KeyVault/vaults/keys@2023-02-01"
  parent_id = azapi_resource.vault_replica.id
  name      = var.resource_name
}

resource "azapi_resource_action" "key_replica" {
  type        = "Microsoft.KeyVault/vaults/keys@2023-02-01"
  resource_id = data.azapi_resource_id.key_replica.id
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
    azapi_resource.kvRoleAssignmentTf_replica,
    azapi_resource.kvRoleAssignmentIdentity_replica,
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
        userName = var.mongo_admin_username
      }
      authConfig = {
        allowedModes = ["MicrosoftEntraID", "NativeAuth"]
      }
      compute = {
        tier = "M30"
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
      highAvailability = {
        targetMode = "Disabled"
      }
      previewFeatures = [
        "ShardRebalancer"
      ]
      publicNetworkAccess = "Enabled"
      serverVersion       = "5.0"
      sharding = {
        shardCount = 1
      }
      storage = {
        sizeGb = 32
      }
    }
  }
  sensitive_body = {
    properties = {
      administrator = {
        password = var.mongo_admin_password
      }
    }
  }
  tags = {
    Environment = "Test"
  }
  schema_validation_enabled = false
  response_export_values    = ["*"]
}

# Wait for the MongoDB cluster to have backup capability available
resource "time_sleep" "wait_for_backup_ready" {
  depends_on = [azapi_resource.mongoCluster]

  # Wait 5 minutes for backup to be available - MongoDB clusters typically need time to enable backup
  create_duration = "300s"
}

# Data source to get the updated cluster info with backup details
data "azapi_resource" "mongoCluster_backup_check" {
  type        = "Microsoft.DocumentDB/mongoClusters@2025-09-01"
  resource_id = azapi_resource.mongoCluster.id

  response_export_values = ["properties.backup.earliestRestoreTime"]
  depends_on             = [time_sleep.wait_for_backup_ready]
}

resource "azapi_resource" "mongoCluster_PointInTimeRestore" {
  type      = "Microsoft.DocumentDB/mongoClusters@2025-09-01"
  parent_id = azapi_resource.resourceGroup.id
  name      = "${var.resource_name}-restore"
  location  = var.location
  identity {
    type         = "UserAssigned"
    identity_ids = [azapi_resource.userAssignedIdentity.id]
  }
  body = {
    properties = {
      createMode = "PointInTimeRestore"
      administrator = {
        userName = var.mongo_admin_username
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
      restoreParameters = {
        pointInTimeUTC   = data.azapi_resource.mongoCluster_backup_check.output.properties.backup.earliestRestoreTime
        sourceResourceId = azapi_resource.mongoCluster.id
      }
    }
  }
  sensitive_body = {
    properties = {
      administrator = {
        password = var.mongo_restore_admin_password
      }
    }
  }
  schema_validation_enabled = false
  ignore_casing             = false
  ignore_missing_property   = false

  depends_on = [
    data.azapi_resource.mongoCluster_backup_check
  ]
}

resource "azapi_resource" "mongoCluster_GeoReplica" {
  type      = "Microsoft.DocumentDB/mongoClusters@2025-09-01"
  parent_id = azapi_resource.resourceGroup.id
  name      = "${var.resource_name}-repl"
  location  = var.replica_location
  identity {
    type         = "UserAssigned"
    identity_ids = [azapi_resource.userAssignedIdentity.id]
  }
  body = {
    properties = {
      createMode = "GeoReplica"
      encryption = {
        customerManagedKeyEncryption = {
          keyEncryptionKeyIdentity = {
            identityType                   = "UserAssignedIdentity"
            userAssignedIdentityResourceId = azapi_resource.userAssignedIdentity.id
          }
          keyEncryptionKeyUrl = azapi_resource_action.key_replica.output.properties.keyUri
        }
      }
      replicaParameters = {
        sourceLocation   = var.location # Source location matches the primary cluster
        sourceResourceId = azapi_resource.mongoCluster.id
      }
    }
  }
  schema_validation_enabled = false
  response_export_values    = ["*"]
}
