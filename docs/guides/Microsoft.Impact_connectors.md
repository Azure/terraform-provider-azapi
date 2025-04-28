---
subcategory: "Microsoft.Impact - Impact"
page_title: "connectors"
description: |-
  Manages a Impact Connectors.
---

# Microsoft.Impact/connectors - Impact Connectors

This article demonstrates how to use `azapi` provider to manage the Impact Connectors resource in Azure.

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

data "azapi_client_config" "current" {}

resource "azapi_resource" "connector" {
  type      = "Microsoft.Impact/connectors@2024-05-01-preview"
  parent_id = "/subscriptions/${data.azapi_client_config.current.subscription_id}"
  name      = var.resource_name
  body = {
    properties = {
      connectorType = "AzureMonitor"
    }
  }
  schema_validation_enabled = false
}

```



## Arguments Reference

The following arguments are supported:

* `type` - (Required) The type of the resource. This should be set to `Microsoft.Impact/connectors@api-version`. The available api-versions for this resource are: [`2024-05-01-preview`].

* `parent_id` - (Required) The ID of the azure resource in which this resource is created. The allowed values are:  
  `/subscriptions/{subscriptionId}`

* `name` - (Required) Specifies the name of the azure resource. Changing this forces a new resource to be created.

* `body` - (Required) Specifies the configuration of the resource. More information about the arguments in `body` can be found in the [Microsoft documentation](https://learn.microsoft.com/en-us/azure/templates/Microsoft.Impact/connectors?pivots=deployment-language-terraform).

## Import

 ```shell
 # Azure resource can be imported using the resource id, e.g.
 terraform import azapi_resource.example /subscriptions/{subscriptionId}/providers/Microsoft.Impact/connectors/{resourceName}
 
 # It also supports specifying API version by using the resource id with api-version as a query parameter, e.g.
 terraform import azapi_resource.example /subscriptions/{subscriptionId}/providers/Microsoft.Impact/connectors/{resourceName}?api-version=2024-05-01-preview
 ```
