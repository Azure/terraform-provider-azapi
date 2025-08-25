---
subcategory: "Microsoft.Web - App Service, Azure Functions"
page_title: "staticSites/customDomains"
description: |-
  Manages a Static Web App Custom Domain.
---

# Microsoft.Web/staticSites/customDomains - Static Web App Custom Domain

This article demonstrates how to use `azapi` provider to manage the Static Web App Custom Domain resource in Azure.

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
  default = "westus"
}

resource "azapi_resource" "resourceGroup" {
  type     = "Microsoft.Resources/resourceGroups@2020-06-01"
  name     = var.resource_name
  location = var.location
}

resource "azapi_resource" "staticSite" {
  type      = "Microsoft.Web/staticSites@2021-02-01"
  parent_id = azapi_resource.resourceGroup.id
  name      = var.resource_name
  location  = var.location
  body = {
    properties = {}
    sku = {
      name = "Free"
    }
  }
}

resource "azapi_resource" "customDomain" {
  type      = "Microsoft.Web/staticSites/customDomains@2021-02-01"
  parent_id = azapi_resource.staticSite.id
  name      = "${var.resource_name}.contoso.com"
  body = {
    properties = {
      validationMethod = "dns-txt-token"
    }
  }
}


```



## Arguments Reference

The following arguments are supported:

* `type` - (Required) The type of the resource. This should be set to `Microsoft.Web/staticSites/customDomains@api-version`. The available api-versions for this resource are: [`2019-08-01`, `2020-06-01`, `2020-09-01`, `2020-10-01`, `2020-12-01`, `2021-01-01`, `2021-01-15`, `2021-02-01`, `2021-03-01`, `2022-03-01`, `2022-09-01`, `2023-01-01`, `2023-12-01`, `2024-04-01`, `2024-11-01`].

* `parent_id` - (Required) The ID of the azure resource in which this resource is created. The allowed values are:  
  `/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Web/staticSites/{resourceName}`

* `name` - (Required) Specifies the name of the azure resource. Changing this forces a new resource to be created.

* `body` - (Required) Specifies the configuration of the resource. More information about the arguments in `body` can be found in the [Microsoft documentation](https://learn.microsoft.com/en-us/azure/templates/Microsoft.Web/staticSites/customDomains?pivots=deployment-language-terraform).

For other arguments, please refer to the [azapi_resource](https://registry.terraform.io/providers/Azure/azapi/latest/docs/resources/resource) documentation.

## Import

 ```shell
 # Azure resource can be imported using the resource id, e.g.
 terraform import azapi_resource.example /subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Web/staticSites/{resourceName}/customDomains/{resourceName}
 
 # It also supports specifying API version by using the resource id with api-version as a query parameter, e.g.
 terraform import azapi_resource.example /subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Web/staticSites/{resourceName}/customDomains/{resourceName}?api-version=2024-11-01
 ```
