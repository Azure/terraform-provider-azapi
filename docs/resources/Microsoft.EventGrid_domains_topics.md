---
subcategory: "Microsoft.EventGrid - Event Grid"
page_title: "domains/topics"
description: |-
  Manages a EventGrid Domain Topic.
---

# Microsoft.EventGrid/domains/topics - EventGrid Domain Topic

This article demonstrates how to use `azapi` provider to manage the EventGrid Domain Topic resource in Azure.



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

resource "azapi_resource" "domain" {
  type      = "Microsoft.EventGrid/domains@2021-12-01"
  parent_id = azapi_resource.resourceGroup.id
  name      = var.resource_name
  location  = var.location
  body = {
    properties = {
      autoCreateTopicWithFirstSubscription = true
      autoDeleteTopicWithLastSubscription  = true
      disableLocalAuth                     = false
      inputSchema                          = "EventGridSchema"
      inputSchemaMapping                   = null
      publicNetworkAccess                  = "Enabled"
    }
  }
  schema_validation_enabled = false
  response_export_values    = ["*"]
}

resource "azapi_resource" "topic" {
  type                      = "Microsoft.EventGrid/domains/topics@2021-12-01"
  parent_id                 = azapi_resource.domain.id
  name                      = var.resource_name
  schema_validation_enabled = false
  response_export_values    = ["*"]
}


```



## Arguments Reference

The following arguments are supported:

* `type` - (Required) The type of the resource. This should be set to `Microsoft.EventGrid/domains/topics@api-version`. The available api-versions for this resource are: [`2018-09-15-preview`, `2019-02-01-preview`, `2019-06-01`, `2020-01-01-preview`, `2020-04-01-preview`, `2020-06-01`, `2020-10-15-preview`, `2021-06-01-preview`, `2021-10-15-preview`, `2021-12-01`, `2022-06-15`, `2023-06-01-preview`, `2023-12-15-preview`, `2024-06-01-preview`, `2024-12-15-preview`, `2025-02-15`, `2025-04-01-preview`].

* `parent_id` - (Required) The ID of the azure resource in which this resource is created. The allowed values are:  
  `/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.EventGrid/domains/{resourceName}`

* `name` - (Required) Specifies the name of the azure resource. Changing this forces a new resource to be created.

* `body` - (Required) Specifies the configuration of the resource. More information about the arguments in `body` can be found in the [Microsoft documentation](https://learn.microsoft.com/en-us/azure/templates/Microsoft.EventGrid/domains/topics?pivots=deployment-language-terraform).

For other arguments, please refer to the [azapi_resource](https://registry.terraform.io/providers/Azure/azapi/latest/docs/resources/resource) documentation.

## Import

 ```shell
 # Azure resource can be imported using the resource id, e.g.
 terraform import azapi_resource.example /subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.EventGrid/domains/{resourceName}/topics/{resourceName}
 
 # It also supports specifying API version by using the resource id with api-version as a query parameter, e.g.
 terraform import azapi_resource.example /subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.EventGrid/domains/{resourceName}/topics/{resourceName}?api-version=2025-04-01-preview
 ```
