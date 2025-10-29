---
subcategory: "Microsoft.Resources - Azure Resource Manager"
page_title: "resourceGroups"
description: |-
  Manages a Resource Group.
---

# Microsoft.Resources/resourceGroups - Resource Group

This article demonstrates how to use `azapi` provider to manage the Resource Group resource in Azure.



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


```



## Arguments Reference

The following arguments are supported:

* `type` - (Required) The type of the resource. This should be set to `Microsoft.Resources/resourceGroups@api-version`. The available api-versions for this resource are: [`2015-11-01`, `2016-02-01`, `2016-07-01`, `2016-09-01`, `2017-05-10`, `2018-02-01`, `2018-05-01`, `2019-03-01`, `2019-05-01`, `2019-05-10`, `2019-07-01`, `2019-08-01`, `2019-10-01`, `2020-06-01`, `2020-08-01`, `2020-10-01`, `2021-01-01`, `2021-04-01`, `2022-09-01`, `2023-07-01`, `2024-03-01`, `2024-07-01`, `2024-11-01`, `2025-03-01`, `2025-04-01`].

* `parent_id` - (Required) The ID of the azure resource in which this resource is created. The allowed values are:  
  `/subscriptions/{subscriptionId}`

* `name` - (Required) Specifies the name of the azure resource. Changing this forces a new resource to be created.

* `body` - (Required) Specifies the configuration of the resource. More information about the arguments in `body` can be found in the [Microsoft documentation](https://learn.microsoft.com/en-us/azure/templates/Microsoft.Resources/resourceGroups?pivots=deployment-language-terraform).

For other arguments, please refer to the [azapi_resource](https://registry.terraform.io/providers/Azure/azapi/latest/docs/resources/resource) documentation.

## Import

 ```shell
 # Azure resource can be imported using the resource id, e.g.
 terraform import azapi_resource.example /subscriptions/{subscriptionId}/providers/Microsoft.Resources/resourceGroups/{resourceName}
 
 # It also supports specifying API version by using the resource id with api-version as a query parameter, e.g.
 terraform import azapi_resource.example /subscriptions/{subscriptionId}/providers/Microsoft.Resources/resourceGroups/{resourceName}?api-version=2025-04-01
 ```
