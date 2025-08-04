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

resource "azapi_resource" "registry" {
  type      = "Microsoft.ContainerRegistry/registries@2023-11-01-preview"
  parent_id = azapi_resource.resourceGroup.id
  name      = "${var.resource_name}registry"
  location  = var.location
  body = {
    properties = {
      adminUserEnabled         = false
      anonymousPullEnabled     = false
      dataEndpointEnabled      = false
      networkRuleBypassOptions = "AzureServices"
      policies = {
        exportPolicy = {
          status = "enabled"
        }
        quarantinePolicy = {
          status = "disabled"
        }
        retentionPolicy = {}
        trustPolicy     = {}
      }
      publicNetworkAccess = "Enabled"
      zoneRedundancy      = "Disabled"
    }
    sku = {
      name = "Basic"
    }
  }
}

resource "azapi_resource" "cacheRule" {
  type      = "Microsoft.ContainerRegistry/registries/cacheRules@2023-07-01"
  parent_id = azapi_resource.registry.id
  name      = "${var.resource_name}-cache-rule"
  body = {
    properties = {
      sourceRepository = "mcr.microsoft.com/hello-world"
      targetRepository = "target"
    }
  }
}
