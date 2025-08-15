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
