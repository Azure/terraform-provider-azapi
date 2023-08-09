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

resource "azapi_resource" "account" {
  type = "Microsoft.LabServices/labaccounts@2018-10-15"
  name = "myAccount"
  parent_id = azurerm_resource_group.test.id
  location = azurerm_resource_group.test.location
  body = jsonencode({
    properties = {
      enabledRegionSelection = false
    }
  })
}

resource "azapi_resource" "lab" {
  type = "Microsoft.LabServices/labaccounts/labs@2018-10-15"
  name = "myLab"
  parent_id = azapi_resource.account.id

  body = jsonencode({
    properties = {
      maxUsersInLab = 10
      userAccessMode = "Restricted"
    }
  })
}