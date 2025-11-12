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

resource "azapi_resource" "mongoClusterSSDv2" {
  type      = "Microsoft.DocumentDB/mongoClusters@2025-09-01"
  parent_id = azapi_resource.resourceGroup.id
  name      = "${var.resource_name}-ssdv2"
  location  = var.location
  body = {
    properties = {
      authConfig = {
        allowedModes = ["MicrosoftEntraID"]
      }
      compute = {
        tier = "M30"
      }
      highAvailability = {
        targetMode = "Disabled"
      }
      serverVersion       = "6.0"
      publicNetworkAccess = "Disabled"
      sharding = {
        shardCount = 1
      }
      storage = {
        sizeGb = 64
        type   = "PremiumSSDv2"
      }
    }
  }
  schema_validation_enabled = false
  response_export_values    = ["*"]
}

resource "azapi_resource" "mongoUser_EntraServicePrincipal" {
  type      = "Microsoft.DocumentDB/mongoClusters/users@2025-09-01"
  name      = azapi_resource.userAssignedIdentity.output.properties.principalId
  parent_id = azapi_resource.mongoClusterSSDv2.id

  body = {
    properties = {
      roles = [
        {
          role = "root"
          db   = "admin"
        }
      ]
      identityProvider = {
        type = "MicrosoftEntraID"
        properties = {
          principalType = "ServicePrincipal"
        }
      }
    }
  }
  schema_validation_enabled = false
  response_export_values    = ["*"]
}
