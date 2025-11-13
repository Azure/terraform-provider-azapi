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

data "azapi_client_config" "current" {}

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
  description = "The administrator password for the virtual machine scale set"
}

resource "azapi_resource" "resourceGroup" {
  type     = "Microsoft.Resources/resourceGroups@2020-06-01"
  name     = var.resource_name
  location = var.location
}

resource "azapi_resource" "assessmentMetadata" {
  type      = "Microsoft.Security/assessmentMetadata@2021-06-01"
  parent_id = "/subscriptions/${data.azapi_client_config.current.subscription_id}"
  name      = "fdaaa62c-1d42-45ab-be2f-2af194dd1700"
  body = {
    properties = {
      assessmentType = "CustomerManaged"
      description    = "Test Description"
      displayName    = "Test Display Name"
      severity       = "Medium"
    }
  }
}

resource "azapi_resource" "pricing" {
  type      = "Microsoft.Security/pricings@2023-01-01"
  parent_id = "/subscriptions/${data.azapi_client_config.current.subscription_id}"
  name      = "VirtualMachines"
  body = {
    properties = {
      extensions  = []
      pricingTier = "Standard"
      subPlan     = "P2"
    }
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

resource "azapi_resource" "virtualMachineScaleSet" {
  type      = "Microsoft.Compute/virtualMachineScaleSets@2024-11-01"
  parent_id = azapi_resource.resourceGroup.id
  name      = "${var.resource_name}-vmss"
  location  = var.location
  body = {
    properties = {
      additionalCapabilities                 = {}
      doNotRunExtensionsOnOverprovisionedVMs = false
      orchestrationMode                      = "Uniform"
      overprovision                          = true
      singlePlacementGroup                   = true
      upgradePolicy = {
        mode = "Manual"
      }
      virtualMachineProfile = {
        diagnosticsProfile = {
          bootDiagnostics = {
            enabled    = false
            storageUri = ""
          }
        }
        extensionProfile = {
          extensionsTimeBudget = "PT1H30M"
        }
        networkProfile = {
          networkInterfaceConfigurations = [{
            name = "example"
            properties = {
              dnsSettings = {
                dnsServers = []
              }
              enableAcceleratedNetworking = false
              enableIPForwarding          = false
              ipConfigurations = [{
                name = "internal"
                properties = {
                  applicationGatewayBackendAddressPools = []
                  applicationSecurityGroups             = []
                  loadBalancerBackendAddressPools       = []
                  loadBalancerInboundNatPools           = []
                  primary                               = true
                  privateIPAddressVersion               = "IPv4"
                  subnet = {
                    id = azapi_resource.subnet.id
                  }
                }
              }]
              primary = true
            }
          }]
        }
        osProfile = {
          adminPassword            = var.admin_password
          adminUsername            = "adminuser"
          allowExtensionOperations = true
          computerNamePrefix       = "${var.resource_name}-vmss"
          linuxConfiguration = {
            disablePasswordAuthentication = false
            provisionVMAgent              = true
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
              storageAccountType = "Standard_LRS"
            }
            osType                  = "Linux"
            writeAcceleratorEnabled = false
          }
        }
      }
    }
    sku = {
      capacity = 1
      name     = "Standard_B1s"
    }
  }
}

resource "azapi_resource" "assessment" {
  type      = "Microsoft.Security/assessments@2020-01-01"
  parent_id = azapi_resource.virtualMachineScaleSet.id
  name      = "fdaaa62c-1d42-45ab-be2f-2af194dd1700"
  body = {
    properties = {
      additionalData = {}
      resourceDetails = {
        source = "Azure"
      }
      status = {
        cause       = ""
        code        = "Healthy"
        description = ""
      }
    }
  }
}

