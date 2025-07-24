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

variable "attached_resource_name" {
  type    = string
  default = "acctest0002"
}

locals {
  os_disk_name          = "myosdisk1"
  attached_os_disk_name = "myosdisk2"
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
      }
    }
  }
  schema_validation_enabled = false
  response_export_values    = ["*"]
}

data "azapi_resource" "managedDisk" {
  type      = "Microsoft.Compute/disks@2023-10-02"
  parent_id = azapi_resource.resourceGroup.id
  name      = local.os_disk_name

  depends_on = [azapi_resource.virtualMachine]
}

resource "azapi_resource" "snapshot" {
  type      = "Microsoft.Compute/snapshots@2023-10-02"
  parent_id = azapi_resource.resourceGroup.id
  name      = var.resource_name
  location  = var.location
  body = {
    sku = {
      name = "Standard_ZRS"
    }
    properties = {
      creationData = {
        createOption     = "Copy"
        sourceResourceId = data.azapi_resource.managedDisk.id
      }
      diskSizeGB = 30
      encryption = {
        type = "EncryptionAtRestWithPlatformKey"
      }
      networkAccessPolicy = "AllowAll"
      osType              = "Linux"
      hyperVGeneration    = "V1"
      incremental         = true
      publicNetworkAccess = "Enabled"
      supportedCapabilities = {
        architecture = "x64"
      }
    }
  }
  schema_validation_enabled = false
  response_export_values    = ["*"]
}

resource "azapi_resource" "attachedManagedDisk" {
  type      = "Microsoft.Compute/disks@2023-10-02"
  parent_id = azapi_resource.resourceGroup.id
  name      = local.attached_os_disk_name
  location  = var.location
  body = {
    properties = {
      creationData = {
        createOption     = "Copy",
        sourceResourceId = azapi_resource.snapshot.id
      }

      diskSizeGB = 30
      encryption = {
        type = "EncryptionAtRestWithPlatformKey"
      }
      networkAccessPolicy = "AllowAll"
      osType              = "Linux"
      hyperVGeneration    = "V1"
      publicNetworkAccess = "Enabled"
      supportedCapabilities = {
        architecture = "x64"
      }
    }
    sku = {
      name = "Standard_LRS"
    }
    zones = [
      "1"
    ]
  }

  schema_validation_enabled = false
  response_export_values    = ["*"]
}

resource "azapi_resource" "attachedNetworkInterface" {
  type      = "Microsoft.Network/networkInterfaces@2022-07-01"
  parent_id = azapi_resource.resourceGroup.id
  name      = var.attached_resource_name
  location  = var.location
  body = {
    properties = {
      enableAcceleratedNetworking = false
      enableIPForwarding          = false
      ipConfigurations = [
        {
          name = "testconfiguration2"
          properties = {
            primary                   = true
            privateIPAddressVersion   = "IPv4"
            privateIPAllocationMethod = "Dynamic"
            subnet = {
              id = azapi_resource.subnet.id
            }
          }
        }
      ]
    }
  }
  schema_validation_enabled = false
  response_export_values    = ["*"]
}

resource "azapi_resource" "attachedVirtualMachine" {
  type      = "Microsoft.Compute/virtualMachines@2023-03-01"
  parent_id = azapi_resource.resourceGroup.id
  name      = var.attached_resource_name
  location  = var.location
  body = {
    properties = {
      hardwareProfile = {
        vmSize = "Standard_F2"
      }
      networkProfile = {
        networkInterfaces = [
          {
            id = azapi_resource.attachedNetworkInterface.id
            properties = {
              primary = false
            }
          },
        ]
      }
      storageProfile = {
        osDisk = {
          caching                 = "ReadWrite"
          createOption            = "Attach"
          name                    = local.attached_os_disk_name
          osType                  = "Linux",
          writeAcceleratorEnabled = false
          managedDisk = {
            id = azapi_resource.attachedManagedDisk.id
          }
        }
      }
    }
  }
  schema_validation_enabled = false
  response_export_values    = ["*"]
}

