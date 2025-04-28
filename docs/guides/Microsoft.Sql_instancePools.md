---
subcategory: "Microsoft.Sql - Azure SQL Database, Azure SQL Managed Instance, Azure Synapse Analytics"
page_title: "instancePools"
description: |-
  Manages a SQL Instance Pools.
---

# Microsoft.Sql/instancePools - SQL Instance Pools

This article demonstrates how to use `azapi` provider to manage the SQL Instance Pools resource in Azure.

## Example Usage

### default

```hcl
terraform {
  required_providers {
    azapi = {
      source = "Azure/azapi"
    }
    azurerm = {
      source = "hashicorp/azurerm"
    }
  }
}

provider "azurerm" {
  features {
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
  type     = "Microsoft.Resources/resourceGroups@2022-09-01"
  name     = var.resource_name
  location = var.location
}

resource "azapi_resource" "networkSecurityGroup" {
  type      = "Microsoft.Network/networkSecurityGroups@2023-04-01"
  parent_id = azapi_resource.resourceGroup.id
  name      = var.resource_name
  location  = azapi_resource.resourceGroup.location
  body = {
    properties = {
      securityRules = [
        {
          name = "allow_tds_inbound"
          properties = {
            description              = "Allow access to data"
            protocol                 = "TCP"
            sourcePortRange          = "*"
            destinationPortRange     = "1433"
            sourceAddressPrefix      = "VirtualNetwork"
            destinationAddressPrefix = "*"
            access                   = "Allow"
            priority                 = 1000
            direction                = "Inbound"
          }
        },
        {
          name = "allow_redirect_inbound"
          properties = {
            description              = "Allow inbound redirect traffic to Managed Instance inside the virtual network"
            protocol                 = "Tcp"
            sourcePortRange          = "*"
            destinationPortRange     = "11000-11999"
            sourceAddressPrefix      = "VirtualNetwork"
            destinationAddressPrefix = "*"
            access                   = "Allow"
            priority                 = 1100
            direction                = "Inbound"
          }
        },
        {
          name = "allow_geodr_inbound"
          properties = {
            description              = "Allow inbound geodr traffic inside the virtual network"
            protocol                 = "Tcp"
            sourcePortRange          = "*"
            destinationPortRange     = "5022"
            sourceAddressPrefix      = "VirtualNetwork"
            destinationAddressPrefix = "*"
            access                   = "Allow"
            priority                 = 1200
            direction                = "Inbound"
          }
        },
        {
          name = "deny_all_inbound"
          properties = {
            description              = "Deny all other inbound traffic"
            protocol                 = "*"
            sourcePortRange          = "*"
            destinationPortRange     = "*"
            sourceAddressPrefix      = "*"
            destinationAddressPrefix = "*"
            access                   = "Deny"
            priority                 = 4096
            direction                = "Inbound"
          }
        },
        {
          name = "allow_linkedserver_outbound"
          properties = {
            description              = "Allow outbound linkedserver traffic inside the virtual network"
            protocol                 = "Tcp"
            sourcePortRange          = "*"
            destinationPortRange     = "1433"
            sourceAddressPrefix      = "*"
            destinationAddressPrefix = "VirtualNetwork"
            access                   = "Allow"
            priority                 = 1000
            direction                = "Outbound"
          }
        },
        {
          name = "allow_redirect_outbound"
          properties = {
            description              = "Allow outbound redirect traffic to Managed Instance inside the virtual network"
            protocol                 = "Tcp"
            sourcePortRange          = "*"
            destinationPortRange     = "11000-11999"
            sourceAddressPrefix      = "*"
            destinationAddressPrefix = "VirtualNetwork"
            access                   = "Allow"
            priority                 = 1100
            direction                = "Outbound"
          }
        },
        {
          name = "allow_geodr_outbound"
          properties = {
            description              = "Allow outbound geodr traffic inside the virtual network"
            protocol                 = "Tcp"
            sourcePortRange          = "*"
            destinationPortRange     = "5022"
            sourceAddressPrefix      = "*"
            destinationAddressPrefix = "VirtualNetwork"
            access                   = "Allow"
            priority                 = 1200
            direction                = "Outbound"
          }
        },
        {
          name = "deny_all_outbound"
          properties = {
            description              = "Deny all other outbound traffic"
            protocol                 = "*"
            sourcePortRange          = "*"
            destinationPortRange     = "*"
            sourceAddressPrefix      = "*"
            destinationAddressPrefix = "*"
            access                   = "Deny"
            priority                 = 4096
            direction                = "Outbound"
          }
        }
      ]
    }
  }
}

resource "azapi_resource" "routeTable" {
  type      = "Microsoft.Network/routeTables@2023-04-01"
  parent_id = azapi_resource.resourceGroup.id
  name      = var.resource_name
  location  = azapi_resource.resourceGroup.location
  body = {
    properties = {
      disableBgpRoutePropagation = false
    }
  }
}

resource "azapi_resource" "virtualNetwork" {
  type      = "Microsoft.Network/virtualNetworks@2023-04-01"
  parent_id = azapi_resource.resourceGroup.id
  name      = var.resource_name
  location  = azapi_resource.resourceGroup.location
  body = {
    properties = {
      addressSpace = {
        addressPrefixes = ["10.0.0.0/16"]
      }
      subnets = [
        {
          name = "Default"
          properties = {
            addressPrefix = "10.0.0.0/24"
          }
        },
        {
          name = var.resource_name
          properties = {
            addressPrefix = "10.0.1.0/24"
            networkSecurityGroup = {
              id = azapi_resource.networkSecurityGroup.id
            }
            routeTable = {
              id = azapi_resource.routeTable.id
            }
            delegations = [
              {
                name = "miDelegation"
                properties = {
                  serviceName = "Microsoft.Sql/managedInstances"
                }
              }
            ]
          }
        }
      ]
    }
  }
}

data "azapi_resource" "subnet" {
  type      = "Microsoft.Network/virtualNetworks/subnets@2023-04-01"
  parent_id = azapi_resource.virtualNetwork.id
  name      = var.resource_name
}


resource "azapi_resource" "instancePool" {
  type      = "Microsoft.Sql/instancePools@2022-05-01-preview"
  parent_id = azapi_resource.resourceGroup.id
  name      = var.resource_name
  location  = azapi_resource.resourceGroup.location
  body = {
    properties = {
      licenseType = "LicenseIncluded"
      subnetId    = data.azapi_resource.subnet.id
      vCores      = 8
    }
    sku = {
      family = "Gen5"
      name   = "GP_Gen5"
      tier   = "GeneralPurpose"
    }
  }

  timeouts {
    create = "300m"
    update = "300m"
    delete = "300m"
  }
}

```



## Arguments Reference

The following arguments are supported:

* `type` - (Required) The type of the resource. This should be set to `Microsoft.Sql/instancePools@api-version`. The available api-versions for this resource are: [`2018-06-01-preview`, `2020-02-02-preview`, `2020-08-01-preview`, `2020-11-01-preview`, `2021-02-01-preview`, `2021-05-01-preview`, `2021-08-01-preview`, `2021-11-01`, `2021-11-01-preview`, `2022-02-01-preview`, `2022-05-01-preview`, `2022-08-01-preview`, `2022-11-01-preview`, `2023-02-01-preview`, `2023-05-01-preview`, `2023-08-01-preview`, `2024-05-01-preview`].

* `parent_id` - (Required) The ID of the azure resource in which this resource is created. The allowed values are:  
  `/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}`

* `name` - (Required) Specifies the name of the azure resource. Changing this forces a new resource to be created.

* `body` - (Required) Specifies the configuration of the resource. More information about the arguments in `body` can be found in the [Microsoft documentation](https://learn.microsoft.com/en-us/azure/templates/Microsoft.Sql/instancePools?pivots=deployment-language-terraform).

## Import

 ```shell
 # Azure resource can be imported using the resource id, e.g.
 terraform import azapi_resource.example /subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Sql/instancePools/{resourceName}
 
 # It also supports specifying API version by using the resource id with api-version as a query parameter, e.g.
 terraform import azapi_resource.example /subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Sql/instancePools/{resourceName}?api-version=2024-05-01-preview
 ```
