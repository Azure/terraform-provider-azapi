---
subcategory: "Microsoft.HealthcareApis - Healthcare APIs"
page_title: "workspaces/iotConnectors/fhirDestinations"
description: |-
  Manages a Healthcare Med Tech (Internet of Medical Things) Service Fhir Destination.
---

# Microsoft.HealthcareApis/workspaces/iotConnectors/fhirDestinations - Healthcare Med Tech (Internet of Medical Things) Service Fhir Destination

This article demonstrates how to use `azapi` provider to manage the Healthcare Med Tech (Internet of Medical Things) Service Fhir Destination resource in Azure.

## Example Usage

### default

```hcl
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


```



## Arguments Reference

The following arguments are supported:

* `type` - (Required) The type of the resource. This should be set to `Microsoft.HealthcareApis/workspaces/iotConnectors/fhirDestinations@api-version`. The available api-versions for this resource are: [`2021-06-01-preview`, `2021-11-01`, `2022-01-31-preview`, `2022-05-15`, `2022-06-01`, `2022-10-01-preview`, `2022-12-01`, `2023-02-28`, `2023-09-06`, `2023-11-01`, `2023-12-01`, `2024-03-01`, `2024-03-31`, `2025-03-01-preview`].

* `parent_id` - (Required) The ID of the azure resource in which this resource is created. The allowed values are:  
  `/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.HealthcareApis/workspaces/{resourceName}/iotConnectors/{resourceName}`

* `name` - (Required) Specifies the name of the azure resource. Changing this forces a new resource to be created.

* `body` - (Required) Specifies the configuration of the resource. More information about the arguments in `body` can be found in the [Microsoft documentation](https://learn.microsoft.com/en-us/azure/templates/Microsoft.HealthcareApis/workspaces/iotConnectors/fhirDestinations?pivots=deployment-language-terraform).

For other arguments, please refer to the [azapi_resource](https://registry.terraform.io/providers/Azure/azapi/latest/docs/resources/resource) documentation.

## Import

 ```shell
 # Azure resource can be imported using the resource id, e.g.
 terraform import azapi_resource.example /subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.HealthcareApis/workspaces/{resourceName}/iotConnectors/{resourceName}/fhirDestinations/{resourceName}
 
 # It also supports specifying API version by using the resource id with api-version as a query parameter, e.g.
 terraform import azapi_resource.example /subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.HealthcareApis/workspaces/{resourceName}/iotConnectors/{resourceName}/fhirDestinations/{resourceName}?api-version=2025-03-01-preview
 ```
