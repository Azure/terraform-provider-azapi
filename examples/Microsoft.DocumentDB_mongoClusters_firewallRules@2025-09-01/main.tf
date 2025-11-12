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

data "azapi_client_config" "current" {}

resource "azapi_resource" "mongoCluster" {
  type      = "Microsoft.DocumentDB/mongoClusters@2025-09-01"
  parent_id = azapi_resource.resourceGroup.id
  name      = var.resource_name
  location  = var.location
  body = {
    properties = {
      authConfig = {
        allowedModes = ["MicrosoftEntraID"]
      }
      compute = {
        tier = "M40"
      }
      highAvailability = {
        targetMode = "Disabled"
      }
      previewFeatures = [
        "ShardRebalancer"
      ]
      publicNetworkAccess = "Enabled"
      serverVersion       = "5.0"
      sharding = {
        shardCount = 1
      }
      storage = {
        sizeGb = 32
      }
    }
  }
  tags = {
    Environment = "Test"
  }
  schema_validation_enabled = false
  response_export_values    = ["*"]
}

resource "azapi_resource" "firewallRule" {
  type      = "Microsoft.DocumentDB/mongoClusters/firewallRules@2025-09-01"
  parent_id = azapi_resource.mongoCluster.id
  name      = var.resource_name
  body = {
    properties = {
      endIpAddress   = "0.0.0.0"
      startIpAddress = "0.0.0.0"
    }
  }
  schema_validation_enabled = false
  response_export_values    = ["*"]
}
