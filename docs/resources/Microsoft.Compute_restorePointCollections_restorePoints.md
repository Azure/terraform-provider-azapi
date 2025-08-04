---
subcategory: "Microsoft.Compute - Virtual Machines, Virtual Machine Scale Sets"
page_title: "restorePointCollections/restorePoints"
description: |-
  Manages a Virtual Machine Restore Point.
---

# Microsoft.Compute/restorePointCollections/restorePoints - Virtual Machine Restore Point

This article demonstrates how to use `azapi` provider to manage the Virtual Machine Restore Point resource in Azure.

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

resource "azapi_resource" "resourceGroup" {
  type     = "Microsoft.Resources/resourceGroups@2020-06-01"
  name     = var.resource_name
  location = var.location
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
      addressPrefix                     = "10.0.0.0/24"
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
        name = "internal"
        properties = {
          primary                   = false
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
      additionalCapabilities = {}
      applicationProfile = {
        galleryApplications = []
      }
      diagnosticsProfile = {
        bootDiagnostics = {
          enabled    = false
          storageUri = ""
        }
      }
      extensionsTimeBudget = "PT1H30M"
      hardwareProfile = {
        vmSize = "Standard_F2"
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
        adminUsername            = "adminuser"
        allowExtensionOperations = true
        computerName             = "${var.resource_name}-vm"
        linuxConfiguration = {
          disablePasswordAuthentication = true
          patchSettings = {
            assessmentMode = "ImageDefault"
            patchMode      = "ImageDefault"
          }
          provisionVMAgent = true
          ssh = {
            publicKeys = [{
              keyData = "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQC+wWK73dCr+jgQOAxNsHAnNNNMEMWOHYEccp6wJm2gotpr9katuF/ZAdou5AaW1C61slRkHRkpRRX9FA9CYBiitZgvCCz+3nWNN7l/Up54Zps/pHWGZLHNJZRYyAB6j5yVLMVHIHriY49d/GZTZVNB8GoJv9Gakwc/fuEZYYl4YDFiGMBP///TzlI4jhiJzjKnEvqPFki5p2ZRJqcbCiF4pJrxUQR/RXqVFQdbRLZgYfJ8xGB878RENq3yQ39d8dVOkq4edbkzwcUmwwwkYVPIoDGsYLaRHnG+To7FvMeyO7xDVQkMKzopTQV8AuKpyvpqu0a9pWOMaiCyDytO7GGN you@me.com"
              path    = "/home/adminuser/.ssh/authorized_keys"
            }]
          }
        }
        secrets = []
      }
      priority = "Regular"
      storageProfile = {
        dataDisks = []
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
          osType                  = "Linux"
          writeAcceleratorEnabled = false
        }
      }
    }
  }
}

resource "azapi_resource" "restorePointCollection" {
  type      = "Microsoft.Compute/restorePointCollections@2024-03-01"
  parent_id = azapi_resource.resourceGroup.id
  name      = "${var.resource_name}-rpc"
  location  = var.location
  body = {
    properties = {
      source = {
        id = azapi_resource.virtualMachine.id
      }
    }
  }
  tags = {
    foo = "bar"
  }
}

resource "azapi_resource" "restorePoint" {
  type      = "Microsoft.Compute/restorePointCollections/restorePoints@2024-03-01"
  parent_id = azapi_resource.restorePointCollection.id
  name      = "${var.resource_name}-rp"
  body = {
    properties = {}
  }
}
```



## Arguments Reference

The following arguments are supported:

* `type` - (Required) The type of the resource. This should be set to `Microsoft.Compute/restorePointCollections/restorePoints@api-version`. The available api-versions for this resource are: [`2021-03-01`, `2021-04-01`, `2021-07-01`, `2021-11-01`, `2022-03-01`, `2022-08-01`, `2022-11-01`, `2023-03-01`, `2023-07-01`, `2023-09-01`, `2024-03-01`, `2024-07-01`, `2024-11-01`].

* `parent_id` - (Required) The ID of the azure resource in which this resource is created. The allowed values are:  
  `/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Compute/restorePointCollections/{resourceName}`

* `name` - (Required) Specifies the name of the azure resource. Changing this forces a new resource to be created.

* `body` - (Required) Specifies the configuration of the resource. More information about the arguments in `body` can be found in the [Microsoft documentation](https://learn.microsoft.com/en-us/azure/templates/Microsoft.Compute/restorePointCollections/restorePoints?pivots=deployment-language-terraform).

For other arguments, please refer to the [azapi_resource](https://registry.terraform.io/providers/Azure/azapi/latest/docs/resources/resource) documentation.

## Import

 ```shell
 # Azure resource can be imported using the resource id, e.g.
 terraform import azapi_resource.example /subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Compute/restorePointCollections/{resourceName}/restorePoints/{resourceName}
 
 # It also supports specifying API version by using the resource id with api-version as a query parameter, e.g.
 terraform import azapi_resource.example /subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Compute/restorePointCollections/{resourceName}/restorePoints/{resourceName}?api-version=2024-11-01
 ```
