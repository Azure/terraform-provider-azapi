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

resource "azapi_resource" "signalR" {
  type      = "Microsoft.SignalRService/signalR@2023-02-01"
  parent_id = azapi_resource.resourceGroup.id
  name      = var.resource_name
  location  = var.location
  body = {
    properties = {
      cors = {
      }
      disableAadAuth   = false
      disableLocalAuth = false
      features = [
        {
          flag  = "ServiceMode"
          value = "Default"
        },
        {
          flag  = "EnableConnectivityLogs"
          value = "False"
        },
        {
          flag  = "EnableMessagingLogs"
          value = "False"
        },
        {
          flag  = "EnableLiveTrace"
          value = "False"
        },
      ]
      publicNetworkAccess = "Enabled"
      resourceLogConfiguration = {
        categories = [
          {
            enabled = "false"
            name    = "MessagingLogs"
          },
          {
            enabled = "false"
            name    = "ConnectivityLogs"
          },
          {
            enabled = "false"
            name    = "HttpRequestLogs"
          },
        ]
      }
      serverless = {
        connectionTimeoutInSeconds = 30
      }
      tls = {
        clientCertEnabled = false
      }
      upstream = {
        templates = [
        ]
      }
    }
    sku = {
      capacity = 1
      name     = "Standard_S1"
    }
  }
  schema_validation_enabled = false
  response_export_values    = ["*"]
}

