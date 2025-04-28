---
subcategory: "Microsoft.Security - Security Center"
page_title: "securityContacts"
description: |-
  Manages a subscription's Security Center Contact.
---

# Microsoft.Security/securityContacts - subscription's Security Center Contact

This article demonstrates how to use `azapi` provider to manage the subscription's Security Center Contact resource in Azure.

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
  default = "eastus"
}

data "azurerm_client_config" "current" {
}

resource "azapi_resource" "securityContact" {
  type      = "Microsoft.Security/securityContacts@2017-08-01-preview"
  parent_id = "/subscriptions/${data.azurerm_client_config.current.subscription_id}"
  name      = var.resource_name
  body = {
    properties = {
      alertNotifications = "On"
      alertsToAdmins     = "On"
      email              = "basic@example.com"
      phone              = "+1-555-555-5555"
    }
  }
  schema_validation_enabled = false
  response_export_values    = ["*"]
}


```



## Arguments Reference

The following arguments are supported:

* `type` - (Required) The type of the resource. This should be set to `Microsoft.Security/securityContacts@api-version`. The available api-versions for this resource are: [`2017-08-01-preview`, `2020-01-01-preview`, `2023-12-01-preview`].

* `parent_id` - (Required) The ID of the azure resource in which this resource is created. The allowed values are:  
  `/subscriptions/{subscriptionId}`

* `name` - (Required) Specifies the name of the azure resource. Changing this forces a new resource to be created.

* `body` - (Required) Specifies the configuration of the resource. More information about the arguments in `body` can be found in the [Microsoft documentation](https://learn.microsoft.com/en-us/azure/templates/Microsoft.Security/securityContacts?pivots=deployment-language-terraform).

## Import

 ```shell
 # Azure resource can be imported using the resource id, e.g.
 terraform import azapi_resource.example /subscriptions/{subscriptionId}/providers/Microsoft.Security/securityContacts/{resourceName}
 
 # It also supports specifying API version by using the resource id with api-version as a query parameter, e.g.
 terraform import azapi_resource.example /subscriptions/{subscriptionId}/providers/Microsoft.Security/securityContacts/{resourceName}?api-version=2023-12-01-preview
 ```
