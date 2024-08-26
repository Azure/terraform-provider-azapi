terraform {
  required_providers {
    azapi = {
      source = "Azure/azapi"
    }
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
  type     = "Microsoft.Resources/resourceGroups@2020-06-01"
  name     = var.resource_name
  location = var.location
}

resource "azapi_resource" "automationAccount" {
  type      = "Microsoft.Automation/automationAccounts@2021-06-22"
  parent_id = azapi_resource.resourceGroup.id
  name      = var.resource_name
  location  = var.location
  body = {
    properties = {
      encryption = {
        keySource = "Microsoft.Automation"
      }
      publicNetworkAccess = true
      sku = {
        name = "Basic"
      }
    }
  }
  schema_validation_enabled = false
  response_export_values    = ["*"]
}

resource "azapi_resource" "runbook" {
  type      = "Microsoft.Automation/automationAccounts/runbooks@2019-06-01"
  parent_id = azapi_resource.automationAccount.id
  name      = "Get-AzureVMTutorial"
  location  = var.location
  body = {
    properties = {
      description = "This is a test runbook for terraform acceptance test"
      draft = {
      }
      logActivityTrace = 0
      logProgress      = true
      logVerbose       = true
      runbookType      = "PowerShell"
    }
  }
  schema_validation_enabled = false
  response_export_values    = ["*"]
}

resource "azapi_resource_action" "public_runbook" {
  type        = "Microsoft.Automation/automationAccounts/runbooks@2019-06-01"
  resource_id = azapi_resource.runbook.id
  action      = "publish"
}

resource "azapi_resource" "schedule" {
  type      = "Microsoft.Automation/automationAccounts/schedules@2020-01-13-preview"
  parent_id = azapi_resource.automationAccount.id
  name      = var.resource_name
  body = {
    properties = {
      description = ""
      frequency   = "OneTime"
      startTime   = "2024-07-05T08:51:00+00:00"
      timeZone    = "Etc/UTC"
    }
  }
  schema_validation_enabled = false
  response_export_values    = ["*"]
}

resource "azapi_resource" "jobSchedule" {
  type      = "Microsoft.Automation/automationAccounts/jobSchedules@2020-01-13-preview"
  parent_id = azapi_resource.automationAccount.id
  name      = "194a324f-9e3d-43ee-1234-c968b797edd5"
  body = {
    properties = {
      runbook = {
        name = azapi_resource.runbook.name
      }
      schedule = {
        name = azapi_resource.schedule.name
      }
    }
  }
  schema_validation_enabled = false
  response_export_values    = ["*"]
  depends_on                = [azapi_resource_action.public_runbook]
}

