---
subcategory: "Microsoft.Storage - Storage"
page_title: "storageAccounts/encryptionScopes"
description: |-
  Manages a Storage Encryption Scope.
---

# Microsoft.Storage/storageAccounts/encryptionScopes - Storage Encryption Scope

This article demonstrates how to use `azapi` provider to manage the Storage Encryption Scope resource in Azure.

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


```



## Arguments Reference

The following arguments are supported:

* `type` - (Required) The type of the resource. This should be set to `Microsoft.Storage/storageAccounts/encryptionScopes@api-version`. The available api-versions for this resource are: [`2019-06-01`, `2020-08-01-preview`, `2021-01-01`, `2021-02-01`, `2021-04-01`, `2021-06-01`, `2021-08-01`, `2021-09-01`, `2022-05-01`, `2022-09-01`, `2023-01-01`, `2023-04-01`, `2023-05-01`, `2024-01-01`, `2025-01-01`].

* `parent_id` - (Required) The ID of the azure resource in which this resource is created. The allowed values are:  
  `/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Storage/storageAccounts/{resourceName}`

* `name` - (Required) Specifies the name of the azure resource. Changing this forces a new resource to be created.

* `body` - (Required) Specifies the configuration of the resource. More information about the arguments in `body` can be found in the [Microsoft documentation](https://learn.microsoft.com/en-us/azure/templates/Microsoft.Storage/storageAccounts/encryptionScopes?pivots=deployment-language-terraform).

For other arguments, please refer to the [azapi_resource](https://registry.terraform.io/providers/Azure/azapi/latest/docs/resources/resource) documentation.

## Import

 ```shell
 # Azure resource can be imported using the resource id, e.g.
 terraform import azapi_resource.example /subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Storage/storageAccounts/{resourceName}/encryptionScopes/{resourceName}
 
 # It also supports specifying API version by using the resource id with api-version as a query parameter, e.g.
 terraform import azapi_resource.example /subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Storage/storageAccounts/{resourceName}/encryptionScopes/{resourceName}?api-version=2025-01-01
 ```
