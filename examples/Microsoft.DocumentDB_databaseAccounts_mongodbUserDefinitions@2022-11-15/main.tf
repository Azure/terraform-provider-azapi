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

variable "mongodb_user_password" {
  type        = string
  description = "The password for the MongoDB user"
  sensitive   = true
}

resource "azapi_resource" "resourceGroup" {
  type     = "Microsoft.Resources/resourceGroups@2020-06-01"
  name     = var.resource_name
  location = var.location
}

resource "azapi_resource" "databaseAccount" {
  type      = "Microsoft.DocumentDB/databaseAccounts@2024-08-15"
  parent_id = azapi_resource.resourceGroup.id
  name      = "${var.resource_name}-acct"
  location  = var.location
  body = {
    kind = "MongoDB"
    properties = {
      backupPolicy = null
      capabilities = [{
        name = "EnableMongoRoleBasedAccessControl"
        }, {
        name = "EnableMongo"
      }]
      consistencyPolicy = {
        defaultConsistencyLevel = "Strong"
        maxIntervalInSeconds    = 5
        maxStalenessPrefix      = 100
      }
      databaseAccountOfferType           = "Standard"
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

resource "azapi_resource" "mongodbDatabas" {
  type      = "Microsoft.DocumentDB/databaseAccounts/mongodbDatabases@2021-10-15"
  parent_id = azapi_resource.databaseAccount.id
  name      = "${var.resource_name}-mongodb"
  body = {
    properties = {
      options = {}
      resource = {
        id = "${var.resource_name}-mongodb"
      }
    }
  }
}

resource "azapi_resource" "mongodbUserDefinition" {
  type      = "Microsoft.DocumentDB/databaseAccounts/mongodbUserDefinitions@2022-11-15"
  parent_id = azapi_resource.databaseAccount.id
  name      = "${azapi_resource.mongodbDatabas.name}.myUserName"
  body = {
    properties = {
      databaseName = azapi_resource.mongodbDatabas.name
      mechanisms   = "SCRAM-SHA-256"
      password     = var.mongodb_user_password
      userName     = "myUserName"
    }
  }
}

