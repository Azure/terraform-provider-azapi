---
subcategory: "Microsoft.AzureActiveDirectory - Microsoft Entra ID B2C"
page_title: "b2cDirectories"
description: |-
  Manages a AAD B2C Directory.
---

# Microsoft.AzureActiveDirectory/b2cDirectories - AAD B2C Directory

This article demonstrates how to use `azapi` provider to manage the AAD B2C Directory resource in Azure.

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
  default = "acctest0003"
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

data "azapi_resource_id" "b2cDirectory" {
  type      = "Microsoft.AzureActiveDirectory/b2cDirectories@2021-04-01-preview"
  parent_id = azapi_resource.resourceGroup.id
  name      = "${var.resource_name}.onmicrosoft.com"
}

resource "azapi_resource_action" "b2cDirectory" {
  type        = "Microsoft.AzureActiveDirectory/b2cDirectories@2021-04-01-preview"
  resource_id = data.azapi_resource_id.b2cDirectory.id
  method      = "PUT"
  body = {
    location = "United States"
    properties = {
      createTenantProperties = {
        countryCode = "US"
        displayName = var.resource_name
      }
    }
    sku = {
      name = "PremiumP1"
      tier = "A0"
    }

  }
}

```



## Arguments Reference

The following arguments are supported:

* `type` - (Required) The type of the resource. This should be set to `Microsoft.AzureActiveDirectory/b2cDirectories@api-version`. The available api-versions for this resource are: [`2019-01-01-preview`, `2021-04-01`, `2023-01-18-preview`, `2023-05-17-preview`].

* `parent_id` - (Required) The ID of the azure resource in which this resource is created. The allowed values are:  
  `/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}`

* `name` - (Required) Specifies the name of the azure resource. Changing this forces a new resource to be created.

* `body` - (Required) Specifies the configuration of the resource. More information about the arguments in `body` can be found in the [Microsoft documentation](https://learn.microsoft.com/en-us/azure/templates/Microsoft.AzureActiveDirectory/b2cDirectories?pivots=deployment-language-terraform).

## Import

 ```shell
 # Azure resource can be imported using the resource id, e.g.
 terraform import azapi_resource.example /subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.AzureActiveDirectory/b2cDirectories/{resourceName}
 
 # It also supports specifying API version by using the resource id with api-version as a query parameter, e.g.
 terraform import azapi_resource.example /subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.AzureActiveDirectory/b2cDirectories/{resourceName}?api-version=2023-05-17-preview
 ```
