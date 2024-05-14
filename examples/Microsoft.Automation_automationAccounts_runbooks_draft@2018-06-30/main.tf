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

data "azapi_resource_id" "draft" {
  type      = "Microsoft.Automation/automationAccounts/runbooks/draft@2018-06-30"
  parent_id = azapi_resource.runbook.id
  name      = "content"
}

resource "azapi_resource_action" "draft" {
  type        = "Microsoft.Automation/automationAccounts/runbooks/draft@2018-06-30"
  resource_id = data.azapi_resource_id.draft.id
  method      = "PUT"
  body = {
    name = "content"
  }
  response_export_values = ["*"]
}

