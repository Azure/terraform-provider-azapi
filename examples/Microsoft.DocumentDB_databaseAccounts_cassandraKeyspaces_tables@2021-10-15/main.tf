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
