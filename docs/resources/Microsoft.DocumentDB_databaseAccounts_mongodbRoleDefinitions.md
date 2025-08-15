---
subcategory: "Microsoft.DocumentDB - Azure Cosmos DB"
page_title: "databaseAccounts/mongodbRoleDefinitions"
description: |-
  Manages a Cosmos DB Mongo Role Definition.
---

# Microsoft.DocumentDB/databaseAccounts/mongodbRoleDefinitions - Cosmos DB Mongo Role Definition

This article demonstrates how to use `azapi` provider to manage the Cosmos DB Mongo Role Definition resource in Azure.

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
  # Base name used to derive all other resource names. Override in terraform.tfvars for uniqueness.
}

variable "location" {
  type    = string
  default = "eastus"
}

locals {
  # Cosmos DB account names must be globally unique, 3-44 chars, lowercase letters or numbers only.
  account_name = lower(replace(var.resource_name, "-", ""))
  db_name      = "${lower(var.resource_name)}db"
  role_name    = "${lower(var.resource_name)}role"
}

resource "azapi_resource" "resourceGroup" {
  type     = "Microsoft.Resources/resourceGroups@2020-06-01"
  name     = var.resource_name
  location = var.location
}

resource "azapi_resource" "databaseAccount" {
  type      = "Microsoft.DocumentDB/databaseAccounts@2024-08-15"
  parent_id = azapi_resource.resourceGroup.id
  name      = local.account_name
  location  = var.location
  body = {
    kind = "MongoDB"
    properties = {
      backupPolicy = null
      capabilities = [
        { name = "EnableMongoRoleBasedAccessControl" },
        { name = "EnableMongo" }
      ]
      consistencyPolicy = {
        defaultConsistencyLevel = "Strong"
        maxIntervalInSeconds    = 5
        maxStalenessPrefix      = 100
      }
      databaseAccountOfferType           = "Standard"
      defaultIdentity                    = "FirstPartyIdentity"
      disableKeyBasedMetadataWriteAccess = false
      disableLocalAuth                   = false
      enableAnalyticalStorage            = false
      enableAutomaticFailover            = false
      enableBurstCapacity                = false
      enableFreeTier                     = false
      enableMultipleWriteLocations       = false
      enablePartitionMerge               = false
      ipRules                            = []
      isVirtualNetworkFilterEnabled      = false
      locations = [{
        failoverPriority = 0
        isZoneRedundant  = false
        locationName     = var.location
      }]
      minimalTlsVersion           = "Tls12"
      networkAclBypass            = "None"
      networkAclBypassResourceIds = []
      publicNetworkAccess         = "Enabled"
      virtualNetworkRules         = []
    }
  }
}

resource "azapi_resource" "mongodbDatabase" {
  type      = "Microsoft.DocumentDB/databaseAccounts/mongodbDatabases@2021-10-15"
  parent_id = azapi_resource.databaseAccount.id
  name      = local.db_name
  body = {
    properties = {
      options = {}
      resource = {
        id = local.db_name
      }
    }
  }
}

resource "azapi_resource" "mongodbRoleDefinition" {
  type      = "Microsoft.DocumentDB/databaseAccounts/mongodbRoleDefinitions@2022-11-15"
  parent_id = azapi_resource.databaseAccount.id
  name      = "${local.db_name}.${local.role_name}"
  body = {
    properties = {
      databaseName = local.db_name
      roleName     = local.role_name
      type         = 1
    }
  }
  depends_on = [azapi_resource.mongodbDatabase]
}

```



## Arguments Reference

The following arguments are supported:

* `type` - (Required) The type of the resource. This should be set to `Microsoft.DocumentDB/databaseAccounts/mongodbRoleDefinitions@api-version`. The available api-versions for this resource are: [`2021-10-15-preview`, `2021-11-15-preview`, `2022-02-15-preview`, `2022-05-15-preview`, `2022-08-15`, `2022-08-15-preview`, `2022-11-15`, `2022-11-15-preview`, `2023-03-01-preview`, `2023-03-15`, `2023-03-15-preview`, `2023-04-15`, `2023-09-15`, `2023-09-15-preview`, `2023-11-15`, `2023-11-15-preview`, `2024-02-15-preview`, `2024-05-15`, `2024-05-15-preview`, `2024-08-15`, `2024-09-01-preview`, `2024-11-15`, `2024-12-01-preview`, `2025-04-15`, `2025-05-01-preview`].

* `parent_id` - (Required) The ID of the azure resource in which this resource is created. The allowed values are:  
  `/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.DocumentDB/databaseAccounts/{resourceName}`

* `name` - (Required) Specifies the name of the azure resource. Changing this forces a new resource to be created.

* `body` - (Required) Specifies the configuration of the resource. More information about the arguments in `body` can be found in the [Microsoft documentation](https://learn.microsoft.com/en-us/azure/templates/Microsoft.DocumentDB/databaseAccounts/mongodbRoleDefinitions?pivots=deployment-language-terraform).

For other arguments, please refer to the [azapi_resource](https://registry.terraform.io/providers/Azure/azapi/latest/docs/resources/resource) documentation.

## Import

 ```shell
 # Azure resource can be imported using the resource id, e.g.
 terraform import azapi_resource.example /subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.DocumentDB/databaseAccounts/{resourceName}/mongodbRoleDefinitions/{resourceName}
 
 # It also supports specifying API version by using the resource id with api-version as a query parameter, e.g.
 terraform import azapi_resource.example /subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.DocumentDB/databaseAccounts/{resourceName}/mongodbRoleDefinitions/{resourceName}?api-version=2025-05-01-preview
 ```
