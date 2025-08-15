---
subcategory: "Microsoft.RecoveryServices - Azure Site Recovery"
page_title: "vaults/backupFabrics/protectionContainers"
description: |-
  Manages a storage account container in an Azure Recovery Vault.
---

# Microsoft.RecoveryServices/vaults/backupFabrics/protectionContainers - storage account container in an Azure Recovery Vault

This article demonstrates how to use `azapi` provider to manage the storage account container in an Azure Recovery Vault resource in Azure.

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
  default = "westus"
}

locals {
  # Storage Account name must be 3-24 lowercase alphanumeric
  sa_name    = substr(join("", regexall("[a-z0-9]", lower(var.resource_name))), 0, 24)
  vault_name = "${var.resource_name}-vault"
}

resource "azapi_resource" "resourceGroup" {
  type     = "Microsoft.Resources/resourceGroups@2020-06-01"
  name     = var.resource_name
  location = var.location
}

resource "azapi_resource" "vault" {
  type      = "Microsoft.RecoveryServices/vaults@2024-01-01"
  parent_id = azapi_resource.resourceGroup.id
  name      = local.vault_name
  location  = var.location
  body = {
    properties = {
      publicNetworkAccess = "Enabled"
      redundancySettings = {
        crossRegionRestore            = "Disabled"
        standardTierStorageRedundancy = "GeoRedundant"
      }
    }
    sku = {
      name = "Standard"
    }
  }
}

resource "azapi_resource" "storageAccount" {
  type      = "Microsoft.Storage/storageAccounts@2023-05-01"
  parent_id = azapi_resource.resourceGroup.id
  name      = local.sa_name
  location  = var.location
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

resource "azapi_resource" "protectionContainer" {
  type      = "Microsoft.RecoveryServices/vaults/backupFabrics/protectionContainers@2023-02-01"
  parent_id = "${azapi_resource.vault.id}/backupFabrics/Azure"
  # Format: "StorageContainer;storage;<resourceGroupName>;<storageAccountName>"
  name = "StorageContainer;storage;${var.resource_name};${local.sa_name}"
  body = {
    properties = {
      backupManagementType = "AzureStorage"
      containerType        = "StorageContainer"
      friendlyName         = local.sa_name
      sourceResourceId     = azapi_resource.storageAccount.id
    }
  }
}

```



## Arguments Reference

The following arguments are supported:

* `type` - (Required) The type of the resource. This should be set to `Microsoft.RecoveryServices/vaults/backupFabrics/protectionContainers@api-version`. The available api-versions for this resource are: [`2016-06-01`, `2016-12-01`, `2020-10-01`, `2020-12-01`, `2021-01-01`, `2021-02-01`, `2021-02-01-preview`, `2021-02-10`, `2021-03-01`, `2021-04-01`, `2021-06-01`, `2021-07-01`, `2021-08-01`, `2021-10-01`, `2021-12-01`, `2022-01-01`, `2022-02-01`, `2022-03-01`, `2022-04-01`, `2022-06-01-preview`, `2022-09-01-preview`, `2022-09-30-preview`, `2022-10-01`, `2023-01-01`, `2023-02-01`, `2023-04-01`, `2023-06-01`, `2023-08-01`, `2024-01-01`, `2024-02-01`, `2024-04-01`, `2024-04-30-preview`, `2024-07-30-preview`, `2024-10-01`, `2024-11-01-preview`, `2025-01-01`, `2025-02-01`, `2025-02-28-preview`].

* `parent_id` - (Required) The ID of the azure resource in which this resource is created. The allowed values are:  
  `{any azure resource id}`

* `name` - (Required) Specifies the name of the azure resource. Changing this forces a new resource to be created.

* `body` - (Required) Specifies the configuration of the resource. More information about the arguments in `body` can be found in the [Microsoft documentation](https://learn.microsoft.com/en-us/azure/templates/Microsoft.RecoveryServices/vaults/backupFabrics/protectionContainers?pivots=deployment-language-terraform).

For other arguments, please refer to the [azapi_resource](https://registry.terraform.io/providers/Azure/azapi/latest/docs/resources/resource) documentation.

## Import

 ```shell
 # Azure resource can be imported using the resource id, e.g.
 terraform import azapi_resource.example {any azure resource id}/protectionContainers/{resourceName}
 
 # It also supports specifying API version by using the resource id with api-version as a query parameter, e.g.
 terraform import azapi_resource.example {any azure resource id}/protectionContainers/{resourceName}?api-version=2025-02-28-preview
 ```
