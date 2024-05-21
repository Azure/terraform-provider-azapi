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
  default = "westeurope"
}

resource "azapi_resource" "resourceGroup" {
  type     = "Microsoft.Resources/resourceGroups@2020-06-01"
  name     = var.resource_name
  location = var.location
}

resource "azapi_resource" "databaseAccount" {
  type      = "Microsoft.DocumentDB/databaseAccounts@2021-10-15"
  parent_id = azapi_resource.resourceGroup.id
  name      = var.resource_name
  location  = var.location
  body = {
    kind = "GlobalDocumentDB"
    properties = {
      capabilities = [
        {
          name = "EnableGremlin"
        },
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
      enableFreeTier                     = false
      enableMultipleWriteLocations       = false
      ipRules = [
      ]
      isVirtualNetworkFilterEnabled = false
      locations = [
        {
          failoverPriority = 0
          isZoneRedundant  = false
          locationName     = "West Europe"
        },
      ]
      networkAclBypass = "None"
      networkAclBypassResourceIds = [
      ]
      publicNetworkAccess = "Enabled"
      virtualNetworkRules = [
      ]
    }
  }
  schema_validation_enabled = false
  response_export_values    = ["*"]
}

resource "azapi_resource" "gremlinDatabase" {
  type      = "Microsoft.DocumentDB/databaseAccounts/gremlinDatabases@2023-04-15"
  parent_id = azapi_resource.databaseAccount.id
  name      = var.resource_name
  body = {
    properties = {
      options = {
      }
      resource = {
        id = var.resource_name
      }
    }
  }
  schema_validation_enabled = false
  response_export_values    = ["*"]
}

resource "azapi_resource" "graph" {
  type      = "Microsoft.DocumentDB/databaseAccounts/gremlinDatabases/graphs@2023-04-15"
  parent_id = azapi_resource.gremlinDatabase.id
  name      = var.resource_name
  body = {
    properties = {
      options = {
        throughput = 400
      }
      resource = {
        id = var.resource_name
        partitionKey = {
          kind = "Hash"
          paths = [
            "/test",
          ]
        }
      }
    }
  }
  schema_validation_enabled = false
  response_export_values    = ["*"]
}

