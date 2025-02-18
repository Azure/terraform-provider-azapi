terraform {
  required_providers {
    azapi = {
      source = "Azure/azapi"
    }
  }
}

provider "azurerm" {
  features {
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

resource "azapi_resource" "virtualWan" {
  type      = "Microsoft.Network/virtualWans@2022-07-01"
  parent_id = azapi_resource.resourceGroup.id
  name      = var.resource_name
  location  = var.location
  body = {
    properties = {
      allowBranchToBranchTraffic     = true
      disableVpnEncryption           = false
      office365LocalBreakoutCategory = "None"
      type                           = "Standard"
    }
  }
  schema_validation_enabled = false
  response_export_values    = ["*"]
}

resource "azapi_resource" "virtualHub" {
  type      = "Microsoft.Network/virtualHubs@2022-07-01"
  parent_id = azapi_resource.resourceGroup.id
  name      = var.resource_name
  location  = var.location
  body = {
    properties = {
      addressPrefix        = "10.0.2.0/24"
      hubRoutingPreference = "ExpressRoute"
      virtualRouterAutoScaleConfiguration = {
        minCapacity = 2
      }
      virtualWan = {
        id = azapi_resource.virtualWan.id
      }
    }
  }
  schema_validation_enabled = false
  response_export_values    = ["*"]
}

resource "azurerm_firewall" "test" {
  name                = var.resource_name
  location            = azapi_resource.resourceGroup.location
  resource_group_name = azapi_resource.resourceGroup.name
  sku_name            = "AZFW_Hub"
  sku_tier            = "Standard"

  virtual_hub {
    virtual_hub_id  = azapi_resource.virtualHub.id
    public_ip_count = 1
  }
}

resource "azapi_resource" "routingIntent" {
  name      = var.resource_name
  type      = "Microsoft.Network/virtualHubs/routingIntent@2022-09-01"
  parent_id = azapi_resource.virtualHub.id

  body = {
    properties = {
      routingPolicies = [
        {
          name = "InternetTraffic"
          destinations = [
            "Internet"
          ]
          nextHop = azurerm_firewall.test.id
        },
        {
          name = "PrivateTrafficPolicy"
          destinations = [
            "PrivateTraffic"
          ]
          nextHop = azurerm_firewall.test.id
        }
      ]
    }
  }
}
