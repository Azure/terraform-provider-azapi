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

variable "admin_password" {
  type        = string
  sensitive   = true
  description = "The administrator password for the virtual machine"
}

resource "azapi_resource" "resourceGroup" {
  type     = "Microsoft.Resources/resourceGroups@2020-06-01"
  name     = var.resource_name
  location = var.location
}

resource "azapi_resource" "networkWatcher" {
  type      = "Microsoft.Network/networkWatchers@2024-05-01"
  parent_id = azapi_resource.resourceGroup.id
  name      = "${var.resource_name}-nw"
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
    }
  }
}

resource "azapi_resource" "subnet" {
  type      = "Microsoft.Network/virtualNetworks/subnets@2024-05-01"
  parent_id = azapi_resource.virtualNetwork.id
  name      = "internal"
  body = {
    properties = {
      addressPrefix                     = "10.0.2.0/24"
      defaultOutboundAccess             = true
      delegations                       = []
      privateEndpointNetworkPolicies    = "Disabled"
      privateLinkServiceNetworkPolicies = "Enabled"
      serviceEndpointPolicies           = []
      serviceEndpoints                  = []
    }
  }
}

resource "azapi_resource" "networkInterface" {
  type      = "Microsoft.Network/networkInterfaces@2024-05-01"
  parent_id = azapi_resource.resourceGroup.id
  name      = "${var.resource_name}-nic"
  location  = var.location
  body = {
    properties = {
      enableAcceleratedNetworking = false
      enableIPForwarding          = false
      ipConfigurations = [{
        name = "ipconfig1"
        properties = {
          primary                   = true
          privateIPAddressVersion   = "IPv4"
          privateIPAllocationMethod = "Dynamic"
          subnet = {
            id = azapi_resource.subnet.id
          }
        }
      }]
    }
  }
}

resource "azapi_resource" "virtualMachine" {
  type      = "Microsoft.Compute/virtualMachines@2024-03-01"
  parent_id = azapi_resource.resourceGroup.id
  name      = "${var.resource_name}-vm"
  location  = var.location
  body = {
    properties = {
      hardwareProfile = {
        vmSize = "Standard_B1s"
      }
      networkProfile = {
        networkInterfaces = [{
          id = azapi_resource.networkInterface.id
          properties = {
            primary = true
          }
        }]
      }
      osProfile = {
        adminPassword = var.admin_password
        adminUsername = "testadmin"
        computerName  = "${var.resource_name}-vm"
        linuxConfiguration = {
          disablePasswordAuthentication = false
        }
      }
      storageProfile = {
        imageReference = {
          offer     = "0001-com-ubuntu-server-jammy"
          publisher = "Canonical"
          sku       = "22_04-lts"
          version   = "latest"
        }
        osDisk = {
          caching      = "ReadWrite"
          createOption = "FromImage"
          managedDisk = {
            storageAccountType = "Standard_LRS"
          }
          name                    = "${var.resource_name}-osdisk"
          writeAcceleratorEnabled = false
        }
      }
    }
  }
}

resource "azapi_resource" "extension" {
  type      = "Microsoft.Compute/virtualMachines/extensions@2024-03-01"
  parent_id = azapi_resource.virtualMachine.id
  name      = "network-watcher"
  location  = var.location
  body = {
    properties = {
      autoUpgradeMinorVersion = true
      enableAutomaticUpgrade  = false
      publisher               = "Microsoft.Azure.NetworkWatcher"
      suppressFailures        = false
      type                    = "NetworkWatcherAgentLinux"
      typeHandlerVersion      = "1.4"
    }
  }
}

resource "azapi_resource" "packetCapture" {
  type      = "Microsoft.Network/networkWatchers/packetCaptures@2024-05-01"
  parent_id = azapi_resource.networkWatcher.id
  name      = "${var.resource_name}-pc"
  body = {
    properties = {
      bytesToCapturePerPacket = 0
      storageLocation = {
        filePath = "/var/captures/packet.cap"
      }
      target               = azapi_resource.virtualMachine.id
      targetType           = "AzureVM"
      timeLimitInSeconds   = 18000
      totalBytesPerSession = 1073741824
    }
  }
}
