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

resource "azapi_resource" "netAppAccount" {
  type      = "Microsoft.NetApp/netAppAccounts@2022-05-01"
  parent_id = azapi_resource.resourceGroup.id
  name      = var.resource_name
  location  = var.location
  body = {
    properties = {
      activeDirectories = [
      ]
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
      addressPrefix = "10.0.2.0/24"
      delegations = [
        {
          name = "netapp"
          properties = {
            serviceName = "Microsoft.Netapp/volumes"
          }
        },
      ]
      privateEndpointNetworkPolicies    = "Enabled"
      privateLinkServiceNetworkPolicies = "Enabled"
      serviceEndpointPolicies = [
      ]
      serviceEndpoints = [
      ]
    }
  }
  schema_validation_enabled = false
  response_export_values    = ["*"]
}

resource "azapi_resource" "capacityPool" {
  type      = "Microsoft.NetApp/netAppAccounts/capacityPools@2022-05-01"
  parent_id = azapi_resource.netAppAccount.id
  name      = var.resource_name
  location  = var.location
  body = {
    properties = {
      serviceLevel = "Premium"
      size         = 4.398046511104e+12
    }
  }
  schema_validation_enabled = false
  response_export_values    = ["*"]
}

resource "azapi_resource" "volume" {
  type      = "Microsoft.NetApp/netAppAccounts/capacityPools/volumes@2022-05-01"
  parent_id = azapi_resource.capacityPool.id
  name      = var.resource_name
  location  = var.location
  body = {
    properties = {
      avsDataStore  = "Disabled"
      creationToken = "my-unique-file-path-230630033642692134"
      dataProtection = {
      }
      exportPolicy = {
        rules = [
        ]
      }
      networkFeatures = "Basic"
      protocolTypes = [
        "NFSv3",
      ]
      securityStyle            = "Unix"
      serviceLevel             = "Premium"
      snapshotDirectoryVisible = false
      snapshotId               = ""
      subnetId                 = azapi_resource.subnet.id
      usageThreshold           = 1.073741824e+11
      volumeType               = ""
    }
    zones = [
    ]
  }
  schema_validation_enabled = false
  response_export_values    = ["*"]
}

resource "azapi_resource" "snapshot" {
  type                      = "Microsoft.NetApp/netAppAccounts/capacityPools/volumes/snapshots@2022-05-01"
  parent_id                 = azapi_resource.volume.id
  name                      = var.resource_name
  location                  = var.location
  schema_validation_enabled = false
  response_export_values    = ["*"]
}

