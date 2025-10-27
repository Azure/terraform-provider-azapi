---
subcategory: "Microsoft.Security - Security Center"
page_title: "assessments"
description: |-
  Manages a Security Center Assessment for Azure Security Center.
---

# Microsoft.Security/assessments - Security Center Assessment for Azure Security Center

This article demonstrates how to use `azapi` provider to manage the Security Center Assessment for Azure Security Center resource in Azure.



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

data "azapi_client_config" "current" {}

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
          adminPassword            = "P@ssword1234!"
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


```



## Arguments Reference

The following arguments are supported:

* `type` - (Required) The type of the resource. This should be set to `Microsoft.Security/assessments@api-version`. The available api-versions for this resource are: [`2019-01-01-preview`, `2020-01-01`, `2021-06-01`, `2025-05-04-preview`].

* `parent_id` - (Required) The ID of the azure resource in which this resource is created. The allowed values are:  
  `{any azure resource id}`

* `name` - (Required) Specifies the name of the azure resource. Changing this forces a new resource to be created.

* `body` - (Required) Specifies the configuration of the resource. More information about the arguments in `body` can be found in the [Microsoft documentation](https://learn.microsoft.com/en-us/azure/templates/Microsoft.Security/assessments?pivots=deployment-language-terraform).

For other arguments, please refer to the [azapi_resource](https://registry.terraform.io/providers/Azure/azapi/latest/docs/resources/resource) documentation.

## Import

 ```shell
 # Azure resource can be imported using the resource id, e.g.
 terraform import azapi_resource.example {any azure resource id}/providers/Microsoft.Security/assessments/{resourceName}
 
 # It also supports specifying API version by using the resource id with api-version as a query parameter, e.g.
 terraform import azapi_resource.example {any azure resource id}/providers/Microsoft.Security/assessments/{resourceName}?api-version=2025-05-04-preview
 ```
