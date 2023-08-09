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

data "azurerm_client_config" "test" {}

resource "azurerm_resource_group" "test" {
  name     = "acctest5076"
  location = "West Europe"
}

resource "azurerm_application_insights" "test" {
  name                = "acctest5076"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  application_type    = "web"
}

resource "azurerm_key_vault" "test" {
  name                = "acctest5076"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  tenant_id           = data.azurerm_client_config.test.tenant_id
  sku_name            = "premium"
}

resource "azurerm_storage_account" "test" {
  name                     = "acctest5076"
  location                 = azurerm_resource_group.test.location
  resource_group_name      = azurerm_resource_group.test.name
  account_tier             = "Standard"
  account_replication_type = "GRS"
}

resource "azurerm_machine_learning_workspace" "test" {
  name                    = "acctest5076"
  location                = azurerm_resource_group.test.location
  resource_group_name     = azurerm_resource_group.test.name
  application_insights_id = azurerm_application_insights.test.id
  key_vault_id            = azurerm_key_vault.test.id
  storage_account_id      = azurerm_storage_account.test.id

  identity {
    type = "SystemAssigned"
  }
}


resource "azapi_resource" "test" {
  type      = "Microsoft.MachineLearningServices/workspaces/computes@2021-07-01"
  name      = "acctest6032"
  parent_id = azurerm_machine_learning_workspace.test.id

  location = "eastus"
  body = jsonencode({
    properties = {
      computeType = "ComputeInstance"
      properties = {
        vmSize = "STANDARD_NC6"
      }
    }
  })
}
