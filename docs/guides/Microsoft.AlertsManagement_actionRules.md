---
subcategory: "Microsoft.AlertsManagement - Azure Monitor"
page_title: "actionRules"
description: |-
  Manages a Alert Processing Rule which apply action group.
---

# Microsoft.AlertsManagement/actionRules - Alert Processing Rule which apply action group

This article demonstrates how to use `azapi` provider to manage the Alert Processing Rule which apply action group resource in Azure.

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

resource "azapi_resource" "actionRule" {
  type      = "Microsoft.AlertsManagement/actionRules@2021-08-08"
  parent_id = azapi_resource.resourceGroup.id
  name      = var.resource_name
  location  = "global"
  body = {
    properties = {
      actions = [
        {
          actionType = "RemoveAllActionGroups"
        },
      ]
      description = ""
      enabled     = true
      scopes = [
        azapi_resource.resourceGroup.id,
      ]
    }

  }
  schema_validation_enabled = false
  response_export_values    = ["*"]
}


```



## Arguments Reference

The following arguments are supported:

* `type` - (Required) The type of the resource. This should be set to `Microsoft.AlertsManagement/actionRules@api-version`. The available api-versions for this resource are: [`2019-05-05-preview`, `2021-08-08`, `2021-08-08-preview`].

* `parent_id` - (Required) The ID of the azure resource in which this resource is created. The allowed values are:  
  `/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}`

* `name` - (Required) Specifies the name of the azure resource. Changing this forces a new resource to be created.

* `body` - (Required) Specifies the configuration of the resource. More information about the arguments in `body` can be found in the [Microsoft documentation](https://learn.microsoft.com/en-us/azure/templates/Microsoft.AlertsManagement/actionRules?pivots=deployment-language-terraform).

## Import

 ```shell
 # Azure resource can be imported using the resource id, e.g.
 terraform import azapi_resource.example /subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.AlertsManagement/actionRules/{resourceName}
 
 # It also supports specifying API version by using the resource id with api-version as a query parameter, e.g.
 terraform import azapi_resource.example /subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.AlertsManagement/actionRules/{resourceName}?api-version=2021-08-08-preview
 ```
