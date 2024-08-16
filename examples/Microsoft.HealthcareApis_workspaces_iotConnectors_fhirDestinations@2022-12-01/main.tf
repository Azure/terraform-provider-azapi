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

resource "azapi_resource" "namespace" {
  type      = "Microsoft.EventHub/namespaces@2022-01-01-preview"
  parent_id = azapi_resource.resourceGroup.id
  name      = var.resource_name
  location  = var.location
  body = {
    properties = {
      disableLocalAuth     = false
      isAutoInflateEnabled = false
      publicNetworkAccess  = "Enabled"
      zoneRedundant        = false
    }
    sku = {
      capacity = 1
      name     = "Standard"
      tier     = "Standard"
    }
  }
  schema_validation_enabled = false
  response_export_values    = ["*"]
}

resource "azapi_resource" "eventhub" {
  type      = "Microsoft.EventHub/namespaces/eventhubs@2021-11-01"
  parent_id = azapi_resource.namespace.id
  name      = var.resource_name
  body = {
    properties = {
      messageRetentionInDays = 1
      partitionCount         = 2
      status                 = "Active"
    }
  }
  schema_validation_enabled = false
  response_export_values    = ["*"]
}

resource "azapi_resource" "consumerGroup" {
  type      = "Microsoft.EventHub/namespaces/eventhubs/consumerGroups@2021-11-01"
  parent_id = azapi_resource.eventhub.id
  name      = var.resource_name
  body = {
    properties = {
      userMetadata = ""
    }
  }
  schema_validation_enabled = false
  response_export_values    = ["*"]
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

resource "azapi_resource" "iotConnector" {
  type      = "Microsoft.HealthcareApis/workspaces/iotConnectors@2022-12-01"
  parent_id = azapi_resource.workspace.id
  name      = var.resource_name
  location  = var.location
  body = {
    properties = {
      deviceMapping = {
        content = {
          template = [
          ]
          templateType = "CollectionContent"
        }
      }
      ingestionEndpointConfiguration = {
        consumerGroup                   = azapi_resource.consumerGroup.id
        eventHubName                    = azapi_resource.eventhub.name
        fullyQualifiedEventHubNamespace = "${azapi_resource.namespace.name}.servicebus.windows.net"
      }
    }
  }
  schema_validation_enabled = false
  response_export_values    = ["*"]
}

resource "azapi_resource" "fhirDestination" {
  type      = "Microsoft.HealthcareApis/workspaces/iotConnectors/fhirDestinations@2022-12-01"
  parent_id = azapi_resource.iotConnector.id
  name      = var.resource_name
  location  = var.location
  body = {
    properties = {
      fhirMapping = {
        content = {
          template = [
          ]
          templateType = "CollectionFhirTemplate"
        }
      }
      fhirServiceResourceId          = azapi_resource.fhirService.id
      resourceIdentityResolutionType = "Create"
    }
  }
  schema_validation_enabled = false
  response_export_values    = ["*"]
}

