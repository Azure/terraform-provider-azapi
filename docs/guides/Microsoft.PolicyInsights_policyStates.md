---
subcategory: "Microsoft.PolicyInsights - Azure Policy"
page_title: "policyStates"
description: |-
  Manages a Policy Insights Policy States.
---

# Microsoft.PolicyInsights/policyStates - Policy Insights Policy States

This article demonstrates how to use `azapi` provider to manage the Policy Insights Policy States resource in Azure.

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

resource "azapi_resource" "resourceGroup" {
  type     = "Microsoft.Resources/resourceGroups@2020-06-01"
  name     = var.resource_name
  location = var.location
}

resource "azapi_resource_action" "triggerEvaluation" {
  type        = "Microsoft.PolicyInsights/policyStates@2019-10-01"
  resource_id = "${azapi_resource.resourceGroup.id}/providers/Microsoft.PolicyInsights/policyStates/latest"
  action      = "triggerEvaluation"
}

```



## Arguments Reference

The following arguments are supported:

* `type` - (Required) The type of the resource. This should be set to `Microsoft.PolicyInsights/policyStates@api-version`. The available api-versions for this resource are: [].

* `parent_id` - (Required) The ID of the azure resource in which this resource is created. The allowed values are:  
  

* `name` - (Required) Specifies the name of the azure resource. Changing this forces a new resource to be created.

* `body` - (Required) Specifies the configuration of the resource. More information about the arguments in `body` can be found in the [Microsoft documentation](https://learn.microsoft.com/en-us/azure/templates/Microsoft.PolicyInsights/policyStates?pivots=deployment-language-terraform).

## Import

 ```shell
 # Azure resource can be imported using the resource id, e.g.
 terraform import azapi_resource.example 
 
 # It also supports specifying API version by using the resource id with api-version as a query parameter, e.g.
 terraform import azapi_resource.example ?api-version=API_VERSION
 ```
