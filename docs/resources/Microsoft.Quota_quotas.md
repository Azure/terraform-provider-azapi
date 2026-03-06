---
subcategory: "Microsoft.Quota - Azure Quota"
page_title: "quotas"
description: |-
  Manages a Azure Quota.
---

# Microsoft.Quota/quotas - Azure Quota

This article demonstrates how to use `azapi` provider to manage the Azure Quota resource in Azure.



## Example Usage



## Arguments Reference

The following arguments are supported:

* `type` - (Required) The type of the resource. This should be set to `Microsoft.Quota/quotas@api-version`. The available api-versions for this resource are: [`2021-03-15-preview`, `2023-02-01`, `2023-06-01-preview`, `2024-10-15-preview`, `2024-12-18-preview`, `2025-03-01`, `2025-03-15-preview`, `2025-07-15`, `2025-09-01`].

* `parent_id` - (Required) The ID of the azure resource in which this resource is created. The allowed values are:  
  `/`  
  `/providers/Microsoft.Management/managementGroups/{managementGroupId}`  
  `/subscriptions/{subscriptionId}`  
  `/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}`  
  `{any azure resource id}`

* `name` - (Required) Specifies the name of the azure resource. Changing this forces a new resource to be created.

* `body` - (Required) Specifies the configuration of the resource. More information about the arguments in `body` can be found in the [Microsoft documentation](https://learn.microsoft.com/en-us/azure/templates/Microsoft.Quota/quotas?pivots=deployment-language-terraform).

For other arguments, please refer to the [azapi_resource](https://registry.terraform.io/providers/Azure/azapi/latest/docs/resources/resource) documentation.

## Import

 ```shell
 # Azure resource can be imported using the resource id, e.g.
 terraform import azapi_resource.example //providers/Microsoft.Quota/quotas/{resourceName}
 
 # It also supports specifying API version by using the resource id with api-version as a query parameter, e.g.
 terraform import azapi_resource.example //providers/Microsoft.Quota/quotas/{resourceName}?api-version=2025-09-01
 ```
