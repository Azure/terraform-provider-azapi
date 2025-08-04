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

resource "azapi_resource" "resourceGroup" {
  type     = "Microsoft.Resources/resourceGroups@2020-06-01"
  name     = var.resource_name
  location = var.location
}

resource "azapi_resource" "redisEnterprise" {
  type      = "Microsoft.Cache/redisEnterprise@2024-10-01"
  parent_id = azapi_resource.resourceGroup.id
  name      = var.resource_name
  location  = var.location
  body = {
    properties = {
      minimumTlsVersion = "1.2"
    }
    sku = {
      capacity = 4
      name     = "Enterprise_E20"
    }
  }
}

resource "azapi_resource" "databas" {
  type      = "Microsoft.Cache/redisEnterprise/databases@2024-10-01"
  parent_id = azapi_resource.redisEnterprise.id
  name      = "default"
  body = {
    properties = {
      clientProtocol   = "Encrypted"
      clusteringPolicy = "OSSCluster"
      evictionPolicy   = "VolatileLRU"
      modules          = []
      port             = 10000
    }
  }
}

