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
  type     = "Microsoft.Resources/resourceGroups@2024-11-01"
  name     = var.resource_name
  location = var.location
}

resource "azapi_resource" "publicIPPrefix" {
  type      = "Microsoft.Network/publicIPPrefixes@2025-01-01"
  parent_id = azapi_resource.resourceGroup.id
  name      = var.resource_name
  location  = var.location
  body = {
    properties = {
      prefixLength           = 30
      publicIPAddressVersion = "IPv4"
    }
    sku = {
      name = "StandardV2"
    }
  }
  schema_validation_enabled = false
  response_export_values    = ["*"]
}

resource "azapi_resource" "publicIPPrefixV6" {
  type      = "Microsoft.Network/publicIPPrefixes@2025-01-01"
  parent_id = azapi_resource.resourceGroup.id
  name      = "${var.resource_name}V6"
  location  = var.location
  body = {
    properties = {
      prefixLength           = 125
      publicIPAddressVersion = "IPv6"
    }
    sku = {
      name = "StandardV2"
    }
  }
  schema_validation_enabled = false
  response_export_values    = ["*"]
}

resource "azapi_resource" "publicIPAddress" {
  type      = "Microsoft.Network/publicIPAddresses@2025-01-01"
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
      name = "StandardV2"
      tier = "Regional"
    }
  }
  schema_validation_enabled = false
  response_export_values    = ["*"]
}

resource "azapi_resource" "publicIPAddressV6" {
  type      = "Microsoft.Network/publicIPAddresses@2025-01-01"
  parent_id = azapi_resource.resourceGroup.id
  name      = "${var.resource_name}V6"
  location  = var.location
  body = {
    properties = {
      ddosSettings = {
        protectionMode = "VirtualNetworkInherited"
      }
      idleTimeoutInMinutes     = 4
      publicIPAddressVersion   = "IPv6"
      publicIPAllocationMethod = "Static"
    }
    sku = {
      name = "StandardV2"
      tier = "Regional"
    }
  }
  schema_validation_enabled = false
  response_export_values    = ["*"]
}

resource "azapi_resource" "virtualNetwork" {
  type      = "Microsoft.Network/virtualNetworks@2025-01-01"
  parent_id = azapi_resource.resourceGroup.id
  name      = var.resource_name
  location  = var.location
  body = {
    properties = {
      addressSpace = {
        addressPrefixes = [
          "10.0.0.0/16",
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

resource "azapi_resource" "natGateway" {
  type      = "Microsoft.Network/natGateways@2025-01-01"
  parent_id = azapi_resource.resourceGroup.id
  name      = var.resource_name
  location  = var.location
  body = {
    properties = {
      publicIpAddresses = [
        {
          id = azapi_resource.publicIPAddress.id,
        },
      ]
      publicIpAddressesV6 = [
        {
          id = azapi_resource.publicIPAddressV6.id,
        },
      ]
      publicIpPrefixes = [
        {
          id = azapi_resource.publicIPPrefix.id,
        },
      ]
      publicIpPrefixesV6 = [
        {
          id = azapi_resource.publicIPPrefixV6.id,
        },
      ]
    }
    sku = {
      name = "StandardV2"
    }
  }
  schema_validation_enabled = false
  ignore_casing             = false
  ignore_missing_property   = false
}
