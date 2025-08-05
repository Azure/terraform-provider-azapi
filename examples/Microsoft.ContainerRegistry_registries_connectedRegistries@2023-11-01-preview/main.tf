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
      dataEndpointEnabled      = true
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
      name = "Premium"
    }
  }
}

resource "azapi_resource" "scopeMap" {
  type      = "Microsoft.ContainerRegistry/registries/scopeMaps@2023-11-01-preview"
  parent_id = azapi_resource.registry.id
  name      = "${var.resource_name}scopemap"
  body = {
    properties = {
      actions     = ["repositories/hello-world/content/delete", "repositories/hello-world/content/read", "repositories/hello-world/content/write", "repositories/hello-world/metadata/read", "repositories/hello-world/metadata/write", "gateway/${var.resource_name}connectedregistry/config/read", "gateway/${var.resource_name}connectedregistry/config/write", "gateway/${var.resource_name}connectedregistry/message/read", "gateway/${var.resource_name}connectedregistry/message/write"]
      description = ""
    }
  }
}

resource "azapi_resource" "token" {
  type      = "Microsoft.ContainerRegistry/registries/tokens@2023-11-01-preview"
  parent_id = azapi_resource.registry.id
  name      = "${var.resource_name}token"
  body = {
    properties = {
      scopeMapId = azapi_resource.scopeMap.id
      status     = "enabled"
    }
  }
}

resource "azapi_resource" "connectedRegistry" {
  type      = "Microsoft.ContainerRegistry/registries/connectedRegistries@2023-11-01-preview"
  parent_id = azapi_resource.registry.id
  name      = "${var.resource_name}connectedregistry"
  body = {
    properties = {
      clientTokenIds = null
      logging = {
        auditLogStatus = "Disabled"
        logLevel       = "None"
      }
      mode = "ReadWrite"
      parent = {
        syncProperties = {
          messageTtl = "P1D"
          schedule   = "* * * * *"
          syncWindow = ""
          tokenId    = azapi_resource.token.id
        }
      }
    }
  }
}