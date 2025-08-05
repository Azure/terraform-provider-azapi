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

resource "azapi_resource" "redis" {
  type      = "Microsoft.Cache/redis@2024-11-01"
  parent_id = azapi_resource.resourceGroup.id
  name      = var.resource_name
  location  = var.location
  body = {
    properties = {
      disableAccessKeyAuthentication = false
      enableNonSslPort               = true
      minimumTlsVersion              = "1.2"
      publicNetworkAccess            = "Enabled"
      redisConfiguration = {
        maxmemory-policy                       = "volatile-lru"
        preferred-data-persistence-auth-method = ""
      }
      redisVersion = "6"
      sku = {
        capacity = 1
        family   = "C"
        name     = "Basic"
      }
    }
  }
}

resource "azapi_resource" "accessPolicy" {
  type      = "Microsoft.Cache/redis/accessPolicies@2024-11-01"
  parent_id = azapi_resource.redis.id
  name      = "${var.resource_name}-accessPolicy"
  body = {
    properties = {
      permissions = "+@read +@connection +cluster|info allkeys"
    }
  }
}

