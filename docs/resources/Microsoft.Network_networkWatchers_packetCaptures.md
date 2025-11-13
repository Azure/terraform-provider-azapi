---
subcategory: "Microsoft.Network - Various Networking Services"
page_title: "networkWatchers/packetCaptures"
description: |-
  Manages a Configures Packet Capturing against a Virtual Machine using a Network Watcher.
---

# Microsoft.Network/networkWatchers/packetCaptures - Configures Packet Capturing against a Virtual Machine using a Network Watcher

This article demonstrates how to use `azapi` provider to manage the Configures Packet Capturing against a Virtual Machine using a Network Watcher resource in Azure.



## Example Usage

### default

```hcl
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

```



## Arguments Reference

The following arguments are supported:

* `type` - (Required) The type of the resource. This should be set to `Microsoft.Network/networkWatchers/packetCaptures@api-version`. The available api-versions for this resource are: [`2016-09-01`, `2016-12-01`, `2017-03-01`, `2017-03-30`, `2017-06-01`, `2017-08-01`, `2017-09-01`, `2017-10-01`, `2017-11-01`, `2018-01-01`, `2018-02-01`, `2018-04-01`, `2018-06-01`, `2018-07-01`, `2018-08-01`, `2018-10-01`, `2018-11-01`, `2018-12-01`, `2019-02-01`, `2019-04-01`, `2019-06-01`, `2019-07-01`, `2019-08-01`, `2019-09-01`, `2019-11-01`, `2019-12-01`, `2020-03-01`, `2020-04-01`, `2020-05-01`, `2020-06-01`, `2020-07-01`, `2020-08-01`, `2020-11-01`, `2021-02-01`, `2021-03-01`, `2021-05-01`, `2021-08-01`, `2022-01-01`, `2022-05-01`, `2022-07-01`, `2022-09-01`, `2022-11-01`, `2023-02-01`, `2023-04-01`, `2023-05-01`, `2023-06-01`, `2023-09-01`, `2023-11-01`, `2024-01-01`, `2024-03-01`, `2024-05-01`, `2024-07-01`, `2024-10-01`].

* `parent_id` - (Required) The ID of the azure resource in which this resource is created. The allowed values are:  
  `/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Network/networkWatchers/{resourceName}`

* `name` - (Required) Specifies the name of the azure resource. Changing this forces a new resource to be created.

* `body` - (Required) Specifies the configuration of the resource. More information about the arguments in `body` can be found in the [Microsoft documentation](https://learn.microsoft.com/en-us/azure/templates/Microsoft.Network/networkWatchers/packetCaptures?pivots=deployment-language-terraform).

For other arguments, please refer to the [azapi_resource](https://registry.terraform.io/providers/Azure/azapi/latest/docs/resources/resource) documentation.

## Import

 ```shell
 # Azure resource can be imported using the resource id, e.g.
 terraform import azapi_resource.example /subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Network/networkWatchers/{resourceName}/packetCaptures/{resourceName}
 
 # It also supports specifying API version by using the resource id with api-version as a query parameter, e.g.
 terraform import azapi_resource.example /subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Network/networkWatchers/{resourceName}/packetCaptures/{resourceName}?api-version=2024-10-01
 ```
