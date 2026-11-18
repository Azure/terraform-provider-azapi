---
subcategory: "Microsoft.NetworkSecurity - Azure Adaptive Virtual Patching"
page_title: "neuroShields"
description: |-
  Manages a Adaptive Shield.
---

# Microsoft.NetworkSecurity/neuroShields - Adaptive Shield

This article demonstrates how to use `azapi` provider to manage the Adaptive Shield resource in Azure.



## Example Usage

### basic

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
  default = "centraluseuap"
}

resource "azapi_resource" "resourceGroup" {
  type     = "Microsoft.Resources/resourceGroups@2020-06-01"
  name     = var.resource_name
  location = var.location
}

resource "azapi_resource" "neuroShield" {
  type                      = "Microsoft.NetworkSecurity/neuroShields@2026-10-01-preview"
  parent_id                 = azapi_resource.resourceGroup.id
  name                      = var.resource_name
  location                  = var.location
  schema_validation_enabled = false

  body = {
    properties = {
      responseMode      = "Manual"
      severityThreshold = "Medium"
      description       = "Adaptive virtual patching example"
    }
  }
}

```



## Arguments Reference

The following arguments are supported:

* `type` - (Required) The type of the resource. This should be set to `Microsoft.NetworkSecurity/neuroShields@api-version`. The available api-versions for this resource are: [`2026-10-01-preview`].

* `parent_id` - (Required) The ID of the azure resource in which this resource is created. The allowed values are:  
  `/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}`

* `name` - (Required) Specifies the name of the azure resource. Changing this forces a new resource to be created.

* `body` - (Required) Specifies the configuration of the resource. More information about the arguments in `body` can be found in the [Microsoft documentation](https://learn.microsoft.com/en-us/azure/templates/Microsoft.NetworkSecurity/neuroShields?pivots=deployment-language-terraform).

For other arguments, please refer to the [azapi_resource](https://registry.terraform.io/providers/Azure/azapi/latest/docs/resources/resource) documentation.

## Import

 ```shell
 # Azure resource can be imported using the resource id, e.g.
 terraform import azapi_resource.example /subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.NetworkSecurity/neuroShields/{resourceName}
 
 # It also supports specifying API version by using the resource id with api-version as a query parameter, e.g.
 terraform import azapi_resource.example /subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.NetworkSecurity/neuroShields/{resourceName}?api-version=2026-10-01-preview
 ```
