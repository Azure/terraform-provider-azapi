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
  body = jsonencode({
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
  })
}

resource "azapi_resource" "routeTable" {
  type      = "Microsoft.Network/routeTables@2023-04-01"
  parent_id = azapi_resource.resourceGroup.id
  name      = var.resource_name
  location  = azapi_resource.resourceGroup.location
  body = jsonencode({
    properties = {
      disableBgpRoutePropagation = false
    }
  })
}

resource "azapi_resource" "virtualNetwork" {
  type      = "Microsoft.Network/virtualNetworks@2023-04-01"
  parent_id = azapi_resource.resourceGroup.id
  name      = var.resource_name
  location  = azapi_resource.resourceGroup.location
  body = jsonencode({
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
  })
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
  body = jsonencode({
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
  })

  timeouts {
    create = "300m"
    update = "300m"
    delete = "300m"
  }
}
