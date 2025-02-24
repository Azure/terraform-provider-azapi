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
    }
  }
  schema_validation_enabled = false
  lifecycle {
    ignore_changes = [body.properties.subnets]
  }
}

resource "azapi_resource" "networkSecurityGroup" {
  type      = "Microsoft.Network/networkSecurityGroups@2024-05-01"
  parent_id = azapi_resource.resourceGroup.id
  name      = var.resource_name
  location  = azapi_resource.resourceGroup.location
  body = {
    properties = {
      securityRules = [{
        name = "MSSQLRule"
        properties = {
          access                     = "Allow"
          destinationAddressPrefix   = "*"
          destinationAddressPrefixes = []
          destinationPortRange       = "1433"
          destinationPortRanges      = []
          direction                  = "Inbound"
          priority                   = 1001
          protocol                   = "Tcp"
          sourceAddressPrefix        = "167.220.255.0/25"
          sourceAddressPrefixes      = []
          sourcePortRange            = "*"
          sourcePortRanges           = []
        }
      }]
    }
  }
}

resource "azapi_resource" "subnet" {
  type      = "Microsoft.Network/virtualNetworks/subnets@2022-07-01"
  parent_id = azapi_resource.virtualNetwork.id
  name      = var.resource_name
  body = {
    properties = {
      addressPrefix = "10.0.0.0/24"
      networkSecurityGroup = {
        id = azapi_resource.networkSecurityGroup.id
      }
    }
  }
}

resource "azapi_resource" "publicIPAddress" {
  type      = "Microsoft.Network/publicIPAddresses@2024-05-01"
  parent_id = azapi_resource.resourceGroup.id
  name      = var.resource_name
  location  = azapi_resource.resourceGroup.location
  body = {
    properties = {
      ddosSettings = {
        protectionMode = "VirtualNetworkInherited"
      }
      idleTimeoutInMinutes     = 4
      ipTags                   = []
      publicIPAddressVersion   = "IPv4"
      publicIPAllocationMethod = "Dynamic"
    }
    sku = {
      name = "Basic"
      tier = "Regional"
    }
  }
}

resource "azapi_resource" "networkInterface" {
  type      = "Microsoft.Network/networkInterfaces@2024-05-01"
  parent_id = azapi_resource.resourceGroup.id
  name      = var.resource_name
  location  = azapi_resource.resourceGroup.location
  body = {
    properties = {
      disableTcpStateTracking = false
      dnsSettings = {
        dnsServers = []
      }
      enableAcceleratedNetworking = false
      enableIPForwarding          = false
      ipConfigurations = [
        {
          type = "Microsoft.Network/networkInterfaces/ipConfigurations"
          name = "testconfiguration1"
          properties = {
            privateIPAddressVersion   = "IPv4"
            privateIPAllocationMethod = "Dynamic"
            publicIPAddress = {
              id = azapi_resource.publicIPAddress.id
            }
            subnet = {
              id = azapi_resource.subnet.id
            }
            primary          = true
            privateIPAddress = "10.0.0.4"
          }
        }
      ]
      nicType       = "Standard"
      auxiliaryMode = "None"
      auxiliarySku  = "None"
    }
  }
}

resource "azapi_resource" "virtualMachine" {
  type      = "Microsoft.Compute/virtualMachines@2024-07-01"
  parent_id = azapi_resource.resourceGroup.id
  name      = var.resource_name
  location  = azapi_resource.resourceGroup.location
  body = {
    properties = {
      osProfile = {
        adminUsername            = "testadmin"
        adminPassword            = "Password1234!"
        allowExtensionOperations = true
        computerName             = "winhost01"
        secrets                  = []
        windowsConfiguration = {
          timeZone               = "Pacific Standard Time"
          enableAutomaticUpdates = true
          patchSettings = {
            patchMode      = "AutomaticByOS"
            assessmentMode = "ImageDefault"
          }
          provisionVMAgent = true
        }
      }
      storageProfile = {
        dataDisks = []
        imageReference = {
          offer     = "SQL2017-WS2016"
          publisher = "MicrosoftSQLServer"
          sku       = "SQLDEV"
          version   = "latest"
        }
        osDisk = {
          diskSizeGB = 127
          managedDisk = {
            storageAccountType = "Premium_LRS"
          }
          name                    = "acctvm-250116171212663925OSDisk"
          osType                  = "Windows"
          writeAcceleratorEnabled = false
          caching                 = "ReadOnly"
          createOption            = "FromImage"
          deleteOption            = "Detach"
        }
      }
      hardwareProfile = {
        vmSize = "Standard_F2s"
      }
      networkProfile = {
        networkInterfaces = [
          {
            properties = {
              primary = false
            }
            id = azapi_resource.networkInterface.id
          }
        ]
      }
    }
  }
}


resource "azapi_resource" "sqlvirtualMachine" {
  type      = "Microsoft.SqlVirtualMachine/sqlVirtualMachines@2023-10-01"
  parent_id = azapi_resource.resourceGroup.id
  name      = azapi_resource.virtualMachine.name
  location  = azapi_resource.virtualMachine.location
  body = {
    properties = {
      sqlServerLicenseType     = "PAYG"
      virtualMachineResourceId = azapi_resource.virtualMachine.id
      enableAutomaticUpgrade   = true
      leastPrivilegeMode       = "Enabled"
      sqlImageOffer            = "SQL2017-WS2016"
      sqlImageSku              = "Developer"
      sqlManagement            = "Full"
    }
  }
}