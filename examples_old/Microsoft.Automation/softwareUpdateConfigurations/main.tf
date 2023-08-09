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

resource "azapi_resource" "auto" {
  type      = "Microsoft.Automation/automationAccounts@2021-06-22"
  name      = "myAccount"
  parent_id = azurerm_resource_group.test.id

  location = azurerm_resource_group.test.location
  body = jsonencode({
    properties = {
      sku = {
        name = "Basic"
      }
    }
  })
}

resource "azapi_resource" "sfuc" {
  type      = "Microsoft.Automation/automationAccounts/softwareUpdateConfigurations@2019-06-01"
  name      = "myConfig"
  parent_id = azapi_resource.auto.id

  body = jsonencode({
    properties = {
      scheduleInfo = {
        startTime               = "2022-04-08T07:00:00+00:00"
        isEnabled               = true
        interval                = 1
        frequency               = "Day"
        timeZone                = "Etc/UTC"
      }
      updateConfiguration = {
        operatingSystem = "Windows"
        windows = {
          includedUpdateClassifications = "Critical"
          rebootSetting                 = "IfRequired"
        }
        targets = {
          azureQueries = [
            {
              locations = ["westus"]
              scope = [
                azurerm_resource_group.test.id
              ]
            }
          ]
        }
      }
    }
  })
}