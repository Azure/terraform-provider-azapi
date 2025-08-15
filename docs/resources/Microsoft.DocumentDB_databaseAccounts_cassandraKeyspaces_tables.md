---
subcategory: "Microsoft.DocumentDB - Azure Cosmos DB"
page_title: "databaseAccounts/cassandraKeyspaces/tables"
description: |-
  Manages a Cassandra Table within a Cosmos DB Cassandra Keyspace.
---

# Microsoft.DocumentDB/databaseAccounts/cassandraKeyspaces/tables - Cassandra Table within a Cosmos DB Cassandra Keyspace

This article demonstrates how to use `azapi` provider to manage the Cassandra Table within a Cosmos DB Cassandra Keyspace resource in Azure.

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
  # This base name will be used to derive all other resource names.
}

variable "location" {
  type    = string
  default = "eastus"
}

locals {
  # Cosmos DB account names must be globally unique, 3-44 chars, lowercase letters or numbers only.
  # Derive names from the base resource_name; users can override via terraform.tfvars to ensure uniqueness.
  account_name  = lower(replace(var.resource_name, "-", ""))
  keyspace_name = "${lower(var.resource_name)}ks"
  table_name    = "${lower(var.resource_name)}tbl"
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
    kind = "GlobalDocumentDB"
    properties = {
      backupPolicy = null
      capabilities = [{
        name = "EnableCassandra"
      }]
      consistencyPolicy = {
        defaultConsistencyLevel = "Strong"
        maxIntervalInSeconds    = 5
        maxStalenessPrefix      = 100
      }
      databaseAccountOfferType           = "Standard"
      defaultIdentity                    = "FirstPartyIdentity"
      disableKeyBasedMetadataWriteAccess = false
      disableLocalAuth                   = false
      enableAnalyticalStorage            = true
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

resource "azapi_resource" "cassandraKeyspace" {
  type      = "Microsoft.DocumentDB/databaseAccounts/cassandraKeyspaces@2021-10-15"
  parent_id = azapi_resource.databaseAccount.id
  name      = local.keyspace_name
  body = {
    properties = {
      options = {}
      resource = {
        id = local.keyspace_name
      }
    }
  }
}

resource "azapi_resource" "table" {
  type      = "Microsoft.DocumentDB/databaseAccounts/cassandraKeyspaces/tables@2021-10-15"
  parent_id = azapi_resource.cassandraKeyspace.id
  name      = local.table_name
  body = {
    properties = {
      options = {}
      resource = {
        analyticalStorageTtl = 1
        id                   = local.table_name
        schema = {
          clusterKeys = []
          columns = [{
            name = "test1"
            type = "ascii"
            }, {
            name = "test2"
            type = "int"
          }]
          partitionKeys = [{
            name = "test1"
          }]
        }
      }
    }
  }
}

```



## Arguments Reference

The following arguments are supported:

* `type` - (Required) The type of the resource. This should be set to `Microsoft.DocumentDB/databaseAccounts/cassandraKeyspaces/tables@api-version`. The available api-versions for this resource are: [`2019-08-01`, `2019-12-12`, `2020-03-01`, `2020-04-01`, `2020-06-01-preview`, `2020-09-01`, `2021-01-15`, `2021-03-01-preview`, `2021-03-15`, `2021-04-01-preview`, `2021-04-15`, `2021-05-15`, `2021-06-15`, `2021-07-01-preview`, `2021-10-15`, `2021-10-15-preview`, `2021-11-15-preview`, `2022-02-15-preview`, `2022-05-15`, `2022-05-15-preview`, `2022-08-15`, `2022-08-15-preview`, `2022-11-15`, `2022-11-15-preview`, `2023-03-01-preview`, `2023-03-15`, `2023-03-15-preview`, `2023-04-15`, `2023-09-15`, `2023-09-15-preview`, `2023-11-15`, `2023-11-15-preview`, `2024-02-15-preview`, `2024-05-15`, `2024-05-15-preview`, `2024-08-15`, `2024-09-01-preview`, `2024-11-15`, `2024-12-01-preview`, `2025-04-15`, `2025-05-01-preview`].

* `parent_id` - (Required) The ID of the azure resource in which this resource is created. The allowed values are:  
  `/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.DocumentDB/databaseAccounts/{resourceName}/cassandraKeyspaces/{resourceName}`

* `name` - (Required) Specifies the name of the azure resource. Changing this forces a new resource to be created.

* `body` - (Required) Specifies the configuration of the resource. More information about the arguments in `body` can be found in the [Microsoft documentation](https://learn.microsoft.com/en-us/azure/templates/Microsoft.DocumentDB/databaseAccounts/cassandraKeyspaces/tables?pivots=deployment-language-terraform).

For other arguments, please refer to the [azapi_resource](https://registry.terraform.io/providers/Azure/azapi/latest/docs/resources/resource) documentation.

## Import

 ```shell
 # Azure resource can be imported using the resource id, e.g.
 terraform import azapi_resource.example /subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.DocumentDB/databaseAccounts/{resourceName}/cassandraKeyspaces/{resourceName}/tables/{resourceName}
 
 # It also supports specifying API version by using the resource id with api-version as a query parameter, e.g.
 terraform import azapi_resource.example /subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.DocumentDB/databaseAccounts/{resourceName}/cassandraKeyspaces/{resourceName}/tables/{resourceName}?api-version=2025-05-01-preview
 ```
