terraform {
  required_providers {
    azapi = {
      source = "Azure/azapi"
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
  location = "australiaeast"
}

resource "azurerm_virtual_network" "vnet" {
  address_space       = ["10.0.0.0/16"]
  location            = azurerm_resource_group.test.location
  name                = "myvnet"
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_subnet" "inbounddnssub" {
  name                 = "inbounddns"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.vnet.name
  address_prefixes     = ["10.0.0.0/28"]

  delegation {
    name = "Microsoft.Network.dnsResolvers"
    service_delegation {
      actions = ["Microsoft.Network/virtualNetworks/subnets/join/action"]
      name    = "Microsoft.Network/dnsResolvers"
    }
  }
}

resource "azurerm_subnet" "outbounddnssub" {
  name                 = "outbounddns"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.vnet.name
  address_prefixes     = ["10.0.0.64/28"]

  delegation {
    name = "Microsoft.Network.dnsResolvers"
    service_delegation {
      actions = ["Microsoft.Network/virtualNetworks/subnets/join/action"]
      name    = "Microsoft.Network/dnsResolvers"
    }
  }
}

resource "azapi_resource" "testresolver" {
  type      = "Microsoft.Network/dnsResolvers@2022-07-01"
  name      = "myresolver"
  parent_id = azurerm_resource_group.test.id
  location  = azurerm_resource_group.test.location

  body = jsonencode({
    properties = {
      virtualNetwork = {
        id = azurerm_virtual_network.vnet.id
      }
    }
  })
}

resource "azapi_resource" "inboundendpoint" {
  type      = "Microsoft.Network/dnsResolvers/inboundEndpoints@2022-07-01"
  name      = "inboundendpoint"
  parent_id = azapi_resource.testresolver.id
  location  = azapi_resource.testresolver.location

  body = jsonencode({
    properties = {
      ipConfigurations = [{ subnet = { id = azurerm_subnet.inbounddnssub.id } }]
    }
  })
}

resource "azapi_resource" "outboundendpoint" {
  type      = "Microsoft.Network/dnsResolvers/outboundEndpoints@2022-07-01"
  name      = "outboundendpoint"
  parent_id = azapi_resource.testresolver.id
  location  = azapi_resource.testresolver.location

  body = jsonencode({
    properties = {
      subnet = {
        id = azurerm_subnet.outbounddnssub.id
      }
    }
  })
}

resource "azapi_resource" "ruleset" {
  type      = "Microsoft.Network/dnsForwardingRulesets@2022-07-01"
  name      = "henglu921rgset"
  parent_id = azurerm_resource_group.test.id
  location  = azurerm_resource_group.test.location

  body = jsonencode({
    properties = {
      dnsResolverOutboundEndpoints = [{
        id = azapi_resource.outboundendpoint.id
      }]
    }
  })
}

resource "azapi_resource" "resolvervnetlink" {
  type      = "Microsoft.Network/dnsForwardingRulesets/virtualNetworkLinks@2022-07-01"
  name      = "testresolvervnetlink"
  parent_id = azapi_resource.ruleset.id

  body = jsonencode({
    properties = {
      virtualNetwork = {
        id = azurerm_virtual_network.vnet.id
      }
    }
  })
}


resource "azapi_resource" "forwardingrule" {
  type      = "Microsoft.Network/dnsForwardingRulesets/forwardingRules@2022-07-01"
  name      = "testforwardingrule"
  parent_id = azapi_resource.ruleset.id

  body = jsonencode({
    properties = {
      domainName          = "onprem.local."
      forwardingRuleState = "Enabled"
      targetDnsServers = [{
        ipAddress = "10.10.0.1"
        port      = 53
      }]
    }
  })
}