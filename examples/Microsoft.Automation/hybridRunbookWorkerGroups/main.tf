terraform {
  required_providers {
    azapi = {
      source = "azure/azapi"
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
  location = "westus"
}

resource "azurerm_automation_account" "test" {
  name                = "myAccount"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  sku_name            = "Basic"
}

resource "azapi_resource" "test" {
  type      = "Microsoft.Automation/automationAccounts/hybridRunbookWorkerGroups@2021-06-22"
  name      = "myWorkerGroup"
  parent_id = azurerm_automation_account.test.id

  body = jsonencode({
  })
}

data "azurerm_virtual_machine" "test" {
  name                = "myVM"
  resource_group_name = "myVirtualMachines"
}

resource "azapi_resource" "test1" {
  type      = "Microsoft.Automation/automationAccounts/hybridRunbookWorkerGroups/hybridRunbookWorkers@2021-06-22"
  name      = "myRunbookWorker"
  parent_id = azapi_resource.test.id

  body = jsonencode({
    properties = {
      vmResourceId = data.azurerm_virtual_machine.test.id
    }
  })
}