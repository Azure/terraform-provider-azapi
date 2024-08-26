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

resource "azapi_resource" "scheduledQueryRule" {
  type      = "Microsoft.Insights/scheduledQueryRules@2021-08-01"
  parent_id = azapi_resource.resourceGroup.id
  name      = var.resource_name
  location  = var.location
  body = {
    kind = "LogAlert"
    properties = {
      autoMitigate                          = false
      checkWorkspaceAlertsStorageConfigured = false
      criteria = {
        allOf = [
          {
            dimensions      = null
            operator        = "Equal"
            query           = " requests\n| summarize CountByCountry=count() by client_CountryOrRegion\n"
            threshold       = 5
            timeAggregation = "Count"
          },
        ]
      }
      enabled             = true
      evaluationFrequency = "PT5M"
      scopes = [
        azapi_resource.component.id,
      ]
      severity            = 3
      skipQueryValidation = false
      targetResourceTypes = null
      windowSize          = "PT5M"
    }

  }
  schema_validation_enabled = false
  response_export_values    = ["*"]
}

