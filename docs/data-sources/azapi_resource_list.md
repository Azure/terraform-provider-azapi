---
subcategory: ""
layout: "azapi"
page_title: "Azure Resource List Data Source: azapi_resource_list"
description: |-
  List all resources of a specific type under a scope.
---

# azapi_resource_list

This resource can list all resources of a specific type under a scope. If the API supports paging, it will automatically fetch all pages and return the full list.

## Example Usage

```hcl
terraform {
  required_providers {
    azapi = {
      source = "Azure/azapi"
    }
  }
}

provider "azapi" {
}

data "azapi_resource_list" "listBySubscription" {
  type                   = "Microsoft.Automation/automationAccounts@2021-06-22"
  parent_id              = "/subscriptions/00000000-0000-0000-0000-000000000000"
  response_export_values = ["*"]
}

data "azapi_resource_list" "listByResourceGroup" {
  type                   = "Microsoft.Automation/automationAccounts@2021-06-22"
  parent_id              = "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/rg1"
  response_export_values = ["*"]
}

data "azapi_resource_list" "listSubnetsByVnet" {
  type                   = "Microsoft.Network/virtualNetworks/subnets@2021-02-01"
  parent_id              = "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/rg1/providers/Microsoft.Network/virtualNetworks/vnet1"
  response_export_values = ["*"]
}

```

## Arguments Reference

The following arguments are supported:

* `type` - (Required) It is in a format like `<resource-type>@<api-version>`. `<resource-type>` is the Azure resource type, for example, `Microsoft.Storage/storageAccounts`.
  `<api-version>` is version of the API used to manage this azure resource.

* `parent_id` - (Required) The parent resource ID to list resources under. e.g. `/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/myResourceGroup`.

---

* `response_export_values` - (Optional) A list of path that needs to be exported from response body.
  Setting it to `["*"]` will export the full response body.
  Here's an example. If it sets to `["value"]`, it will set the following HCL object to computed property `output`.
```
{
  value = [
    {
      id = "id1"
      Permissions = "Full"
    },
    {
      id = "id2"
      Permissions = "Full"
    }
  ]
}
```

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the azure resource list.

* `output` - The output HCL object containing the properties specified in `response_export_values`. Here are some examples to use the values.


## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `read` - (Defaults to 30 minutes) Used when retrieving the azure resource.
