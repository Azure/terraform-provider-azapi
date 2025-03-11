terraform {
  required_providers {
    azapi = {
      source = "Azure/azapi"
    }
    azurerm = {
      source = "Hashicorp/azurerm"
    }
    time = {
      source = "Hashicorp/time"
    }
  }
}

provider "azurerm" {
  features {}
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

data "azurerm_client_config" "current" {
}

data "azapi_resource" "subscription" {
  type                   = "Microsoft.Resources/subscriptions@2021-01-01"
  resource_id            = "/subscriptions/${data.azurerm_client_config.current.subscription_id}"
  response_export_values = ["*"]
}

resource "azapi_resource" "networkManager" {
  type      = "Microsoft.Network/networkManagers@2024-05-01"
  parent_id = azapi_resource.resourceGroup.id
  name      = var.resource_name
  location  = var.location
  body = {
    properties = {
      description                 = ""
      networkManagerScopeAccesses = []
      networkManagerScopes = {
        managementGroups = []
        subscriptions = [
          data.azapi_resource.subscription.id,
        ]
      }
    }
    tags = { "sampleTag" : "sampleTag" }
  }
  schema_validation_enabled = false
  response_export_values    = ["*"]
}

resource "time_sleep" "wait_60_seconds" {
  depends_on = [azapi_resource.networkManager]

  destroy_duration = "60s"
}

resource "azapi_resource" "ipamPool" {
  type      = "Microsoft.Network/networkManagers/ipamPools@2024-05-01"
  parent_id = azapi_resource.networkManager.id
  name      = var.resource_name
  location  = var.location
  body = {
    properties = {
      addressPrefixes = [
        "10.0.0.0/24",
      ]
      description    = "Test description."
      parentPoolName = ""
      displayName    = "testDisplayName"
    }
    tags = { "myTag" : "testTag" }
  }
  schema_validation_enabled = false
  ignore_casing             = false
  ignore_missing_property   = false

  depends_on = [time_sleep.wait_60_seconds]
}

resource "azapi_resource" "vnet_withIPAM" {
  type      = "Microsoft.Network/virtualNetworks@2024-05-01"
  parent_id = azapi_resource.resourceGroup.id
  name      = var.resource_name
  location  = var.location
  body = {
    properties = {
      addressSpace = {
        ipamPoolPrefixAllocations = [
          {
            numberOfIpAddresses = "100",
            pool = {
              id = azapi_resource.ipamPool.id
            }
          }
        ]
      }
    }
  }
  schema_validation_enabled = false
  ignore_casing             = false
  ignore_missing_property   = false

  # ignore that addressPrefixes is created with empty [], and response has address prefix of pool after creation/association of vnet and pool
  lifecycle {
    ignore_changes = [body.properties.addressSpace.addressPrefixes]
  }

  depends_on = [azapi_resource.ipamPool] # This line forces vnet_withIPAM to destroy first
}
