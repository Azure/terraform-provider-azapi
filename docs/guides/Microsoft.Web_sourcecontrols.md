---
subcategory: "Microsoft.Web - App Service, Azure Functions"
page_title: "sourcecontrols"
description: |-
  Manages a App Service GitHub Token.
---

# Microsoft.Web/sourcecontrols - App Service GitHub Token

This article demonstrates how to use `azapi` provider to manage the App Service GitHub Token resource in Azure.

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
  default = "eastus"
}

resource "azapi_update_resource" "sourcecontrol" {
  type      = "Microsoft.Web/sourcecontrols@2021-02-01"
  parent_id = "/"
  name      = "GitHub"
  body = {
    properties = {
      token       = "abcdefghijklmnopqrstuvwxyz"
      tokenSecret = "abcdefghijklmnopqrstuvwxyz"
    }
  }
  response_export_values = ["*"]
}


```



## Arguments Reference

The following arguments are supported:

* `type` - (Required) The type of the resource. This should be set to `Microsoft.Web/sourcecontrols@api-version`. The available api-versions for this resource are: [`2015-08-01`, `2016-03-01`, `2018-02-01`, `2019-08-01`, `2020-06-01`, `2020-09-01`, `2020-10-01`, `2020-12-01`, `2021-01-01`, `2021-01-15`, `2021-02-01`, `2021-03-01`, `2022-03-01`, `2022-09-01`, `2023-01-01`, `2023-12-01`, `2024-04-01`].

* `parent_id` - (Required) The ID of the azure resource in which this resource is created. The allowed values are:  
  `/`

* `name` - (Required) Specifies the name of the azure resource. Changing this forces a new resource to be created.

* `body` - (Required) Specifies the configuration of the resource. More information about the arguments in `body` can be found in the [Microsoft documentation](https://learn.microsoft.com/en-us/azure/templates/Microsoft.Web/sourcecontrols?pivots=deployment-language-terraform).

## Import

 ```shell
 # Azure resource can be imported using the resource id, e.g.
 terraform import azapi_resource.example //providers/Microsoft.Web/sourcecontrols/{resourceName}
 
 # It also supports specifying API version by using the resource id with api-version as a query parameter, e.g.
 terraform import azapi_resource.example //providers/Microsoft.Web/sourcecontrols/{resourceName}?api-version=2024-04-01
 ```
