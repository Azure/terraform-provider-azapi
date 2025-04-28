---
subcategory: "Microsoft.DigitalTwins - Azure Digital Twins"
page_title: "digitalTwinsInstances/endpoints"
description: |-
  Manages a Digital Twins Endpoint.
---

# Microsoft.DigitalTwins/digitalTwinsInstances/endpoints - Digital Twins Endpoint

This article demonstrates how to use `azapi` provider to manage the Digital Twins Endpoint resource in Azure.

## Example Usage

### default

```hcl
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

resource "azapi_resource" "namespace" {
  type      = "Microsoft.ServiceBus/namespaces@2022-01-01-preview"
  parent_id = azapi_resource.resourceGroup.id
  name      = var.resource_name
  location  = var.location
  body = {
    properties = {
      disableLocalAuth    = false
      publicNetworkAccess = "Enabled"
      zoneRedundant       = false
    }
    sku = {
      capacity = 0
      name     = "Standard"
      tier     = "Standard"
    }
  }
  schema_validation_enabled = false
  response_export_values    = ["*"]
}

resource "azapi_resource" "digitalTwinsInstance" {
  type      = "Microsoft.DigitalTwins/digitalTwinsInstances@2020-12-01"
  parent_id = azapi_resource.resourceGroup.id
  name      = var.resource_name
  location  = var.location
  body = {
  }
  schema_validation_enabled = false
  response_export_values    = ["*"]
}

resource "azapi_resource" "topic" {
  type      = "Microsoft.ServiceBus/namespaces/topics@2021-06-01-preview"
  parent_id = azapi_resource.namespace.id
  name      = var.resource_name
  body = {
    properties = {
      enableBatchedOperations    = false
      enableExpress              = false
      enablePartitioning         = false
      maxSizeInMegabytes         = 5120
      requiresDuplicateDetection = false
      status                     = "Active"
      supportOrdering            = false
    }
  }
  schema_validation_enabled = false
  response_export_values    = ["*"]
}

resource "azapi_resource" "authorizationRule" {
  type      = "Microsoft.ServiceBus/namespaces/topics/authorizationRules@2021-06-01-preview"
  parent_id = azapi_resource.topic.id
  name      = var.resource_name
  body = {
    properties = {
      rights = [
        "Send",
      ]
    }
  }
  schema_validation_enabled = false
  response_export_values    = ["*"]
}

data "azapi_resource_action" "listKeys" {
  type                   = "Microsoft.ServiceBus/namespaces/topics/authorizationRules@2021-06-01-preview"
  resource_id            = azapi_resource.authorizationRule.id
  action                 = "listKeys"
  response_export_values = ["*"]
}

resource "azapi_resource" "endpoint" {
  type      = "Microsoft.DigitalTwins/digitalTwinsInstances/endpoints@2020-12-01"
  parent_id = azapi_resource.digitalTwinsInstance.id
  name      = var.resource_name
  body = {
    properties = {
      authenticationType        = "KeyBased"
      deadLetterSecret          = ""
      endpointType              = "ServiceBus"
      primaryConnectionString   = data.azapi_resource_action.listKeys.output.primaryConnectionString
      secondaryConnectionString = data.azapi_resource_action.listKeys.output.secondaryConnectionString
    }
  }
  schema_validation_enabled = false
  response_export_values    = ["*"]
  lifecycle {
    ignore_changes = [body.properties.primaryConnectionString, body.properties.secondaryConnectionString]
  }
}


```



## Arguments Reference

The following arguments are supported:

* `type` - (Required) The type of the resource. This should be set to `Microsoft.DigitalTwins/digitalTwinsInstances/endpoints@api-version`. The available api-versions for this resource are: [`2020-03-01-preview`, `2020-10-31`, `2020-12-01`, `2021-06-30-preview`, `2022-05-31`, `2022-10-31`, `2023-01-31`].

* `parent_id` - (Required) The ID of the azure resource in which this resource is created. The allowed values are:  
  `/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.DigitalTwins/digitalTwinsInstances/{resourceName}`

* `name` - (Required) Specifies the name of the azure resource. Changing this forces a new resource to be created.

* `body` - (Required) Specifies the configuration of the resource. More information about the arguments in `body` can be found in the [Microsoft documentation](https://learn.microsoft.com/en-us/azure/templates/Microsoft.DigitalTwins/digitalTwinsInstances/endpoints?pivots=deployment-language-terraform).

## Import

 ```shell
 # Azure resource can be imported using the resource id, e.g.
 terraform import azapi_resource.example /subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.DigitalTwins/digitalTwinsInstances/{resourceName}/endpoints/{resourceName}
 
 # It also supports specifying API version by using the resource id with api-version as a query parameter, e.g.
 terraform import azapi_resource.example /subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.DigitalTwins/digitalTwinsInstances/{resourceName}/endpoints/{resourceName}?api-version=2023-01-31
 ```
