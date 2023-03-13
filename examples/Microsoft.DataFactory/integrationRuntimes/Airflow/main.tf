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

resource "azurerm_resource_group" "example_rg" {
  name     = "example-resource-group"
  location = "West Europe"
}

resource "azurerm_data_factory" "example_adf" {
  name                = "example"
  location            = azurerm_resource_group.example_rg.location
  resource_group_name = azurerm_resource_group.example_rg.name
}

# This feature is in public preview at time of writing, hence the `schema_validation_enabled = false`
resource "azapi_resource" "example_ir" {
  type                      = "Microsoft.DataFactory/factories/integrationRuntimes@2018-06-01"
  name                      = "example"
  parent_id                 = azurerm_data_factory.example_adf.id
  schema_validation_enabled = false

  body = jsonencode({
    properties = {
      type = "Airflow"
      typeProperties = {
        computeProperties = {
          location    = "West Europe"
          computeSize = "Small"
          extraNodes  = 0
        }
        airflowProperties = {
          version                  = "2.2.2"
          enableAADIntegration     = true
          airflowRequiredArguments = ["airflow.providers.microsoft.azure"]
        }
      }
    }
  })
}
