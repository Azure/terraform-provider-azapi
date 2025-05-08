---
subcategory: "{{.subcategory}}"
page_title: "{{.page_title}}"
description: |-
  Manages a {{.resource_type_friendly_name}}.
---

# {{.resource_type}} - {{.resource_type_friendly_name}}

This article demonstrates how to use `azapi` provider to manage the {{.resource_type_friendly_name}} resource in Azure.

## Example Usage

{{.example}}

## Arguments Reference

The following arguments are supported:

* `type` - (Required) The type of the resource. This should be set to `{{.resource_type}}@api-version`. The available api-versions for this resource are: [{{.api_versions}}].

* `parent_id` - (Required) The ID of the azure resource in which this resource is created. The allowed values are:  
  {{.parent_id}}

* `name` - (Required) Specifies the name of the azure resource. Changing this forces a new resource to be created.

* `body` - (Required) Specifies the configuration of the resource. More information about the arguments in `body` can be found in the [Microsoft documentation]({{.reference_link}}).

## Import

 ```shell
 # Azure resource can be imported using the resource id, e.g.
 terraform import azapi_resource.example {{.resource_id}}
 
 # It also supports specifying API version by using the resource id with api-version as a query parameter, e.g.
 terraform import azapi_resource.example {{.resource_id}}?api-version={{.api_version}}
 ```
