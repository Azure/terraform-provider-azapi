---
subcategory: "Microsoft.Relay - Azure Relay"
page_title: "namespaces/hybridConnections/authorizationRules"
description: |-
  Manages a Azure Relay Hybrid Connection Authorization Rule.
---

# Microsoft.Relay/namespaces/hybridConnections/authorizationRules - Azure Relay Hybrid Connection Authorization Rule

This article demonstrates how to use `azapi` provider to manage the Azure Relay Hybrid Connection Authorization Rule resource in Azure.

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
  type      = "Microsoft.Relay/namespaces@2017-04-01"
  parent_id = azapi_resource.resourceGroup.id
  name      = var.resource_name
  location  = var.location
  body = {
    properties = {
    }
    sku = {
      name = "Standard"
      tier = "Standard"
    }
  }
  schema_validation_enabled = false
  response_export_values    = ["*"]
}

resource "azapi_resource" "hybridConnection" {
  type      = "Microsoft.Relay/namespaces/hybridConnections@2017-04-01"
  parent_id = azapi_resource.namespace.id
  name      = var.resource_name
  body = {
    properties = {
      requiresClientAuthorization = true
      userMetadata                = ""
    }
  }
  schema_validation_enabled = false
  response_export_values    = ["*"]
}

resource "azapi_resource" "authorizationRule" {
  type      = "Microsoft.Relay/namespaces/hybridConnections/authorizationRules@2017-04-01"
  parent_id = azapi_resource.hybridConnection.id
  name      = var.resource_name
  body = {
    properties = {
      rights = [
        "Listen",
        "Send",
      ]
    }
  }
  schema_validation_enabled = false
  response_export_values    = ["*"]
}


```



## Arguments Reference

The following arguments are supported:

* `type` - (Required) The type of the resource. This should be set to `Microsoft.Relay/namespaces/hybridConnections/authorizationRules@api-version`. The available api-versions for this resource are: [`2016-07-01`, `2017-04-01`, `2021-11-01`, `2024-01-01`].

* `parent_id` - (Required) The ID of the azure resource in which this resource is created. The allowed values are:  
  `/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Relay/namespaces/{resourceName}/hybridConnections/{resourceName}`

* `name` - (Required) Specifies the name of the azure resource. Changing this forces a new resource to be created.

* `body` - (Required) Specifies the configuration of the resource. More information about the arguments in `body` can be found in the [Microsoft documentation](https://learn.microsoft.com/en-us/azure/templates/Microsoft.Relay/namespaces/hybridConnections/authorizationRules?pivots=deployment-language-terraform).

## Import

 ```shell
 # Azure resource can be imported using the resource id, e.g.
 terraform import azapi_resource.example /subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Relay/namespaces/{resourceName}/hybridConnections/{resourceName}/authorizationRules/{resourceName}
 
 # It also supports specifying API version by using the resource id with api-version as a query parameter, e.g.
 terraform import azapi_resource.example /subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Relay/namespaces/{resourceName}/hybridConnections/{resourceName}/authorizationRules/{resourceName}?api-version=2024-01-01
 ```
