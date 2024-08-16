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
  type                      = "Microsoft.HealthcareApis/workspaces@2022-12-01"
  parent_id                 = azapi_resource.resourceGroup.id
  name                      = var.resource_name
  location                  = var.location
  schema_validation_enabled = false
  response_export_values    = ["*"]
}

resource "azapi_resource" "fhirService" {
  type      = "Microsoft.HealthcareApis/workspaces/fhirServices@2022-12-01"
  parent_id = azapi_resource.workspace.id
  name      = var.resource_name
  location  = var.location
  body = {
    kind = "fhir-R4"
    properties = {
      acrConfiguration = {
      }
      authenticationConfiguration = {
        audience          = "https://acctestfhir.fhir.azurehealthcareapis.com"
        authority         = "https://login.microsoftonline.com/${data.azurerm_client_config.current.tenant_id}"
        smartProxyEnabled = false
      }
      corsConfiguration = {
        allowCredentials = false
        headers = [
        ]
        methods = [
        ]
        origins = [
        ]
      }
    }
  }
  schema_validation_enabled = false
  response_export_values    = ["*"]
}

resource "azapi_resource" "fhirService2" {
  type      = "Microsoft.HealthcareApis/workspaces/fhirServices@2022-12-01"
  parent_id = azapi_resource.workspace.id
  name      = var.resource_name
  location  = var.location
  body = {
    kind = "fhir-R4"
    properties = {
      acrConfiguration = {
      }
      authenticationConfiguration = {
        audience          = azapi_resource.fhirService.output.properties.authenticationConfiguration.audience
        authority         = azapi_resource.fhirService.output.properties.authenticationConfiguration.authority
        smartProxyEnabled = false
      }
      corsConfiguration = {
        allowCredentials = false
        headers = [
        ]
        methods = [
        ]
        origins = [
        ]
      }
    }
  }
  schema_validation_enabled = false
  response_export_values    = ["*"]
}

