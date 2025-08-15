---
subcategory: "Microsoft.OperationsManagement - Azure Monitor"
page_title: "solutions"
description: |-
  Manages a Log Analytics (formally Operational Insights) Solution.
---

# Microsoft.OperationsManagement/solutions - Log Analytics (formally Operational Insights) Solution

This article demonstrates how to use `azapi` provider to manage the Log Analytics (formally Operational Insights) Solution resource in Azure.

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

data "azapi_client_config" "current" {}

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

resource "azapi_resource" "workspace" {
  type      = "Microsoft.OperationalInsights/workspaces@2022-10-01"
  parent_id = azapi_resource.resourceGroup.id
  name      = var.resource_name
  location  = var.location
  body = {
    properties = {
      sku = {
        name = "PerGB2018"
      }
    }
  }
}

resource "azapi_resource" "solution" {
  type      = "Microsoft.OperationsManagement/solutions@2015-11-01-preview"
  parent_id = azapi_resource.resourceGroup.id
  name      = "ContainerInsights(${var.resource_name})"
  location  = var.location
  body = {
    plan = {
      name          = "ContainerInsights(${var.resource_name})"
      product       = "OMSGallery/ContainerInsights"
      promotionCode = ""
      publisher     = "Microsoft"
    }
    properties = {
      workspaceResourceId = azapi_resource.workspace.id
    }
  }
  tags = {
    Environment = "Test"
  }
}

```



## Arguments Reference

The following arguments are supported:

* `type` - (Required) The type of the resource. This should be set to `Microsoft.OperationsManagement/solutions@api-version`. The available api-versions for this resource are: [`2015-11-01-preview`].

* `parent_id` - (Required) The ID of the azure resource in which this resource is created. The allowed values are:  
  `/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}`

* `name` - (Required) Specifies the name of the azure resource. Changing this forces a new resource to be created.

* `body` - (Required) Specifies the configuration of the resource. More information about the arguments in `body` can be found in the [Microsoft documentation](https://learn.microsoft.com/en-us/azure/templates/Microsoft.OperationsManagement/solutions?pivots=deployment-language-terraform).

For other arguments, please refer to the [azapi_resource](https://registry.terraform.io/providers/Azure/azapi/latest/docs/resources/resource) documentation.

## Import

 ```shell
 # Azure resource can be imported using the resource id, e.g.
 terraform import azapi_resource.example /subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.OperationsManagement/solutions/{resourceName}
 
 # It also supports specifying API version by using the resource id with api-version as a query parameter, e.g.
 terraform import azapi_resource.example /subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.OperationsManagement/solutions/{resourceName}?api-version=2015-11-01-preview
 ```
