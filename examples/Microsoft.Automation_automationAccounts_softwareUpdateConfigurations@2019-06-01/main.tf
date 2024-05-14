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

resource "azapi_resource" "softwareUpdateConfiguration" {
  type      = "Microsoft.Automation/automationAccounts/softwareUpdateConfigurations@2019-06-01"
  parent_id = azapi_resource.automationAccount.id
  name      = var.resource_name
  body = {
    properties = {
      scheduleInfo = {
        description             = ""
        expiryTimeOffsetMinutes = 0
        frequency               = "OneTime"
        interval                = 0
        isEnabled               = true
        nextRunOffsetMinutes    = 0
        startTimeOffsetMinutes  = 0
        timeZone                = "Etc/UTC"
      }
      updateConfiguration = {
        duration = "PT2H"
        linux = {
          excludedPackageNameMasks = [
          ]
          includedPackageClassifications = "Security"
          includedPackageNameMasks = [
          ]
          rebootSetting = "IfRequired"
        }
        operatingSystem = "Linux"
        targets = {
          azureQueries = [
            {
              locations = [
                "westeurope",
              ]
              scope = [
                azapi_resource.resourceGroup.id,
              ]
            },
          ]
        }
      }
    }
  }
  schema_validation_enabled = false
  response_export_values    = ["*"]
}

