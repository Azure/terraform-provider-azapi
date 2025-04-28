---
subcategory: "Microsoft.Maintenance - Azure Maintenance"
page_title: "configurationAssignments"
description: |-
  Manages a Maintenance Assignment.
---

# Microsoft.Maintenance/configurationAssignments - Maintenance Assignment

This article demonstrates how to use `azapi` provider to manage the Maintenance Assignment resource in Azure.

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

resource "azapi_resource" "maintenanceConfiguration" {
  type      = "Microsoft.Maintenance/maintenanceConfigurations@2022-07-01-preview"
  parent_id = azapi_resource.resourceGroup.id
  name      = var.resource_name
  location  = var.location
  body = {
    properties = {
      extensionProperties = {
      }
      maintenanceScope = "SQLDB"
      namespace        = "Microsoft.Maintenance"
      visibility       = "Custom"
    }
  }
  schema_validation_enabled = false
  response_export_values    = ["*"]
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
      additionalCapabilities = {
      }
      applicationProfile = {
        galleryApplications = [
        ]
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
        networkInterfaces = [
          {
            id = azapi_resource.networkInterface.id
            properties = {
              primary = true
            }
          },
        ]
      }
      osProfile = {
        adminPassword            = "P@$$w0rd1234!"
        adminUsername            = "adminuser"
        allowExtensionOperations = true
        computerName             = var.resource_name
        linuxConfiguration = {
          disablePasswordAuthentication = false
          patchSettings = {
            assessmentMode = "ImageDefault"
            patchMode      = "ImageDefault"
          }
          provisionVMAgent = true
          ssh = {
            publicKeys = [
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
  schema_validation_enabled = false
  response_export_values    = ["*"]
}

resource "azapi_resource" "configurationAssignment" {
  type      = "Microsoft.Maintenance/configurationAssignments@2022-07-01-preview"
  parent_id = azapi_resource.virtualMachine.id
  name      = var.resource_name
  body = {
    location = "westeurope"
    properties = {
      maintenanceConfigurationId = azapi_resource.maintenanceConfiguration.id
      resourceId                 = azapi_resource.virtualMachine.id
    }
  }
  ignore_casing             = true
  schema_validation_enabled = false
  response_export_values    = ["*"]
}


```



## Arguments Reference

The following arguments are supported:

* `type` - (Required) The type of the resource. This should be set to `Microsoft.Maintenance/configurationAssignments@api-version`. The available api-versions for this resource are: [`2018-06-01-preview`, `2020-04-01`, `2020-07-01-preview`, `2021-04-01-preview`, `2021-05-01`, `2021-09-01-preview`, `2022-07-01-preview`, `2022-11-01-preview`, `2023-04-01`, `2023-09-01-preview`, `2023-10-01-preview`].

* `parent_id` - (Required) The ID of the azure resource in which this resource is created. The allowed values are:  
  `{any azure resource id}`

* `name` - (Required) Specifies the name of the azure resource. Changing this forces a new resource to be created.

* `body` - (Required) Specifies the configuration of the resource. More information about the arguments in `body` can be found in the [Microsoft documentation](https://learn.microsoft.com/en-us/azure/templates/Microsoft.Maintenance/configurationAssignments?pivots=deployment-language-terraform).

## Import

 ```shell
 # Azure resource can be imported using the resource id, e.g.
 terraform import azapi_resource.example {any azure resource id}/providers/Microsoft.Maintenance/configurationAssignments/{resourceName}
 
 # It also supports specifying API version by using the resource id with api-version as a query parameter, e.g.
 terraform import azapi_resource.example {any azure resource id}/providers/Microsoft.Maintenance/configurationAssignments/{resourceName}?api-version=2023-10-01-preview
 ```
