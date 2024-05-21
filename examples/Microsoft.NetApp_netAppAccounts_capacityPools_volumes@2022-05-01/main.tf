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
  default = "centralus"
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
          "10.6.0.0/16",
        ]
      }
      dhcpOptions = {
        dnsServers = [
        ]
      }
      subnets = [
      ]
    }
    tags = {
      SkipASMAzSecPack = "true"
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
    tags = {
      SkipASMAzSecPack = "true"
    }
  }
  schema_validation_enabled = false
  response_export_values    = ["*"]
}

resource "azapi_resource" "subnet" {
  type      = "Microsoft.Network/virtualNetworks/subnets@2022-07-01"
  parent_id = azapi_resource.virtualNetwork.id
  name      = "GatewaySubnet"
  body = {
    properties = {
      addressPrefix = "10.6.1.0/24"
      delegations = [
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

resource "azapi_resource" "subnet2" {
  type      = "Microsoft.Network/virtualNetworks/subnets@2022-07-01"
  parent_id = azapi_resource.virtualNetwork.id
  name      = var.resource_name
  body = {
    properties = {
      addressPrefix = "10.6.2.0/24"
      delegations = [
        {
          name = "testdelegation"
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
      serviceLevel = "Standard"
      size         = 4.398046511104e+12
    }
    tags = {
      SkipASMAzSecPack = "true"
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
      avsDataStore  = "Enabled"
      creationToken = "my-unique-file-path-230630034120103726"
      dataProtection = {
      }
      exportPolicy = {
        rules = [
          {
            allowedClients = "0.0.0.0/0"
            cifs           = false
            hasRootAccess  = true
            nfsv3          = true
            nfsv41         = false
            ruleIndex      = 1
            unixReadOnly   = false
            unixReadWrite  = true
          },
        ]
      }
      networkFeatures = "Basic"
      protocolTypes = [
        "NFSv3",
      ]
      serviceLevel             = "Standard"
      snapshotDirectoryVisible = true
      subnetId                 = azapi_resource.subnet2.id
      usageThreshold           = 1.073741824e+11
      volumeType               = ""
    }
    tags = {
      SkipASMAzSecPack = "true"
    }
    zones = [
    ]
  }
  schema_validation_enabled = false
  response_export_values    = ["*"]
}

