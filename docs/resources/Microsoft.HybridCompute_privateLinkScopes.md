---
subcategory: "Microsoft.HybridCompute - Azure Arc-enabled servers"
page_title: "privateLinkScopes"
description: |-
  Manages a Azure Arc Private Link Scope.
---

# Microsoft.HybridCompute/privateLinkScopes - Azure Arc Private Link Scope

This article demonstrates how to use `azapi` provider to manage the Azure Arc Private Link Scope resource in Azure.



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

resource "azapi_resource" "privateLinkScope" {
  type      = "Microsoft.HybridCompute/privateLinkScopes@2022-11-10"
  parent_id = azapi_resource.resourceGroup.id
  name      = var.resource_name
  location  = var.location
  body = {
    properties = {
      publicNetworkAccess = "Disabled"
    }

  }
  schema_validation_enabled = false
  response_export_values    = ["*"]
}


```



## Arguments Reference

The following arguments are supported:

* `type` - (Required) The type of the resource. This should be set to `Microsoft.HybridCompute/privateLinkScopes@api-version`. The available api-versions for this resource are: [`2020-08-15-preview`, `2021-01-28-preview`, `2021-03-25-preview`, `2021-04-22-preview`, `2021-05-17-preview`, `2021-05-20`, `2021-06-10-preview`, `2021-12-10-preview`, `2022-03-10`, `2022-05-10-preview`, `2022-08-11-preview`, `2022-11-10`, `2022-12-27`, `2022-12-27-preview`, `2023-03-15-preview`, `2023-06-20-preview`, `2023-10-03-preview`, `2024-03-31-preview`, `2024-05-20-preview`, `2024-07-10`, `2024-07-31-preview`, `2024-09-10-preview`, `2024-11-10-preview`, `2025-01-13`, `2025-02-19-preview`].

* `parent_id` - (Required) The ID of the azure resource in which this resource is created. The allowed values are:  
  `/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}`

* `name` - (Required) Specifies the name of the azure resource. Changing this forces a new resource to be created.

* `body` - (Required) Specifies the configuration of the resource. More information about the arguments in `body` can be found in the [Microsoft documentation](https://learn.microsoft.com/en-us/azure/templates/Microsoft.HybridCompute/privateLinkScopes?pivots=deployment-language-terraform).

For other arguments, please refer to the [azapi_resource](https://registry.terraform.io/providers/Azure/azapi/latest/docs/resources/resource) documentation.

## Import

 ```shell
 # Azure resource can be imported using the resource id, e.g.
 terraform import azapi_resource.example /subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.HybridCompute/privateLinkScopes/{resourceName}
 
 # It also supports specifying API version by using the resource id with api-version as a query parameter, e.g.
 terraform import azapi_resource.example /subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.HybridCompute/privateLinkScopes/{resourceName}?api-version=2025-02-19-preview
 ```
