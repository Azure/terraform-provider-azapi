---
subcategory: ""
layout: "azapi"
page_title: "Azure Resource ID Data Source: azapi_resource_id"
description: |-
  Parses an Azure resource ID into its separate fields.
---

# azapi_resource_list

This resource can parse an Azure resource ID into its separate fields.

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

data "azapi_resource_id" "account" {
  type        = "Microsoft.Automation/automationAccounts@2021-06-22"
  resource_id = "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/rg1/providers/Microsoft.Automation/automationAccounts/automationAccount1"
}

output "account_name" {
  value = data.azapi_resource_id.account.name
}

output "account_resource_group" {
  value = data.azapi_resource_id.account.resource_group_name
}

output "account_subscription" {
  value = data.azapi_resource_id.account.subscription_id
}

output "account_parent_id" {
  value = data.azapi_resource_id.account.parent_id
}

data "azapi_resource_id" "vnet" {
  type      = "Microsoft.Network/virtualNetworks@2021-02-01"
  parent_id = "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/rg1"
  name      = "vnet1"
}

output "vnet_id" {
  value = data.azapi_resource_id.vnet.id
}

output "vnet_resource_group" {
  value = data.azapi_resource_id.vnet.resource_group_name
}

output "vnet_subscription" {
  value = data.azapi_resource_id.vnet.subscription_id
}
```

## Arguments Reference

The following arguments are supported:

* `type` - (Required) It is in a format like `<resource-type>@<api-version>`. `<resource-type>` is the Azure resource type, for example, `Microsoft.Storage/storageAccounts`.
  `<api-version>` is version of the API used to manage this azure resource.

* `name` - (Optional) Specifies the name of the azure resource.

* `parent_id` - (Optional) The ID of the azure resource in which this resource is created. It supports different kinds of deployment scope for **top level** resources:
  - resource group scope: `parent_id` should be the ID of a resource group, it's recommended to manage a resource group by [azurerm_resource_group](https://registry.terraform.io/providers/hashicorp/azurerm/latest/docs/resources/resource_group).
  - management group scope: `parent_id` should be the ID of a management group, it's recommended to manage a management group by [azurerm_management_group](https://registry.terraform.io/providers/hashicorp/azurerm/latest/docs/resources/management_group).
  - extension scope: `parent_id` should be the ID of the resource you're adding the extension to.
  - subscription scope: `parent_id` should be like `/subscriptions/00000000-0000-0000-0000-000000000000`
  - tenant scope: `parent_id` should be `/`

  For child level resources, the `parent_id` should be the ID of its parent resource, for example, subnet resource's `parent_id` is the ID of the vnet.

* `resource_id` - (Optional) The ID of an existing azure source.

~> **Note:** Configuring `name` and `parent_id` is an alternative way to configure `resource_id`.


## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the azure resource.

* `name` - The name of the azure resource.

* `parent_id` - The ID of the azure resource in which this resource is created.

* `provider_namespace` - The azure resource provider namespace of the azure resource.

* `resource_group_name` - The resource group name of the azure resource.

* `subscription_id` - The subscription ID of the azure resource.

* `parts` - The map of the resource ID parts, where the key is the part name and the value is the part value. e.g. `virtualNetworks=myVnet`.


## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `read` - (Defaults to 30 minutes) Used when retrieving the azure resource.
