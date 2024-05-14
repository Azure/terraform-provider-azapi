terraform {
  required_providers {
    azapi = {
      source = "Azure/azapi"
    }
    azurerm = {
      source = "hashicorp/azurerm"
    }
  }
}

provider "azurerm" {
  features {
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

data "azurerm_client_config" "current" {
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

resource "azapi_resource" "connection" {
  type      = "Microsoft.Automation/automationAccounts/connections@2020-01-13-preview"
  parent_id = azapi_resource.automationAccount.id
  name      = var.resource_name
  body = {
    properties = {
      connectionType = {
        name = "AzureServicePrincipal"
      }
      description = ""
      fieldDefinitionValues = {
        ApplicationId         = "00000000-0000-0000-0000-000000000000"
        CertificateThumbprint = "AEB97B81A68E8988850972916A8B8B6CD8F39813\n"
        SubscriptionId        = data.azurerm_client_config.current.subscription_id
        TenantId              = data.azurerm_client_config.current.tenant_id
      }
    }
  }
  schema_validation_enabled = false
  response_export_values    = ["*"]
}

