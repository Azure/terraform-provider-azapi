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

resource "azapi_resource" "vault" {
  type      = "Microsoft.RecoveryServices/vaults@2024-01-01"
  parent_id = azapi_resource.resourceGroup.id
  name      = "${var.resource_name}-rsv"
  location  = var.location
  body = {
    properties = {
      publicNetworkAccess = "Enabled"
      redundancySettings = {
        crossRegionRestore            = "Disabled"
        standardTierStorageRedundancy = "GeoRedundant"
      }
    }
    sku = {
      name = "Standard"
    }
  }
}

resource "azapi_resource" "replicationFabric" {
  type      = "Microsoft.RecoveryServices/vaults/replicationFabrics@2024-04-01"
  parent_id = azapi_resource.vault.id
  name      = "${var.resource_name}-fabric1"
  body = {
    properties = {
      customDetails = {
        instanceType = "Azure"
        location     = var.location
      }
    }
  }
}

resource "azapi_resource" "replicationFabric_1" {
  type      = "Microsoft.RecoveryServices/vaults/replicationFabrics@2024-04-01"
  parent_id = azapi_resource.vault.id
  name      = "${var.resource_name}-fabric2"
  body = {
    properties = {
      customDetails = {
        instanceType = "Azure"
        location     = "centralus"
      }
    }
  }
}

resource "azapi_resource" "virtualNetwork" {
  type      = "Microsoft.Network/virtualNetworks@2024-05-01"
  parent_id = azapi_resource.resourceGroup.id
  name      = "${var.resource_name}-vnet1"
  location  = var.location
  body = {
    properties = {
      addressSpace = {
        addressPrefixes = ["192.168.1.0/24"]
      }
      dhcpOptions = {
        dnsServers = []
      }
      privateEndpointVNetPolicies = "Disabled"
      subnets                     = []
    }
  }
}

resource "azapi_resource" "virtualNetwork_1" {
  type      = "Microsoft.Network/virtualNetworks@2024-05-01"
  parent_id = azapi_resource.resourceGroup.id
  name      = "${var.resource_name}-vnet2"
  location  = "centralus"
  body = {
    properties = {
      addressSpace = {
        addressPrefixes = ["192.168.2.0/24"]
      }
      dhcpOptions = {
        dnsServers = []
      }
      privateEndpointVNetPolicies = "Disabled"
      subnets                     = []
    }
  }
}

resource "azapi_resource" "replicationNetworkMapping" {
  type      = "Microsoft.RecoveryServices/vaults/replicationFabrics/replicationNetworks/replicationNetworkMappings@2024-04-01"
  parent_id = "${azapi_resource.replicationFabric.id}/replicationNetworks/${basename(azapi_resource.virtualNetwork.id)}"
  name      = "${var.resource_name}-mapping"
  body = {
    properties = {
      fabricSpecificDetails = {
        instanceType     = "AzureToAzure"
        primaryNetworkId = azapi_resource.virtualNetwork.id
      }
      recoveryFabricName = azapi_resource.replicationFabric_1.name
      recoveryNetworkId  = azapi_resource.virtualNetwork_1.id
    }
  }
}
