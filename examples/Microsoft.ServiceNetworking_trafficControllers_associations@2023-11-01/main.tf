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

resource "azapi_resource" "trafficController" {
  type      = "Microsoft.ServiceNetworking/trafficControllers@2023-11-01"
  parent_id = azapi_resource.resourceGroup.id
  name      = "${var.resource_name}-tc"
  location  = var.location
}

resource "azapi_resource" "virtualNetwork" {
  type      = "Microsoft.Network/virtualNetworks@2024-05-01"
  parent_id = azapi_resource.resourceGroup.id
  name      = "${var.resource_name}-vnet"
  location  = var.location
  body = {
    properties = {
      addressSpace = {
        addressPrefixes = ["10.0.0.0/16"]
      }
      dhcpOptions = {
        dnsServers = []
      }
      privateEndpointVNetPolicies = "Disabled"
      subnets                     = []
    }
  }
}

resource "azapi_resource" "subnet" {
  type      = "Microsoft.Network/virtualNetworks/subnets@2024-05-01"
  parent_id = azapi_resource.virtualNetwork.id
  name      = "${var.resource_name}-subnet"
  body = {
    properties = {
      addressPrefix         = "10.0.1.0/24"
      defaultOutboundAccess = true
      delegations = [{
        name = "delegation"
        properties = {
          serviceName = "Microsoft.ServiceNetworking/trafficControllers"
        }
      }]
      privateEndpointNetworkPolicies    = "Disabled"
      privateLinkServiceNetworkPolicies = "Enabled"
      serviceEndpointPolicies           = []
      serviceEndpoints                  = []
    }
  }
}

resource "azapi_resource" "association" {
  type      = "Microsoft.ServiceNetworking/trafficControllers/associations@2023-11-01"
  parent_id = azapi_resource.trafficController.id
  name      = "${var.resource_name}-assoc"
  location  = var.location
  body = {
    properties = {
      associationType = "subnets"
      subnet = {
        id = azapi_resource.subnet.id
      }
    }
  }
}

