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
  type = string
  default = "acctest0001"
}

//variable "location" {
//  type    = string
//  default = "westeurope"
//}

data "azapi_resource" "resourceGroup" {
  type = "Microsoft.Resources/resourceGroups@2020-06-01"
  // name = var.resource_name
  name = "terraform-testing-qumulo"
  //  location = var.location
}

resource "azapi_resource" "vnet" {
  type      = "Microsoft.Network/virtualNetworks@2024-05-01"
  name      = var.resource_name
  parent_id = data.azapi_resource.resourceGroup.id
  location  = data.azapi_resource.resourceGroup.location

  body = {
    properties = {
      addressSpace = {
        addressPrefixes = ["10.0.0.0/16"]
      }
      privateEndpointVNetPolicies = "Disabled"
    }
  }

  schema_validation_enabled = false
  response_export_values    = ["*"]
}

resource "azapi_resource" "subnet" {
  type      = "Microsoft.Network/virtualNetworks/subnets@2024-05-01"
  name      = var.resource_name
  parent_id = azapi_resource.vnet.id
  location  = data.azapi_resource.resourceGroup.location

  body = {
    properties = {
      addressPrefix         = "10.0.1.0/24"
      defaultOutboundAccess = true
      delegations = [{
        name = "delegation"
        properties = {
          actions     = ["Microsoft.Network/virtualNetworks/subnets/join/action"]
          serviceName = "Qumulo.Storage/fileSystems"
        }
      }]
      privateEndpointNetworkPolicies    = "Disabled"
      privateLinkServiceNetworkPolicies = "Enabled"
    }
  }

  schema_validation_enabled = false
  response_export_values    = ["*"]
}

resource "azapi_resource" "qumuloFileSystem" {
  type      = "Qumulo.Storage/fileSystems@2024-06-19"
  name      = var.resource_name
  location  = data.azapi_resource.resourceGroup.location
  parent_id = data.azapi_resource.resourceGroup.id

  body = {
    properties = {
      adminPassword     = ")^X#ZX#JRyIY}t9"
      availabilityZone  = "1"
      delegatedSubnetId = azapi_resource.subnet.id
      marketplaceDetails = {
        offerId     = "qumulo-saas-mpp"
        planId      = "azure-native-qumulo-v3"
        publisherId = "qumulo1584033880660"
      }
      storageSku = "Cold_LRS"
      userDetails = {
        email = "test@test.com"
      }
    }
  }

  tags = {
    provisionResource = "55e031c5d497a3508838f44ec8767c84" // this is for test purpose to not to create real qumulo resource
    environment       = "terraform-acctests"
    some_key          = "some-value"
  }

  schema_validation_enabled = false
  response_export_values    = ["*"]
}
