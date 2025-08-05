terraform {
  required_providers {
    azapi = {
      source = "Azure/azapi"
    }
  }
}

variable "resource_name" {
  type    = string
  default = "acctest0001"
}

variable "location" {
  type    = string
  default = "westus"
}

data "azapi_client_config" "current" {}

resource "azapi_resource" "resourceGroup" {
  type     = "Microsoft.Resources/resourceGroups@2020-06-01"
  name     = var.resource_name
  location = var.location
}

resource "azapi_resource" "namespace" {
  type      = "Microsoft.EventHub/namespaces@2024-01-01"
  parent_id = azapi_resource.resourceGroup.id
  name      = "${var.resource_name}-EHN-unique"
  location  = var.location
  body = {
    properties = {
      disableLocalAuth     = false
      isAutoInflateEnabled = false
      minimumTlsVersion    = "1.2"
      publicNetworkAccess  = "Enabled"
    }
    sku = {
      capacity = 1
      name     = "Basic"
      tier     = "Basic"
    }
  }
}

resource "azapi_resource" "eventhub" {
  type      = "Microsoft.EventHub/namespaces/eventhubs@2024-01-01"
  parent_id = azapi_resource.namespace.id
  name      = "${var.resource_name}-EH-unique"
  body = {
    properties = {
      messageRetentionInDays = 1
      partitionCount         = 2
      status                 = "Active"
    }
  }
}

resource "azapi_resource" "authorizationRule" {
  type      = "Microsoft.EventHub/namespaces/authorizationRules@2024-01-01"
  parent_id = azapi_resource.namespace.id
  name      = "example"
  body = {
    properties = {
      rights = ["Listen", "Send", "Manage"]
    }
  }
}

resource "azapi_resource" "diagnosticSetting" {
  type      = "Microsoft.AADIAM/diagnosticSettings@2017-04-01"
  parent_id = "/"
  name      = "${var.resource_name}-DS-unique"
  body = {
    properties = {
      eventHubAuthorizationRuleId = azapi_resource.authorizationRule.id
      eventHubName                = azapi_resource.eventhub.name
      logs = [
        {
          category = "RiskyUsers"
          enabled  = true
        },
        {
          category = "ServicePrincipalSignInLogs"
          enabled  = true
        },
        {
          category = "SignInLogs"
          enabled  = true
        },
        {
          category = "B2CRequestLogs"
          enabled  = true
        },
        {
          category = "UserRiskEvents"
          enabled  = true
        },
        {
          category = "NonInteractiveUserSignInLogs"
          enabled  = true
        },
        {
          category = "AuditLogs"
          enabled  = true
        }
      ]
    }
  }
}
