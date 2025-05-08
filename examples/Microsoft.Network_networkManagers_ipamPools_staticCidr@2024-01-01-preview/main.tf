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

data "azapi_client_config" "current" {}

resource "azapi_resource" "resourceGroup" {
  type     = "Microsoft.Resources/resourceGroups@2020-06-01"
  name     = var.resource_name
  location = var.location
}

resource "azapi_resource" "networkManager" {
  type      = "Microsoft.Network/networkManagers@2022-09-01"
  parent_id = azapi_resource.resourceGroup.id
  name      = var.resource_name
  location  = var.location
  body = {
    properties = {
      description = ""
      networkManagerScopeAccesses = [
        "SecurityAdmin",
      ]
      networkManagerScopes = {
        managementGroups = [
        ]
        subscriptions = [
          "/subscriptions/${data.azapi_client_config.current.subscription_id}",
        ]
      }
    }
  }
  retry = {
    error_message_regex = ["CannotDeleteResource"]
  }
}

resource "azapi_resource" "ipamPool" {
  type      = "Microsoft.Network/networkManagers/ipamPools@2024-01-01-preview"
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
  }

  tags = {
    myTag = "testTag"
  }
}

resource "azapi_resource" "staticCidr" {
  type      = "Microsoft.Network/networkManagers/ipamPools/staticCidrs@2024-01-01-preview"
  parent_id = azapi_resource.ipamPool.id
  name      = var.resource_name
  body = {
    properties = {
      addressPrefixes = [
        "10.0.0.0/25",
      ]
      numberOfIPAddressesToAllocate = ""
      description                   = "test description"
    }
  }
}
