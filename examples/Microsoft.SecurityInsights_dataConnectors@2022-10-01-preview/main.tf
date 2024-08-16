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

resource "azapi_resource" "workspace" {
  type      = "Microsoft.OperationalInsights/workspaces@2022-10-01"
  parent_id = azapi_resource.resourceGroup.id
  name      = var.resource_name
  location  = var.location
  body = {
    properties = {
      features = {
        disableLocalAuth                            = false
        enableLogAccessUsingOnlyResourcePermissions = true
      }
      publicNetworkAccessForIngestion = "Enabled"
      publicNetworkAccessForQuery     = "Enabled"
      retentionInDays                 = 30
      sku = {
        name = "PerGB2018"
      }
      workspaceCapping = {
        dailyQuotaGb = -1
      }
    }
  }
  schema_validation_enabled = false
  response_export_values    = ["*"]
}

resource "azapi_resource" "onboardingState" {
  type      = "Microsoft.SecurityInsights/onboardingStates@2023-06-01-preview"
  parent_id = azapi_resource.workspace.id
  name      = "default"
  body = {
    properties = {
      customerManagedKey = false
    }
  }
}

resource "azapi_resource" "dataConnector" {
  type      = "Microsoft.SecurityInsights/dataConnectors@2022-10-01-preview"
  parent_id = azapi_resource.workspace.id
  name      = var.resource_name
  body = {
    kind = "MicrosoftThreatIntelligence"
    properties = {
      dataTypes = {
        bingSafetyPhishingURL = {
          lookbackPeriod = ""
          state          = "Disabled"
        }
        microsoftEmergingThreatFeed = {
          lookbackPeriod = "1970-01-01T00:00:00Z"
          state          = "enabled"
        }
      }
      tenantId = data.azurerm_client_config.current.tenant_id
    }
  }
  schema_validation_enabled = false
  response_export_values    = ["*"]
  depends_on                = [azapi_resource.onboardingState]
}

