terraform {
  required_providers {
    azurermg = {
      source = "ms-henglu/azurermg"
    }
  }
}

provider "azurerm" {
  features {}
}

provider "azurermg" {
}

resource "azurerm_resource_group" "test" {
  name     = "acctest477"
  location = "West Europe"
}

resource "azurerm_storage_account" "test" {
  name                     = "acctest477"
  resource_group_name      = azurerm_resource_group.test.name
  location                 = azurerm_resource_group.test.location
  account_tier             = "Standard"
  account_replication_type = "LRS"
  account_kind             = "BlobStorage"
}

resource "azurerm_storage_data_lake_gen2_filesystem" "test" {
  name               = "acctest477"
  storage_account_id = azurerm_storage_account.test.id
}

resource "azurerm_synapse_workspace" "test" {
  name                                 = "acctest477"
  resource_group_name                  = azurerm_resource_group.test.name
  location                             = azurerm_resource_group.test.location
  storage_data_lake_gen2_filesystem_id = azurerm_storage_data_lake_gen2_filesystem.test.id
  sql_administrator_login              = "sqladminuser"
  sql_administrator_login_password     = "H@Sh1CoR3!"
}

resource "azurerm_synapse_sql_pool" "test" {
  name                 = "acctest477"
  synapse_workspace_id = azurerm_synapse_workspace.test.id
  sku_name             = "DW100c"
  create_mode          = "Default"
}


resource "azurermg_resource" "test" {
  url = "${azurerm_synapse_sql_pool.test.id}/workloadGroups/smallrc"
  api_version = "2020-12-01"
  body = <<BODY
{
    "properties": {
        "importance": "normal",
        "maxResourcePercent": 100,
        "maxResourcePercentPerRequest": 3,
        "minResourcePercent": 0,
        "minResourcePercentPerRequest": 3,
        "queryExecutionTimeout": 0
    }
}
BODY
}
