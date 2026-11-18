---
subcategory: "Microsoft.NetworkSecurity - Azure Adaptive Virtual Patching"
page_title: "neuroShields/protectedResources"
description: |-
  Manages a Adaptive Shield Protected Resource.
---

# Microsoft.NetworkSecurity/neuroShields/protectedResources - Adaptive Shield Protected Resource

This article demonstrates how to use `azapi` provider to manage the Adaptive Shield Protected Resource resource in Azure.



## Example Usage

### basic

```hcl
terraform {
  required_providers {
    azapi = {
      source = "Azure/azapi"
    }
    tls = {
      source  = "hashicorp/tls"
      version = "~> 4.1"
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
  default = "centraluseuap"
}

resource "tls_private_key" "ssh" {
  algorithm = "RSA"
  rsa_bits  = 4096
}

resource "azapi_resource" "resourceGroup" {
  type     = "Microsoft.Resources/resourceGroups@2020-06-01"
  name     = var.resource_name
  location = var.location
}

resource "azapi_resource" "neuroShield" {
  type                      = "Microsoft.NetworkSecurity/neuroShields@2026-10-01-preview"
  parent_id                 = azapi_resource.resourceGroup.id
  name                      = var.resource_name
  location                  = var.location
  schema_validation_enabled = false

  body = {
    properties = {
      responseMode      = "Manual"
      severityThreshold = "Medium"
      description       = "Adaptive virtual patching example"
    }
  }
}

resource "azapi_resource" "virtualNetwork" {
  type      = "Microsoft.Network/virtualNetworks@2022-07-01"
  parent_id = azapi_resource.resourceGroup.id
  name      = "${var.resource_name}-vnet"
  location  = var.location

  body = {
    properties = {
      addressSpace = {
        addressPrefixes = ["10.0.0.0/16"]
      }
    }
  }

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
    }
  }
}

resource "azapi_resource" "networkInterface" {
  type      = "Microsoft.Network/networkInterfaces@2022-07-01"
  parent_id = azapi_resource.resourceGroup.id
  name      = "${var.resource_name}-nic"
  location  = var.location

  body = {
    properties = {
      ipConfigurations = [
        {
          name = "internal"
          properties = {
            privateIPAllocationMethod = "Dynamic"
            subnet = {
              id = azapi_resource.subnet.id
            }
          }
        }
      ]
    }
  }
}

resource "azapi_resource" "virtualMachine" {
  type      = "Microsoft.Compute/virtualMachines@2023-03-01"
  parent_id = azapi_resource.resourceGroup.id
  name      = var.resource_name
  location  = var.location

  identity {
    type = "SystemAssigned"
  }

  body = {
    properties = {
      hardwareProfile = {
        vmSize = "Standard_D2s_v6"
      }
      networkProfile = {
        networkInterfaces = [
          {
            id = azapi_resource.networkInterface.id
          }
        ]
      }
      osProfile = {
        adminUsername = "azureuser"
        computerName  = var.resource_name
        linuxConfiguration = {
          disablePasswordAuthentication = true
          ssh = {
            publicKeys = [
              {
                path    = "/home/azureuser/.ssh/authorized_keys"
                keyData = tls_private_key.ssh.public_key_openssh
              }
            ]
          }
        }
      }
      storageProfile = {
        imageReference = {
          publisher = "Canonical"
          offer     = "ubuntu-26_04-lts"
          sku       = "server"
          version   = "latest"
        }
        osDisk = {
          caching      = "ReadWrite"
          createOption = "FromImage"
          managedDisk = {
            storageAccountType = "Standard_LRS"
          }
        }
      }
    }
  }
}

resource "azapi_resource" "protectedResource" {
  type                      = "Microsoft.NetworkSecurity/neuroShields/protectedResources@2026-10-01-preview"
  parent_id                 = azapi_resource.neuroShield.id
  name                      = var.resource_name
  schema_validation_enabled = false

  body = {
    properties = {
      targetResourceId = azapi_resource.virtualMachine.id
      targetIdentity = {
        type        = "SystemAssigned"
        principalId = azapi_resource.virtualMachine.identity[0].principal_id
        tenantId    = azapi_resource.virtualMachine.identity[0].tenant_id
      }
    }
  }

  # The target resource and identity are create-only; changes require replacing the registration.
  # https://github.com/Azure/azure-rest-api-specs/blob/e0f840aaf5c7c2daeec7da8abcd310dd7f47f56b/specification/networksecurity/resource-manager/Microsoft.NetworkSecurity/neuroShield/models.tsp#L278-L311
  replace_triggers_refs = [
    "properties.targetResourceId",
    "properties.targetIdentity",
  ]
}

```



## Arguments Reference

The following arguments are supported:

* `type` - (Required) The type of the resource. This should be set to `Microsoft.NetworkSecurity/neuroShields/protectedResources@api-version`. The available api-versions for this resource are: [`2026-10-01-preview`].

* `parent_id` - (Required) The ID of the azure resource in which this resource is created. The allowed values are:  
  `/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.NetworkSecurity/neuroShields/{resourceName}`

* `name` - (Required) Specifies the name of the azure resource. Changing this forces a new resource to be created.

* `body` - (Required) Specifies the configuration of the resource. More information about the arguments in `body` can be found in the [Microsoft documentation](https://learn.microsoft.com/en-us/azure/templates/Microsoft.NetworkSecurity/neuroShields/protectedResources?pivots=deployment-language-terraform).

For other arguments, please refer to the [azapi_resource](https://registry.terraform.io/providers/Azure/azapi/latest/docs/resources/resource) documentation.

## Import

 ```shell
 # Azure resource can be imported using the resource id, e.g.
 terraform import azapi_resource.example /subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.NetworkSecurity/neuroShields/{resourceName}/protectedResources/{resourceName}
 
 # It also supports specifying API version by using the resource id with api-version as a query parameter, e.g.
 terraform import azapi_resource.example /subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.NetworkSecurity/neuroShields/{resourceName}/protectedResources/{resourceName}?api-version=2026-10-01-preview
 ```
