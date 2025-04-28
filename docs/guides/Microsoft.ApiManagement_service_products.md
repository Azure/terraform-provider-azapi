---
subcategory: "Microsoft.ApiManagement - API Management"
page_title: "service/products"
description: |-
  Manages a API Management Product.
---

# Microsoft.ApiManagement/service/products - API Management Product

This article demonstrates how to use `azapi` provider to manage the API Management Product resource in Azure.

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

resource "azapi_resource" "service" {
  type      = "Microsoft.ApiManagement/service@2021-08-01"
  parent_id = azapi_resource.resourceGroup.id
  name      = var.resource_name
  location  = var.location
  body = {
    properties = {
      certificates = [
      ]
      customProperties = {
        "Microsoft.WindowsAzure.ApiManagement.Gateway.Security.Backend.Protocols.Ssl30" = "false"
        "Microsoft.WindowsAzure.ApiManagement.Gateway.Security.Backend.Protocols.Tls10" = "false"
        "Microsoft.WindowsAzure.ApiManagement.Gateway.Security.Backend.Protocols.Tls11" = "false"
        "Microsoft.WindowsAzure.ApiManagement.Gateway.Security.Protocols.Tls10"         = "false"
        "Microsoft.WindowsAzure.ApiManagement.Gateway.Security.Protocols.Tls11"         = "false"
      }
      disableGateway      = false
      publicNetworkAccess = "Enabled"
      publisherEmail      = "pub1@email.com"
      publisherName       = "pub1"
      virtualNetworkType  = "None"
    }
    sku = {
      capacity = 0
      name     = "Consumption"
    }
  }
  schema_validation_enabled = false
  response_export_values    = ["*"]
}

resource "azapi_resource" "product" {
  type      = "Microsoft.ApiManagement/service/products@2021-08-01"
  parent_id = azapi_resource.service.id
  name      = var.resource_name
  body = {
    properties = {
      description          = ""
      displayName          = "Test Product"
      state                = "notPublished"
      subscriptionRequired = false
      terms                = ""
    }
  }
  schema_validation_enabled = false
  response_export_values    = ["*"]
}


```



## Arguments Reference

The following arguments are supported:

* `type` - (Required) The type of the resource. This should be set to `Microsoft.ApiManagement/service/products@api-version`. The available api-versions for this resource are: [`2016-07-07`, `2016-10-10`, `2017-03-01`, `2018-01-01`, `2018-06-01-preview`, `2019-01-01`, `2019-12-01`, `2019-12-01-preview`, `2020-06-01-preview`, `2020-12-01`, `2021-01-01-preview`, `2021-04-01-preview`, `2021-08-01`, `2021-12-01-preview`, `2022-04-01-preview`, `2022-08-01`, `2022-09-01-preview`, `2023-03-01-preview`, `2023-05-01-preview`, `2023-09-01-preview`, `2024-05-01`, `2024-06-01-preview`].

* `parent_id` - (Required) The ID of the azure resource in which this resource is created. The allowed values are:  
  `/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.ApiManagement/service/{resourceName}`

* `name` - (Required) Specifies the name of the azure resource. Changing this forces a new resource to be created.

* `body` - (Required) Specifies the configuration of the resource. More information about the arguments in `body` can be found in the [Microsoft documentation](https://learn.microsoft.com/en-us/azure/templates/Microsoft.ApiManagement/service/products?pivots=deployment-language-terraform).

## Import

 ```shell
 # Azure resource can be imported using the resource id, e.g.
 terraform import azapi_resource.example /subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.ApiManagement/service/{resourceName}/products/{resourceName}
 
 # It also supports specifying API version by using the resource id with api-version as a query parameter, e.g.
 terraform import azapi_resource.example /subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.ApiManagement/service/{resourceName}/products/{resourceName}?api-version=2024-06-01-preview
 ```
