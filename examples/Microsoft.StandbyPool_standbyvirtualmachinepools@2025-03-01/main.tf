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

provider "azurerm" {
  features {}
}

variable "resource_name" {
  type    = string
  default = "acctest0001"
}

variable "location" {
  type    = string
  default = "eastus"
}


resource "azapi_resource" "resourceGroup" {
  type     = "Microsoft.Resources/resourceGroups@2020-06-01"
  name     = "${var.resource_name}-rg"
  location = var.location
}

resource "azapi_resource" "virtualNetwork" {
  type      = "Microsoft.Network/virtualNetworks@2022-07-01"
  parent_id = azapi_resource.resourceGroup.id
  name      = "${var.resource_name}-vnet"
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
  name      = "${var.resource_name}-subnet"
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

resource "azapi_resource" "virtualMachineScaleSet" {
  type      = "Microsoft.Compute/virtualMachineScaleSets@2023-03-01"
  parent_id = azapi_resource.resourceGroup.id
  name      = "${var.resource_name}-vmss"
  location  = var.location
  body = {
    properties = {
      additionalCapabilities = {
      }
      orchestrationMode        = "Flexible"
      platformFaultDomainCount = 1
      singlePlacementGroup     = false
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
          networkApiVersion = "2022-07-01"
          networkInterfaceConfigurations = [
            {
              name = "example"
              properties = {
                dnsSettings = {
                  dnsServers = [
                  ]
                }
                enableAcceleratedNetworking = false
                enableIPForwarding          = false
                ipConfigurations = [
                  {
                    name = "internal"
                    properties = {
                      applicationGatewayBackendAddressPools = [
                      ]
                      applicationSecurityGroups = [
                      ]
                      loadBalancerBackendAddressPools = [
                      ]
                      primary                 = true
                      privateIPAddressVersion = "IPv4"
                      subnet = {
                        id = azapi_resource.subnet.id
                      }
                    }
                  },
                ]
                primary = true
              }
            },
          ]
        }
        osProfile = {
          adminUsername      = "adminuser"
          computerNamePrefix = var.resource_name
          linuxConfiguration = {
            disablePasswordAuthentication = true
            provisionVMAgent              = true
            ssh = {
              publicKeys = [
                {
                  keyData = "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQC+wWK73dCr+jgQOAxNsHAnNNNMEMWOHYEccp6wJm2gotpr9katuF/ZAdou5AaW1C61slRkHRkpRRX9FA9CYBiitZgvCCz+3nWNN7l/Up54Zps/pHWGZLHNJZRYyAB6j5yVLMVHIHriY49d/GZTZVNB8GoJv9Gakwc/fuEZYYl4YDFiGMBP///TzlI4jhiJzjKnEvqPFki5p2ZRJqcbCiF4pJrxUQR/RXqVFQdbRLZgYfJ8xGB878RENq3yQ39d8dVOkq4edbkzwcUmwwwkYVPIoDGsYLaRHnG+To7FvMeyO7xDVQkMKzopTQV8AuKpyvpqu0a9pWOMaiCyDytO7GGN you@me.com"
                  path    = "/home/adminuser/.ssh/authorized_keys"
                },
              ]
            }
          }
          secrets = [
          ]
        }
        priority = "Regular"
        storageProfile = {
          dataDisks = [
          ]
          imageReference = {
            offer     = "UbuntuServer"
            publisher = "Canonical"
            sku       = "16.04-LTS"
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
      name     = "Standard_F2"
      tier     = "Standard"
    }
  }
  schema_validation_enabled = false
  response_export_values    = ["*"]
}

data "azurerm_subscription" "primary" {}

data "azurerm_role_definition" "vm-contributor" {
  name = "Virtual Machine Contributor"
}

data "azurerm_role_definition" "nw-contributor" {
  name = "Network Contributor"
}

data "azurerm_role_definition" "mi-contributor" {
  name = "Managed Identity Contributor"
}

data "azuread_service_principal" "test" {
  display_name = "Standby Pool Resource Provider"
}

resource "azurerm_role_assignment" "vm-contributor" {
  scope              = azapi_resource.resourceGroup.id
  role_definition_id = "${data.azurerm_subscription.primary.id}${data.azurerm_role_definition.vm-contributor.id}"
  principal_id       = data.azuread_service_principal.test.object_id
}

resource "azurerm_role_assignment" "nw-contributor" {
  scope              = azapi_resource.resourceGroup.id
  role_definition_id = "${data.azurerm_subscription.primary.id}${data.azurerm_role_definition.nw-contributor.id}"
  principal_id       = data.azuread_service_principal.test.object_id
}

resource "azurerm_role_assignment" "mi-contributor" {
  scope              = azapi_resource.resourceGroup.id
  role_definition_id = "${data.azurerm_subscription.primary.id}${data.azurerm_role_definition.mi-contributor.id}"
  principal_id       = data.azuread_service_principal.test.object_id
}

resource "azapi_resource" "standbyVirtualMachinePool" {
  type      = "Microsoft.StandbyPool/standbyVirtualMachinePools@2025-03-01"
  parent_id = azapi_resource.resourceGroup.id
  name      = "${var.resource_name}-pool"
  location  = var.location

  body = {
    properties = {
      attachedVirtualMachineScaleSetId = azapi_resource.virtualMachineScaleSet.id
      elasticityProfile = {
        maxReadyCapacity = 2
        minReadyCapacity = 1
      }
      virtualMachineState = "Running"
    }
  }
  schema_validation_enabled = false
  ignore_casing             = false
  ignore_missing_property   = false

  depends_on = [
    azurerm_role_assignment.vm-contributor,
    azurerm_role_assignment.nw-contributor,
    azurerm_role_assignment.mi-contributor,
  ]
}
