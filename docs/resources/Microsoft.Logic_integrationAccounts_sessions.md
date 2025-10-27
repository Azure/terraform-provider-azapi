---
subcategory: "Microsoft.Logic - Logic Apps"
page_title: "integrationAccounts/sessions"
description: |-
  Manages a Logic App Integration Account Session.
---

# Microsoft.Logic/integrationAccounts/sessions - Logic App Integration Account Session

This article demonstrates how to use `azapi` provider to manage the Logic App Integration Account Session resource in Azure.



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

resource "azapi_resource" "integrationAccount" {
  type      = "Microsoft.Logic/integrationAccounts@2019-05-01"
  parent_id = azapi_resource.resourceGroup.id
  name      = var.resource_name
  location  = var.location
  body = {
    properties = {
    }
    sku = {
      name = "Basic"
    }
  }
  schema_validation_enabled = false
  response_export_values    = ["*"]
}

resource "azapi_resource" "session" {
  type      = "Microsoft.Logic/integrationAccounts/sessions@2019-05-01"
  parent_id = azapi_resource.integrationAccount.id
  name      = var.resource_name
  body = {
    properties = {
      content = "\t{\n       \"controlNumber\": \"1234\"\n    }\n"
    }
  }
  schema_validation_enabled = false
  response_export_values    = ["*"]
}


```



## Arguments Reference

The following arguments are supported:

* `type` - (Required) The type of the resource. This should be set to `Microsoft.Logic/integrationAccounts/sessions@api-version`. The available api-versions for this resource are: [`2016-06-01`, `2018-07-01-preview`, `2019-05-01`].

* `parent_id` - (Required) The ID of the azure resource in which this resource is created. The allowed values are:  
  `/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Logic/integrationAccounts/{resourceName}`

* `name` - (Required) Specifies the name of the azure resource. Changing this forces a new resource to be created.

* `body` - (Required) Specifies the configuration of the resource. More information about the arguments in `body` can be found in the [Microsoft documentation](https://learn.microsoft.com/en-us/azure/templates/Microsoft.Logic/integrationAccounts/sessions?pivots=deployment-language-terraform).

For other arguments, please refer to the [azapi_resource](https://registry.terraform.io/providers/Azure/azapi/latest/docs/resources/resource) documentation.

## Import

 ```shell
 # Azure resource can be imported using the resource id, e.g.
 terraform import azapi_resource.example /subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Logic/integrationAccounts/{resourceName}/sessions/{resourceName}
 
 # It also supports specifying API version by using the resource id with api-version as a query parameter, e.g.
 terraform import azapi_resource.example /subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Logic/integrationAccounts/{resourceName}/sessions/{resourceName}?api-version=2019-05-01
 ```
