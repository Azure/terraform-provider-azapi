---
subcategory: "Microsoft.App - Azure Container Apps"
page_title: "managedEnvironments/privateEndpoint"
description: |-
  Manages a Container App Environment Private Endpoint.
---

# Microsoft.App/managedEnvironments/privateEndpoint - Container App Environment Private Endpoint

This article demonstrates how to use `azapi` provider to manage the Container App Environment Private Endpoint resource in Azure.

## Example Usage



## Arguments Reference

The following arguments are supported:

* `type` - (Required) The type of the resource. This should be set to `Microsoft.App/managedEnvironments/privateEndpoint@api-version`. The available api-versions for this resource are: [].

* `parent_id` - (Required) The ID of the azure resource in which this resource is created. The allowed values are:  
  

* `name` - (Required) Specifies the name of the azure resource. Changing this forces a new resource to be created.

* `body` - (Required) Specifies the configuration of the resource. More information about the arguments in `body` can be found in the [Microsoft documentation](https://learn.microsoft.com/en-us/azure/templates/Microsoft.App/managedEnvironments/privateEndpoint?pivots=deployment-language-terraform).

For other arguments, please refer to the [azapi_resource](https://registry.terraform.io/providers/Azure/azapi/latest/docs/resources/resource) documentation.

## Import

 ```shell
 # Azure resource can be imported using the resource id, e.g.
 terraform import azapi_resource.example 
 
 # It also supports specifying API version by using the resource id with api-version as a query parameter, e.g.
 terraform import azapi_resource.example ?api-version=API_VERSION
 ```
