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

resource "azurerm_automation_runbook" "test" {
  name                    = "Get-AzureVMTutorial"
  location                = azurerm_resource_group.test.location
  resource_group_name     = azurerm_resource_group.test.name
  automation_account_name = azurerm_automation_account.test.name
  log_verbose             = "true"
  log_progress            = "true"
  description             = "This is an example runbook"
  runbook_type            = "PowerShellWorkflow"

  publish_content_link {
    uri = "https://raw.githubusercontent.com/Azure/azure-quickstart-templates/c4935ffb69246a6058eb24f54640f53f69d3ac9f/101-automation-runbook-getvms/Runbooks/Get-AzureVMTutorial.ps1"
  }
}

resource "azapi_resource_action" "job" {
  type        = "Microsoft.Automation/automationAccounts/jobs@2019-06-01"
  resource_id = "${azurerm_automation_account.test.id}/jobs/myJob"
  method      = "PUT"
  body = jsonencode({
    properties = {
      runbook = {
        name = azurerm_automation_runbook.test.name
      }
    }
  })
}
