---
subcategory: "Microsoft.Insights - Azure Monitor"
page_title: "privateLinkScopes/scopedResources"
description: |-
  Manages a Azure Monitor Private Link Scoped Service.
---

# Microsoft.Insights/privateLinkScopes/scopedResources - Azure Monitor Private Link Scoped Service

This article demonstrates how to use `azapi` provider to manage the Azure Monitor Private Link Scoped Service resource in Azure.



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

resource "azapi_resource" "component" {
  type      = "Microsoft.Insights/components@2020-02-02"
  parent_id = azapi_resource.resourceGroup.id
  name      = var.resource_name
  location  = var.location
  body = {
    kind = "web"
    properties = {
      Application_Type                = "web"
      DisableIpMasking                = false
      DisableLocalAuth                = false
      ForceCustomerStorageForProfiler = false
      RetentionInDays                 = 90
      SamplingPercentage              = 100
      publicNetworkAccessForIngestion = "Enabled"
      publicNetworkAccessForQuery     = "Enabled"
    }
  }
  schema_validation_enabled = false
  response_export_values    = ["*"]
}

resource "azapi_resource" "privateLinkScope" {
  type      = "Microsoft.Insights/privateLinkScopes@2019-10-17-preview"
  parent_id = azapi_resource.resourceGroup.id
  name      = var.resource_name
  location  = "Global"
  body = {
    properties = {
    }
  }
  schema_validation_enabled = false
  response_export_values    = ["*"]
}

resource "azapi_resource" "scopedResource" {
  type      = "Microsoft.Insights/privateLinkScopes/scopedResources@2019-10-17-preview"
  parent_id = azapi_resource.privateLinkScope.id
  name      = var.resource_name
  body = {
    properties = {
      linkedResourceId = azapi_resource.component.id
    }
  }
  schema_validation_enabled = false
  ignore_casing             = true
  response_export_values    = ["*"]
}


```



## Arguments Reference

The following arguments are supported:

* `type` - (Required) The type of the resource. This should be set to `Microsoft.Insights/privateLinkScopes/scopedResources@api-version`. The available api-versions for this resource are: [`2019-10-17-preview`, `2021-07-01-preview`, `2021-09-01`, `2023-06-01-preview`].

* `parent_id` - (Required) The ID of the azure resource in which this resource is created. The allowed values are:  
  `/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Insights/privateLinkScopes/{resourceName}`

* `name` - (Required) Specifies the name of the azure resource. Changing this forces a new resource to be created.

* `body` - (Required) Specifies the configuration of the resource. More information about the arguments in `body` can be found in the [Microsoft documentation](https://learn.microsoft.com/en-us/azure/templates/Microsoft.Insights/privateLinkScopes/scopedResources?pivots=deployment-language-terraform).

For other arguments, please refer to the [azapi_resource](https://registry.terraform.io/providers/Azure/azapi/latest/docs/resources/resource) documentation.

## Import

 ```shell
 # Azure resource can be imported using the resource id, e.g.
 terraform import azapi_resource.example /subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Insights/privateLinkScopes/{resourceName}/scopedResources/{resourceName}
 
 # It also supports specifying API version by using the resource id with api-version as a query parameter, e.g.
 terraform import azapi_resource.example /subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Insights/privateLinkScopes/{resourceName}/scopedResources/{resourceName}?api-version=2023-06-01-preview
 ```
