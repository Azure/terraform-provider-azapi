---
subcategory: "Microsoft.EventHub - Event Hubs"
page_title: "namespaces/disasterRecoveryConfigs"
description: |-
  Manages a Disaster Recovery Config for an Event Hub Namespace.
---

# Microsoft.EventHub/namespaces/disasterRecoveryConfigs - Disaster Recovery Config for an Event Hub Namespace

This article demonstrates how to use `azapi` provider to manage the Disaster Recovery Config for an Event Hub Namespace resource in Azure.



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

resource "azapi_resource" "namespace2" {
  type      = "Microsoft.EventHub/namespaces@2022-01-01-preview"
  parent_id = azapi_resource.resourceGroup.id
  name      = var.resource_name
  location  = "westus2"
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

resource "azapi_resource" "disasterRecoveryConfig" {
  type      = "Microsoft.EventHub/namespaces/disasterRecoveryConfigs@2021-11-01"
  parent_id = azapi_resource.namespace.id
  name      = var.resource_name
  body = {
    properties = {
      partnerNamespace = azapi_resource.namespace2.id
    }
  }
  schema_validation_enabled = false
  response_export_values    = ["*"]
}


```



## Arguments Reference

The following arguments are supported:

* `type` - (Required) The type of the resource. This should be set to `Microsoft.EventHub/namespaces/disasterRecoveryConfigs@api-version`. The available api-versions for this resource are: [`2017-04-01`, `2018-01-01-preview`, `2021-01-01-preview`, `2021-06-01-preview`, `2021-11-01`, `2022-01-01-preview`, `2022-10-01-preview`, `2023-01-01-preview`, `2024-01-01`, `2024-05-01-preview`, `2025-05-01-preview`].

* `parent_id` - (Required) The ID of the azure resource in which this resource is created. The allowed values are:  
  `/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.EventHub/namespaces/{resourceName}`

* `name` - (Required) Specifies the name of the azure resource. Changing this forces a new resource to be created.

* `body` - (Required) Specifies the configuration of the resource. More information about the arguments in `body` can be found in the [Microsoft documentation](https://learn.microsoft.com/en-us/azure/templates/Microsoft.EventHub/namespaces/disasterRecoveryConfigs?pivots=deployment-language-terraform).

For other arguments, please refer to the [azapi_resource](https://registry.terraform.io/providers/Azure/azapi/latest/docs/resources/resource) documentation.

## Import

 ```shell
 # Azure resource can be imported using the resource id, e.g.
 terraform import azapi_resource.example /subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.EventHub/namespaces/{resourceName}/disasterRecoveryConfigs/{resourceName}
 
 # It also supports specifying API version by using the resource id with api-version as a query parameter, e.g.
 terraform import azapi_resource.example /subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.EventHub/namespaces/{resourceName}/disasterRecoveryConfigs/{resourceName}?api-version=2025-05-01-preview
 ```
