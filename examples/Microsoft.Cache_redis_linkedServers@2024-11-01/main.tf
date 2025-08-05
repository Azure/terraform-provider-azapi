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

resource "azapi_resource" "resourceGroup_secondary" {
  type     = "Microsoft.Resources/resourceGroups@2020-06-01"
  name     = "${var.resource_name}-secondary"
  location = var.location
}

resource "azapi_resource" "redis_secondary" {
  type      = "Microsoft.Cache/redis@2024-11-01"
  parent_id = azapi_resource.resourceGroup.id
  name      = "${var.resource_name}-secondary"
  location  = var.location
  body = {
    properties = {
      disableAccessKeyAuthentication = false
      enableNonSslPort               = false
      minimumTlsVersion              = "1.2"
      publicNetworkAccess            = "Enabled"
      redisConfiguration = {
        maxmemory-delta                        = "642"
        maxmemory-policy                       = "allkeys-lru"
        maxmemory-reserved                     = "642"
        preferred-data-persistence-auth-method = ""
      }
      redisVersion = "6"
      sku = {
        capacity = 1
        family   = "P"
        name     = "Premium"
      }
    }
  }
}

resource "azapi_resource" "redis_primary" {
  type      = "Microsoft.Cache/redis@2024-11-01"
  parent_id = azapi_resource.resourceGroup_secondary.id
  name      = "${var.resource_name}-primary"
  location  = var.location
  body = {
    properties = {
      disableAccessKeyAuthentication = false
      enableNonSslPort               = false
      minimumTlsVersion              = "1.2"
      publicNetworkAccess            = "Enabled"
      redisConfiguration = {
        maxmemory-delta                        = "642"
        maxmemory-policy                       = "allkeys-lru"
        maxmemory-reserved                     = "642"
        preferred-data-persistence-auth-method = ""
      }
      redisVersion = "6"
      sku = {
        capacity = 1
        family   = "P"
        name     = "Premium"
      }
    }
  }
}

resource "azapi_resource" "linkedServer" {
  type      = "Microsoft.Cache/redis/linkedServers@2024-11-01"
  parent_id = azapi_resource.redis_primary.id
  name      = "${var.resource_name}-secondary"
  body = {
    properties = {
      linkedRedisCacheId       = azapi_resource.redis_secondary.id
      linkedRedisCacheLocation = var.location
      serverRole               = "Secondary"
    }
  }
}

