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

variable "admin_username" {
  type        = string
  description = "The administrator username for the virtual machine"
}

variable "admin_password" {
  type        = string
  description = "The administrator password for the virtual machine"
  sensitive   = true
}

locals {
  os_disk_name            = "myosdisk1"
  data_disk_name          = "mydatadisk1"
  attached_data_disk_name = "myattacheddatadisk1"
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

resource "azapi_resource" "subnet" {
  type      = "Microsoft.Network/virtualNetworks/subnets@2022-07-01"
  parent_id = azapi_resource.virtualNetwork.id
  name      = var.resource_name
  body = {
    properties = {
      addressPrefix = "10.0.2.0/24"
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

resource "azapi_resource" "networkInterface" {
  type      = "Microsoft.Network/networkInterfaces@2022-07-01"
  parent_id = azapi_resource.resourceGroup.id
  name      = var.resource_name
  location  = var.location
  body = {
    properties = {
      enableAcceleratedNetworking = false
      enableIPForwarding          = false
      ipConfigurations = [
        {
          name = "testconfiguration1"
          properties = {
            primary                   = true
            privateIPAddressVersion   = "IPv4"
            privateIPAllocationMethod = "Dynamic"
            subnet = {
              id = azapi_resource.subnet.id
            }
          }
        },
      ]
    }
  }
  schema_validation_enabled = false
  response_export_values    = ["*"]
}

resource "azapi_resource" "attachedDisk" {
  type      = "Microsoft.Compute/disks@2022-03-02"
  parent_id = azapi_resource.resourceGroup.id
  name      = local.attached_data_disk_name
  location  = var.location
  body = {
    properties = {
      creationData = {
        createOption = "Empty"
      }
      diskSizeGB = 1
      encryption = {
        type = "EncryptionAtRestWithPlatformKey"
      }
      networkAccessPolicy = "AllowAll"
      osType              = "Linux"
      publicNetworkAccess = "Enabled"
    }
    sku = {
      name = "Standard_LRS"
    }
  }
  schema_validation_enabled = false
  response_export_values    = ["*"]
}

resource "azapi_resource" "virtualMachine" {
  type      = "Microsoft.Compute/virtualMachines@2023-03-01"
  parent_id = azapi_resource.resourceGroup.id
  name      = var.resource_name
  location  = var.location
  body = {
    properties = {
      hardwareProfile = {
        vmSize = "Standard_F2"
      }
      networkProfile = {
        networkInterfaces = [
          {
            id = azapi_resource.networkInterface.id
            properties = {
              primary = false
            }
          },
        ]
      }
      osProfile = {
        adminPassword = var.admin_password
        adminUsername = var.admin_username
        computerName  = "hostname230630032848831819"
        linuxConfiguration = {
          disablePasswordAuthentication = false
        }
      }
      storageProfile = {
        imageReference = {
          offer     = "UbuntuServer"
          publisher = "Canonical"
          sku       = "16.04-LTS"
          version   = "latest"
        }
        osDisk = {
          caching                 = "ReadWrite"
          createOption            = "FromImage"
          name                    = local.os_disk_name
          writeAcceleratorEnabled = false
        }
        dataDisks = [
          {
            caching      = "ReadWrite"
            createOption = "Empty"
            name         = local.data_disk_name
            diskSizeGB   = 1
            lun          = 1
            managedDisk = {
              storageAccountType = "Standard_LRS"
            }
          },
          {
            caching      = "ReadWrite"
            createOption = "Attach"
            name         = azapi_resource.attachedDisk.name
            lun          = 2
            managedDisk = {
              id = azapi_resource.attachedDisk.id
            }
          }
        ]
      }
    }
  }
  schema_validation_enabled = false
  response_export_values    = ["*"]
}
