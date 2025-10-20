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

resource "azapi_resource" "actionGroup" {
  type      = "Microsoft.Insights/actionGroups@2023-01-01"
  parent_id = azapi_resource.resourceGroup.id
  name      = var.resource_name
  location  = "global"
  body = {
    properties = {
      armRoleReceivers = [
      ]
      automationRunbookReceivers = [
      ]
      azureAppPushReceivers = [
      ]
      azureFunctionReceivers = [
      ]
      emailReceivers = [
      ]
      enabled = true
      eventHubReceivers = [
      ]
      groupShortName = "acctestag"
      itsmReceivers = [
      ]
      logicAppReceivers = [
      ]
      smsReceivers = [
      ]
      voiceReceivers = [
      ]
      webhookReceivers = [
      ]
    }
  }
  schema_validation_enabled = false
  response_export_values    = ["*"]
}

resource "azapi_resource" "component" {
  type      = "Microsoft.Insights/components@2020-02-02"
  parent_id = azapi_resource.resourceGroup.id
  name      = var.resource_name
  location  = var.location
  body = {
    kind = "web"
    properties = {
      Application_Type                = "web"
      DisableIpMasking                = false
      DisableLocalAuth                = false
      ForceCustomerStorageForProfiler = false
      RetentionInDays                 = 90
      SamplingPercentage              = 100
      publicNetworkAccessForIngestion = "Enabled"
      publicNetworkAccessForQuery     = "Enabled"
    }
  }
  schema_validation_enabled = false
  response_export_values    = ["*"]
}

resource "azapi_resource" "smartDetectorAlertRule" {
  type      = "microsoft.alertsManagement/smartDetectorAlertRules@2019-06-01"
  parent_id = azapi_resource.resourceGroup.id
  name      = var.resource_name
  location  = "global"
  body = {
    properties = {
      actionGroups = {
        customEmailSubject   = ""
        customWebhookPayload = ""
        groupIds = [
          azapi_resource.actionGroup.id,
        ]
      }
      description = ""
      detector = {
        id = "FailureAnomaliesDetector"
      }
      frequency = "PT1M"
      scope = [
        azapi_resource.component.id,
      ]
      severity = "Sev0"
      state    = "Enabled"
    }
  }
  schema_validation_enabled = false
  response_export_values    = ["*"]
}

