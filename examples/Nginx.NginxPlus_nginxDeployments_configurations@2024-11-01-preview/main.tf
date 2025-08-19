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

resource "azapi_resource" "publicIPAddress" {
  type      = "Microsoft.Network/publicIPAddresses@2024-05-01"
  parent_id = azapi_resource.resourceGroup.id
  name      = "${var.resource_name}-pip"
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
    }
  }
  tags = {
    environment = "Production"
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
      addressPrefix         = "10.0.2.0/24"
      defaultOutboundAccess = true
      delegations = [{
        name = "delegation"
        properties = {
          serviceName = "NGINX.NGINXPLUS/nginxDeployments"
        }
      }]
      privateEndpointNetworkPolicies    = "Disabled"
      privateLinkServiceNetworkPolicies = "Enabled"
      serviceEndpointPolicies           = []
      serviceEndpoints                  = []
    }
  }
}

resource "azapi_resource" "nginxDeployment" {
  type      = "Nginx.NginxPlus/nginxDeployments@2024-11-01-preview"
  parent_id = azapi_resource.resourceGroup.id
  name      = "${var.resource_name}-nginx"
  location  = var.location
  body = {
    properties = {
      autoUpgradeProfile = {
        upgradeChannel = "stable"
      }
      enableDiagnosticsSupport = false
      networkProfile = {
        frontEndIPConfiguration = {
          publicIPAddresses = [{
            id = azapi_resource.publicIPAddress.id
          }]
        }
        networkInterfaceConfiguration = {
          subnetId = azapi_resource.subnet.id
        }
      }
      scalingProperties = {
        capacity = 10
      }
    }
    sku = {
      name = "standardv2_Monthly"
    }
  }
  tags = {
    foo = "bar"
  }
}

resource "azapi_resource" "configuration" {
  type      = "Nginx.NginxPlus/nginxDeployments/configurations@2024-11-01-preview"
  parent_id = azapi_resource.nginxDeployment.id
  name      = "default"
  body = {
    properties = {
      files = [{
        content     = "aHR0cCB7CiAgICBzZXJ2ZXIgewogICAgICAgIGxpc3RlbiA4MDsKICAgICAgICBsb2NhdGlvbiAvIHsKICAgICAgICAgICAgYXV0aF9iYXNpYyAiUHJvdGVjdGVkIEFyZWEiOwogICAgICAgICAgICBhdXRoX2Jhc2ljX3VzZXJfZmlsZSAvb3B0Ly5odHBhc3N3ZDsKICAgICAgICAgICAgZGVmYXVsdF90eXBlIHRleHQvaHRtbDsKICAgICAgICAgICAgcmV0dXJuIDIwMCAnPCFkb2N0eXBlIGh0bWw+PGh0bWwgbGFuZz0iZW4iPjxoZWFkPjwvaGVhZD48Ym9keT4KICAgICAgICAgICAgICAgIDxkaXY+dGhpcyBvbmUgd2lsbCBiZSB1cGRhdGVkPC9kaXY+CiAgICAgICAgICAgICAgICA8ZGl2PmF0IDEwOjM4IGFtPC9kaXY+CiAgICAgICAgICAgIDwvYm9keT48L2h0bWw+JzsKICAgICAgICB9CiAgICAgICAgaW5jbHVkZSBzaXRlLyouY29uZjsKICAgIH0KfQo="
        virtualPath = "/etc/nginx/nginx.conf"
      }]
      protectedFiles = [{
        content     = "dXNlcjokYXByMSRWZVVBNWt0LiRJampSay8vOG1pUnhEc1p2RDRkYUYxCg=="
        virtualPath = "/opt/.htpasswd"
      }]
      rootFile = "/etc/nginx/nginx.conf"
    }
  }
}

