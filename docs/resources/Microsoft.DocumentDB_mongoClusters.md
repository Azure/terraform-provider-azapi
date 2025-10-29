---
subcategory: "Microsoft.DocumentDB - Azure Cosmos DB"
page_title: "mongoClusters"
description: |-
  Manages a MongoDB Cluster using vCore Architecture.
---

# Microsoft.DocumentDB/mongoClusters - MongoDB Cluster using vCore Architecture

This article demonstrates how to use `azapi` provider to manage the MongoDB Cluster using vCore Architecture resource in Azure.



## Example Usage

### default

```hcl
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
        password = "passworD!@#!$*"
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

```



## Arguments Reference

The following arguments are supported:

* `type` - (Required) The type of the resource. This should be set to `Microsoft.DocumentDB/mongoClusters@api-version`. The available api-versions for this resource are: [`2023-03-01-preview`, `2023-03-15-preview`, `2023-09-15-preview`, `2023-11-15-preview`, `2024-02-15-preview`, `2024-03-01-preview`, `2024-06-01-preview`, `2024-07-01`, `2024-10-01-preview`, `2025-04-01-preview`, `2025-07-01-preview`].

* `parent_id` - (Required) The ID of the azure resource in which this resource is created. The allowed values are:  
  `/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}`

* `name` - (Required) Specifies the name of the azure resource. Changing this forces a new resource to be created.

* `body` - (Required) Specifies the configuration of the resource. More information about the arguments in `body` can be found in the [Microsoft documentation](https://learn.microsoft.com/en-us/azure/templates/Microsoft.DocumentDB/mongoClusters?pivots=deployment-language-terraform).

For other arguments, please refer to the [azapi_resource](https://registry.terraform.io/providers/Azure/azapi/latest/docs/resources/resource) documentation.

## Import

 ```shell
 # Azure resource can be imported using the resource id, e.g.
 terraform import azapi_resource.example /subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.DocumentDB/mongoClusters/{resourceName}
 
 # It also supports specifying API version by using the resource id with api-version as a query parameter, e.g.
 terraform import azapi_resource.example /subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.DocumentDB/mongoClusters/{resourceName}?api-version=2025-07-01-preview
 ```
