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

resource "azapi_resource" "automationRule" {
  type      = "Microsoft.SecurityInsights/automationRules@2022-10-01-preview"
  parent_id = azapi_resource.workspace.id
  name      = "3b862818-ad7b-409e-83be-8812f2a06d37"
  body = {
    properties = {
      actions = [
        {
          actionConfiguration = {
            classification        = ""
            classificationComment = ""
            classificationReason  = ""
            severity              = ""
            status                = "Active"
          }
          actionType = "ModifyProperties"
          order      = 1
        },
      ]
      displayName = "acctest-SentinelAutoRule-230630033910945846"
      order       = 1
      triggeringLogic = {
        isEnabled    = true
        triggersOn   = "Incidents"
        triggersWhen = "Created"
      }
    }
  }
  schema_validation_enabled = false
  response_export_values    = ["*"]
  depends_on                = [azapi_resource.onboardingState]
}

