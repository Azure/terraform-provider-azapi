---
subcategory: "Microsoft.Compute - Virtual Machines, Virtual Machine Scale Sets"
page_title: "virtualMachineScaleSets/extensions"
description: |-
  Manages a Extension for a Virtual Machine Scale Set.
---

# Microsoft.Compute/virtualMachineScaleSets/extensions - Extension for a Virtual Machine Scale Set

This article demonstrates how to use `azapi` provider to manage the Extension for a Virtual Machine Scale Set resource in Azure.

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
  name      = "internal"
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
  name      = var.resource_name
  location  = var.location
  body = {
    properties = {
      additionalCapabilities = {
      }
      doNotRunExtensionsOnOverprovisionedVMs = false
      orchestrationMode                      = "Uniform"
      overprovision                          = true
      scaleInPolicy = {
        forceDeletion = false
        rules = [
          "Default",
        ]
      }
      singlePlacementGroup = true
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
                      loadBalancerInboundNatPools = [
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

resource "azapi_resource" "extension" {
  type      = "Microsoft.Compute/virtualMachineScaleSets/extensions@2023-03-01"
  parent_id = azapi_resource.virtualMachineScaleSet.id
  name      = var.resource_name
  body = {
    properties = {
      autoUpgradeMinorVersion = true
      enableAutomaticUpgrade  = false
      provisionAfterExtensions = [
      ]
      publisher = "Microsoft.Azure.Extensions"
      settings = {
        commandToExecute = "echo $HOSTNAME"
      }
      suppressFailures   = false
      type               = "CustomScript"
      typeHandlerVersion = "2.0"
    }
  }
  schema_validation_enabled = false
  response_export_values    = ["*"]
}


```



## Arguments Reference

The following arguments are supported:

* `type` - (Required) The type of the resource. This should be set to `Microsoft.Compute/virtualMachineScaleSets/extensions@api-version`. The available api-versions for this resource are: [`2017-03-30`, `2017-12-01`, `2018-04-01`, `2018-06-01`, `2018-10-01`, `2019-03-01`, `2019-07-01`, `2019-12-01`, `2020-06-01`, `2020-12-01`, `2021-03-01`, `2021-04-01`, `2021-07-01`, `2021-11-01`, `2022-03-01`, `2022-08-01`, `2022-11-01`, `2023-03-01`, `2023-07-01`, `2023-09-01`, `2024-03-01`, `2024-07-01`, `2024-11-01`, `2025-04-01`].

* `parent_id` - (Required) The ID of the azure resource in which this resource is created. The allowed values are:  
  `/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Compute/virtualMachineScaleSets/{resourceName}`

* `name` - (Required) Specifies the name of the azure resource. Changing this forces a new resource to be created.

* `body` - (Required) Specifies the configuration of the resource. More information about the arguments in `body` can be found in the [Microsoft documentation](https://learn.microsoft.com/en-us/azure/templates/Microsoft.Compute/virtualMachineScaleSets/extensions?pivots=deployment-language-terraform).

For other arguments, please refer to the [azapi_resource](https://registry.terraform.io/providers/Azure/azapi/latest/docs/resources/resource) documentation.

## Import

 ```shell
 # Azure resource can be imported using the resource id, e.g.
 terraform import azapi_resource.example /subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Compute/virtualMachineScaleSets/{resourceName}/extensions/{resourceName}
 
 # It also supports specifying API version by using the resource id with api-version as a query parameter, e.g.
 terraform import azapi_resource.example /subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Compute/virtualMachineScaleSets/{resourceName}/extensions/{resourceName}?api-version=2025-04-01
 ```
