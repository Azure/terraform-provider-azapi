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
  default = "eastus"
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

resource "azapi_resource" "userAssignedIdentity" {
  type      = "Microsoft.ManagedIdentity/userAssignedIdentities@2023-01-31"
  parent_id = azapi_resource.resourceGroup.id
  name      = "${var.resource_name}-uai"
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
      subnets                     = []
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
  identity {
    type         = "SystemAssigned, UserAssigned"
    identity_ids = [azapi_resource.userAssignedIdentity.id]
  }
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
        vmSize = "Standard_B2s"
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
        adminPassword            = var.admin_password
        adminUsername            = "adminuser"
        allowExtensionOperations = true
        computerName             = "${var.resource_name}-vm"
        linuxConfiguration = {
          disablePasswordAuthentication = false
          patchSettings = {
            assessmentMode = "ImageDefault"
            patchMode      = "ImageDefault"
          }
          provisionVMAgent = true
          ssh = {
            publicKeys = []
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
            storageAccountType = "Premium_LRS"
          }
          osType                  = "Linux"
          writeAcceleratorEnabled = false
        }
      }
    }
  }
}

resource "azapi_resource" "runCommand" {
  type      = "Microsoft.Compute/virtualMachines/runCommands@2023-03-01"
  parent_id = azapi_resource.virtualMachine.id
  name      = "${var.resource_name}-runcommand"
  location  = var.location
  body = {
    properties = {
      asyncExecution      = false
      errorBlobUri        = ""
      outputBlobUri       = ""
      parameters          = []
      protectedParameters = []
      runAsPassword       = ""
      runAsUser           = ""
      source = {
        script = "echo 'hello world'"
      }
      timeoutInSeconds                = 1200
      treatFailureAsDeploymentFailure = true
    }
  }
}