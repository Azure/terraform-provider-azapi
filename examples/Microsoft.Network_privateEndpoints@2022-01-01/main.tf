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

resource "azapi_resource" "resourceGroup" {
  type     = "Microsoft.Resources/resourceGroups@2020-06-01"
  name     = var.resource_name
  location = var.location
}

resource "azapi_resource" "virtualNetwork" {
  type      = "Microsoft.Network/virtualNetworks@2022-07-01"
  parent_id = azapi_resource.resourceGroup.id
  name      = var.resource_name
  location  = var.location
  body = {
    properties = {
      addressSpace = {
        addressPrefixes = [
          "10.5.0.0/16",
        ]
      }
      dhcpOptions = {
        dnsServers = [
        ]
      }
      subnets = [
      ]
    }
  }
  schema_validation_enabled = false
  response_export_values    = ["*"]
  lifecycle {
    ignore_changes = [body.properties.subnets]
  }
}

resource "azapi_resource" "publicIPAddress" {
  type      = "Microsoft.Network/publicIPAddresses@2022-07-01"
  parent_id = azapi_resource.resourceGroup.id
  name      = var.resource_name
  location  = var.location
  body = {
    properties = {
      ddosSettings = {
        protectionMode = "VirtualNetworkInherited"
      }
      idleTimeoutInMinutes     = 4
      publicIPAddressVersion   = "IPv4"
      publicIPAllocationMethod = "Static"
    }
    sku = {
      name = "Standard"
      tier = "Regional"
    }
  }
  schema_validation_enabled = false
  response_export_values    = ["*"]
}

resource "azapi_resource" "loadBalancer" {
  type      = "Microsoft.Network/loadBalancers@2022-07-01"
  parent_id = azapi_resource.resourceGroup.id
  name      = var.resource_name
  location  = var.location
  body = {
    properties = {
      frontendIPConfigurations = [
        {
          name = var.resource_name
          properties = {
            publicIPAddress = {
              id = azapi_resource.publicIPAddress.id
            }
          }
        },
      ]
    }
    sku = {
      name = "Standard"
      tier = "Regional"
    }
  }
  schema_validation_enabled = false
  response_export_values    = ["*"]
}

resource "azapi_resource" "subnet" {
  type      = "Microsoft.Network/virtualNetworks/subnets@2022-07-01"
  parent_id = azapi_resource.virtualNetwork.id
  name      = var.resource_name
  body = {
    properties = {
      addressPrefix = "10.5.4.0/24"
      delegations = [
      ]
      privateEndpointNetworkPolicies    = "Enabled"
      privateLinkServiceNetworkPolicies = "Disabled"
      serviceEndpointPolicies = [
      ]
      serviceEndpoints = [
      ]
    }
  }
  schema_validation_enabled = false
  response_export_values    = ["*"]
}

resource "azapi_resource" "privateLinkService" {
  type      = "Microsoft.Network/privateLinkServices@2022-07-01"
  parent_id = azapi_resource.resourceGroup.id
  name      = var.resource_name
  location  = var.location
  body = {
    properties = {
      autoApproval = {
        subscriptions = [
        ]
      }
      enableProxyProtocol = false
      fqdns = [
      ]
      ipConfigurations = [
        {
          name = "primaryIpConfiguration-230630033653892379"
          properties = {
            primary                   = true
            privateIPAddress          = ""
            privateIPAddressVersion   = "IPv4"
            privateIPAllocationMethod = "Dynamic"
            subnet = {
              id = azapi_resource.subnet.id
            }
          }
        },
      ]
      loadBalancerFrontendIpConfigurations = [
        {
          id = azapi_resource.loadBalancer.output.properties.frontendIPConfigurations[0].id
        },
      ]
      visibility = {
        subscriptions = [
        ]
      }
    }
  }
  schema_validation_enabled = false
  response_export_values    = ["*"]
}

resource "azapi_resource" "privateEndpoint" {
  type      = "Microsoft.Network/privateEndpoints@2022-07-01"
  parent_id = azapi_resource.resourceGroup.id
  name      = var.resource_name
  location  = var.location
  body = {
    properties = {
      privateLinkServiceConnections = [
        {
          name = azapi_resource.privateLinkService.name
          properties = {
            privateLinkServiceId = azapi_resource.privateLinkService.id
          }
        }
      ]
      subnet = {
        id = azapi_resource.subnet.id
      }
    }
  }
  schema_validation_enabled = false
  response_export_values    = ["*"]
}
