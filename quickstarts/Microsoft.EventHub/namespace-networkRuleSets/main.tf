terraform {
  required_providers {
    azapi = {
      source = "azure/azapi"
    }
  }
}

provider "azurerm" {
  features {}
}

provider "azapi" {
}

resource "azurerm_resource_group" "test" {
  name     = "myResourceGroup"
  location = "westus"
}

resource "azurerm_virtual_network" "test" {
  name                = "myvnet"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  address_space       = ["172.17.0.0/16"]
  dns_servers         = ["10.0.0.4", "10.0.0.5"]
}

resource "azurerm_subnet" "test" {
  name                 = "default"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefixes     = ["172.17.0.0/24"]

  service_endpoints = ["Microsoft.EventHub"]
}

resource "azurerm_eventhub_namespace" "test" {
  name                = "myNamespace"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku                 = "Standard"
  capacity            = 2
}

resource "azapi_update_resource" "test" {
  type      = "Microsoft.EventHub/namespaces/networkRuleSets@2021-11-01"
  name      = "default"
  parent_id = azurerm_eventhub_namespace.test.id

  body = {
    properties = {
      defaultAction       = "Deny"
      publicNetworkAccess = "Enabled"
      virtualNetworkRules = [
        {
          ignoreMissingVnetServiceEndpoint = false
          subnet = {
            // API bug, returned id replaced `resourceGroups` with `resourcegroups`
            id = replace(azurerm_subnet.test.id, "resourceGroups", "resourcegroups")
          }
        }
      ]
      ipRules = [
        {
          action = "Allow"
          ipMask = "1.1.1.1"
        }
      ]
    }
  }

}