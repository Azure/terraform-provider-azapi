---
subcategory: "Microsoft.Network - Various Networking Services"
page_title: "networkManagers/verifierWorkspace/reachabilityAnalysisIntent"
description: |-
  Manages a Network Managers Verifier Workspace Reachability Analysis Intent.
---

# Microsoft.Network/networkManagers/verifierWorkspace/reachabilityAnalysisIntent - Network Managers Verifier Workspace Reachability Analysis Intent

This article demonstrates how to use `azapi` provider to manage the Network Managers Verifier Workspace Reachability Analysis Intent resource in Azure.

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

# Virtual Machine administrator credentials
variable "admin_username" {
  type        = string
  description = "The admin username for the virtual machine"
}

variable "admin_password" {
  type        = string
  description = "The admin password for the virtual machine"
  sensitive   = true
}

data "azapi_client_config" "current" {}

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
  lifecycle {
    ignore_changes = [body.properties.subnets]
  }
  schema_validation_enabled = false
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
  ignore_missing_property = true # ignore serviceEndpointPolicies NOT in GET subnet response
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
          name                    = "myosdisk1"
          writeAcceleratorEnabled = false
        }
      }
    }
  }
  ignore_missing_property = true # ignore adminPassword as NOT part of GET vm response
}

resource "azapi_resource" "networkManager" {
  type      = "Microsoft.Network/networkManagers@2022-09-01"
  parent_id = azapi_resource.resourceGroup.id
  name      = var.resource_name
  location  = var.location
  body = {
    properties = {
      description = ""
      networkManagerScopeAccesses = [
        "SecurityAdmin",
      ]
      networkManagerScopes = {
        managementGroups = [
        ]
        subscriptions = [
          "/subscriptions/${data.azapi_client_config.current.subscription_id}",
        ]
      }
    }
  }
  retry = {
    error_message_regex = ["CannotDeleteResource"]
  }
}

resource "azapi_resource" "verifierWorkspace" {
  type      = "Microsoft.Network/networkManagers/verifierWorkspaces@2024-01-01-preview"
  parent_id = azapi_resource.networkManager.id
  name      = var.resource_name
  location  = var.location
  body = {
    properties = {
      description = "A sample workspace"
    }
  }

  tags = {
    myTag = "testTag"
  }
}

resource "azapi_resource" "reachabilityAnalysisIntent" {
  type      = "Microsoft.Network/networkManagers/verifierWorkspaces/reachabilityAnalysisIntents@2024-01-01-preview"
  parent_id = azapi_resource.verifierWorkspace.id
  name      = var.resource_name
  body = {
    properties = {
      description           = "A sample reachability analysis intent"
      destinationResourceId = azapi_resource.virtualMachine.id
      ipTraffic = {
        destinationIps = [
          "10.4.0.1",
        ]
        destinationPorts = [
          "0",
        ]
        protocols = [
          "Any",
        ]
        sourceIps = [
          "10.4.0.0",
        ]
        sourcePorts = [
          "0",
        ]
      }
      sourceResourceId = azapi_resource.virtualMachine.id
    }
  }
}

resource "azapi_resource" "reachabilityAnalysisRun" {
  type      = "Microsoft.Network/networkManagers/verifierWorkspaces/reachabilityAnalysisRuns@2024-01-01-preview"
  parent_id = azapi_resource.verifierWorkspace.id
  name      = var.resource_name
  body = {
    properties = {
      description = "A sample reachability analysis run"
      intentId    = azapi_resource.reachabilityAnalysisIntent.id
    }
  }
}
```



## Arguments Reference

The following arguments are supported:

* `type` - (Required) The type of the resource. This should be set to `Microsoft.Network/networkManagers/verifierWorkspace/reachabilityAnalysisIntent@api-version`. The available api-versions for this resource are: [].

* `parent_id` - (Required) The ID of the azure resource in which this resource is created. The allowed values are:  
  

* `name` - (Required) Specifies the name of the azure resource. Changing this forces a new resource to be created.

* `body` - (Required) Specifies the configuration of the resource. More information about the arguments in `body` can be found in the [Microsoft documentation](https://learn.microsoft.com/en-us/azure/templates/Microsoft.Network/networkManagers/verifierWorkspace/reachabilityAnalysisIntent?pivots=deployment-language-terraform).

For other arguments, please refer to the [azapi_resource](https://registry.terraform.io/providers/Azure/azapi/latest/docs/resources/resource) documentation.

## Import

 ```shell
 # Azure resource can be imported using the resource id, e.g.
 terraform import azapi_resource.example 
 
 # It also supports specifying API version by using the resource id with api-version as a query parameter, e.g.
 terraform import azapi_resource.example ?api-version=API_VERSION
 ```
