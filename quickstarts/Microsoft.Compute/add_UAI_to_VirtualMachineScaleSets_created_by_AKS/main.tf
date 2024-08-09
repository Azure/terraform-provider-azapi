terraform {
  required_providers {
    azapi = {
      source = "Azure/azapi"
    }
  }
}

provider "azurerm" {
  features {}
}

provider "azapi" {
}

resource "azurerm_resource_group" "test" {
  name     = "myResourceGroup"
  location = "westeurope"
}

resource "azurerm_kubernetes_cluster" "test" {
  name                = "mycluster"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  dns_prefix          = "mycluster"

  default_node_pool {
    name       = "default"
    node_count = 1
    vm_size    = "Standard_D2_v2"
  }

  identity {
    type = "SystemAssigned"
  }
}

resource "azurerm_user_assigned_identity" "test" {
  name                = "myidentity"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
}

data "azurerm_virtual_machine_scale_set" "test" {
  name                = "aks-default-34472765-vmss"
  resource_group_name = "MC_${azurerm_kubernetes_cluster.test.name}_${azurerm_kubernetes_cluster.test.name}_${azurerm_resource_group.test.location}"
}

resource "azapi_resource_action" "test" {
  type        = "Microsoft.Compute/virtualMachineScaleSets@2022-03-01"
  resource_id = data.azurerm_virtual_machine_scale_set.test.id
  // omit `action` field or set it to empty string like `action = ""`, to make request towards the resource
  method    = "PATCH"
  body = {
    identity = {
      type = "UserAssigned"
      userAssignedIdentities = {
        (azurerm_user_assigned_identity.test.id) = {}
      }
    }
  }
}
