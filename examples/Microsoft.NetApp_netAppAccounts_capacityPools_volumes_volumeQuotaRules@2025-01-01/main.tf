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
  tags = {
    CreatedOnDate    = "2023-08-17T08:01:00Z"
    SkipASMAzSecPack = "true"
    SkipNRMSNSG      = "true"
  }
}

resource "azapi_resource" "networkSecurityGroup" {
  type      = "Microsoft.Network/networkSecurityGroups@2024-05-01"
  parent_id = azapi_resource.resourceGroup.id
  name      = "${var.resource_name}-nsg"
  location  = var.location
  body = {
    properties = {
      securityRules = []
    }
  }
  tags = {
    CreatedOnDate    = "2023-08-17T08:01:00Z"
    SkipASMAzSecPack = "true"
  }
}

resource "azapi_resource" "virtualNetwork" {
  type      = "Microsoft.Network/virtualNetworks@2024-05-01"
  parent_id = azapi_resource.resourceGroup.id
  name      = "${var.resource_name}-vnet"
  location  = var.location
  body = {
    properties = {
      addressSpace = {
        addressPrefixes = ["10.88.0.0/16"]
      }
      dhcpOptions = {
        dnsServers = []
      }
      subnets = []
    }
  }
  tags = {
    CreatedOnDate    = "2023-08-17T08:01:00Z"
    SkipASMAzSecPack = "true"
  }
}

resource "azapi_resource" "subnet" {
  type      = "Microsoft.Network/virtualNetworks/subnets@2024-05-01"
  parent_id = azapi_resource.virtualNetwork.id
  name      = "${var.resource_name}-subnet"
  body = {
    properties = {
      addressPrefix         = "10.88.2.0/24"
      defaultOutboundAccess = true
      delegations = [{
        name = "netapp-delegation"
        properties = {
          serviceName = "Microsoft.NetApp/volumes"
        }
      }]
      networkSecurityGroup = {
        id = azapi_resource.networkSecurityGroup.id
      }
      privateEndpointNetworkPolicies    = "Disabled"
      privateLinkServiceNetworkPolicies = "Enabled"
      serviceEndpointPolicies           = []
      serviceEndpoints                  = []
    }
  }
}

resource "azapi_resource" "netAppAccount" {
  type      = "Microsoft.NetApp/netAppAccounts@2025-01-01"
  parent_id = azapi_resource.resourceGroup.id
  name      = "${var.resource_name}-acct"
  location  = var.location
  body = {
    properties = {}
  }
  tags = {
    CreatedOnDate    = "2023-08-17T08:01:00Z"
    SkipASMAzSecPack = "true"
  }
}

resource "azapi_resource" "capacityPool" {
  type      = "Microsoft.NetApp/netAppAccounts/capacityPools@2025-01-01"
  parent_id = azapi_resource.netAppAccount.id
  name      = "${var.resource_name}-pool"
  location  = var.location
  body = {
    properties = {
      coolAccess     = false
      encryptionType = "Single"
      qosType        = "Auto"
      serviceLevel   = "Standard"
      size           = 4398046511104
    }
  }
  tags = {
    CreatedOnDate    = "2023-08-17T08:01:00Z"
    SkipASMAzSecPack = "true"
  }
}

resource "azapi_resource" "volume" {
  type      = "Microsoft.NetApp/netAppAccounts/capacityPools/volumes@2025-01-01"
  parent_id = azapi_resource.capacityPool.id
  name      = "${var.resource_name}-vol"
  location  = var.location
  body = {
    properties = {
      creationToken  = "${var.resource_name}-path"
      dataProtection = {}
      exportPolicy = {
        rules = []
      }
      protocolTypes  = ["NFSv3"]
      serviceLevel   = "Standard"
      subnetId       = azapi_resource.subnet.id
      usageThreshold = 107374182400
    }
  }
  tags = {
    CreatedOnDate    = "2022-07-08T23:50:21Z"
    SkipASMAzSecPack = "true"
  }
}

resource "azapi_resource" "volumeQuotaRule" {
  type      = "Microsoft.NetApp/netAppAccounts/capacityPools/volumes/volumeQuotaRules@2025-01-01"
  parent_id = azapi_resource.volume.id
  name      = "${var.resource_name}-quota"
  location  = var.location
  body = {
    properties = {
      quotaSizeInKiBs = 2048
      quotaType       = "DefaultGroupQuota"
    }
  }
}
