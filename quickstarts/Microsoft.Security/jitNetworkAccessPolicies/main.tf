terraform {
  required_providers {
    azapi = {
      source = "Azure/azapi"
    }
  }
}

provider "azapi" {
}

provider "azurerm" {
  features {}
}

locals {
  public_key = "your_public_key"
}

resource "azurerm_resource_group" "test" {
  name     = "myResourceGroup"
  location = "west europe"
}

resource "azurerm_virtual_network" "test" {
  name                = "myVnet"
  address_space       = ["10.0.0.0/16"]
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_subnet" "test" {
  name                 = "internal"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefixes     = ["10.0.2.0/24"]
}

resource "azurerm_network_interface" "test" {
  name                = "myNetworkInterface"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  ip_configuration {
    name                          = "internal"
    subnet_id                     = azurerm_subnet.test.id
    private_ip_address_allocation = "Dynamic"
  }
}

resource "azurerm_linux_virtual_machine" "test" {
  name                = "myvm"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  size                = "Standard_F2"
  admin_username      = "adminuser"
  network_interface_ids = [
    azurerm_network_interface.test.id,
  ]

  admin_ssh_key {
    username   = "adminuser"
    public_key = local.public_key
  }

  os_disk {
    caching              = "ReadWrite"
    storage_account_type = "Standard_LRS"
  }

  source_image_reference {
    publisher = "Canonical"
    offer     = "UbuntuServer"
    sku       = "16.04-LTS"
    version   = "latest"
  }
}

resource "azapi_resource" "test" {
  type = "Microsoft.Security/locations/jitNetworkAccessPolicies@2020-01-01"
  name = "myPolicy"
  // `Microsoft.Security/locations` is not a valid resource type to be managed, so we must build its id like the following
  parent_id = "${azurerm_resource_group.test.id}/providers/Microsoft.Security/locations/westeurope"
  body = {
    properties = {
      virtualMachines = [
        {
          id = azurerm_linux_virtual_machine.test.id
          ports = [
            {
              maxRequestAccessDuration   = "PT3H"
              number                     = 22
              protocol                   = "*"
              allowedSourceAddressPrefix = "192.127.0.2"
            }
          ]
        }
      ]
    }
    kind = "Basic"
  }
}
