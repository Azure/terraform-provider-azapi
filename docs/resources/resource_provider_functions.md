---
subcategory: ""
layout: "azapi"
page_title: "Azure Functions: Resource Provider Functions"
description: |-
  Provides functions to build and parse Azure resource IDs.
---

# Resource Provider Functions

Azure resource IDs are unique identifiers for resources in Azure, used to manage and reference resources within the platform. This document provides functions to build and parse Azure resource IDs for different scopes, including tenant, subscription, management group, resource group, and extensions.

## Basic Functions

### `build_resource_id`

Constructs an Azure resource ID given the parent ID, resource type, and resource name. This function is useful for creating resource IDs for top-level and nested resources within a specific scope.

#### Arguments

* `parent_id` - (Required) The parent ID of the Azure resource. This is the ID of the resource that contains the resource being created.
* `resource_type` - (Required) The resource type of the Azure resource.
* `name` - (Required) The name of the Azure resource.

#### Attributes

* `resource_id` - The constructed Azure resource ID.

#### Example Usage

```hcl
locals {
  parent_id = "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/myResourceGroup"
  resource_type = "Microsoft.Network/virtualNetworks"
  name = "myVNet"
}

output "resource_id" {
  value = build_resource_id(local.parent_id, local.resource_type, local.name)
}
```

### `parse_resource_id`

Parses an Azure resource ID into its components.

#### Arguments

* `resource_type` - (Required) The resource type of the Azure resource.
* `resource_id` - (Required) The resource ID of the Azure resource to parse.

#### Attributes

* `id` - The full resource ID.
* `type` - The resource type.
* `name` - The name of the resource.
* `parent_id` - The parent resource ID.
* `resource_group_name` - The resource group name.
* `subscription_id` - The subscription ID.
* `provider_namespace` - The provider namespace.
* `parts` - A map of the parsed components of the resource ID.

#### Example Usage

```hcl
locals {
  resource_id = "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/myResourceGroup/providers/Microsoft.Network/virtualNetworks/myVNet"
  resource_type = "Microsoft.Network/virtualNetworks"
}

output "parsed_resource_id" {
  value = parse_resource_id(local.resource_type, local.resource_id)
}
```

## Specialized Provider Functions

### `tenant_resource_id`

Constructs an Azure resource ID for tenant scope resources given the resource type and resource names.

#### Arguments

* `resource_type` - (Required) The resource type of the Azure resource.
* `resource_names` - (Required) A list of resource names to construct the resource ID.

#### Attributes

* `resource_id` - The constructed Azure resource ID.

#### Example Usage

```hcl
locals {
  resource_type = "Microsoft.Billing/billingAccounts/billingProfiles"
  resource_names = ["ba1", "bp1"]
}

output "tenant_resource_id" {
  value = tenant_resource_id(local.resource_type, local.resource_names)
}
```

### `subscription_resource_id`

Constructs an Azure resource ID for subscription scope resources given the subscription ID, resource type, and resource names.

#### Arguments

* `subscription_id` - (Required) The subscription ID of the Azure resource.
* `resource_type` - (Required) The resource type of the Azure resource.
* `resource_names` - (Required) A list of resource names to construct the resource ID.

#### Attributes

* `resource_id` - The constructed Azure resource ID.

#### Example Usage

```hcl
locals {
  subscription_id = "00000000-0000-0000-0000-000000000000"
  resource_type = "Microsoft.Resources/resourceGroups"
  resource_names = ["rg1"]
}

output "subscription_resource_id" {
  value = subscription_resource_id(local.subscription_id, local.resource_type, local.resource_names)
}
```

### `management_group_resource_id`

Constructs an Azure resource ID for management group scope resources given the management group name, resource type, and resource names.

#### Arguments

* `management_group_name` - (Required) The name of the management group.
* `resource_type` - (Required) The resource type of the Azure resource.
* `resource_names` - (Required) A list of resource names to construct the resource ID.

#### Attributes

* `resource_id` - The constructed Azure resource ID.

#### Example Usage

```hcl
locals {
  management_group_name = "mg1"
  resource_type = "Microsoft.Billing/billingAccounts/billingProfiles"
  resource_names = ["ba1", "bp1"]
}

output "management_group_resource_id" {
  value = management_group_resource_id(local.management_group_name, local.resource_type, local.resource_names)
}
```

### `resource_group_resource_id`

Constructs an Azure resource ID for resource group scope resources given the subscription ID, resource group name, resource type, and resource names.

#### Arguments

* `subscription_id` - (Required) The subscription ID of the Azure resource.
* `resource_group_name` - (Required) The name of the resource group.
* `resource_type` - (Required) The resource type of the Azure resource.
* `resource_names` - (Required) A list of resource names to construct the resource ID.

#### Attributes

* `resource_id` - The constructed Azure resource ID.

#### Example Usage

```hcl
locals {
  subscription_id = "00000000-0000-0000-0000-000000000000"
  resource_group_name = "rg1"
  resource_type = "Microsoft.Network/virtualNetworks/subnets"
  resource_names = ["vnet1", "subnet1"]
}

output "resource_group_resource_id" {
  value = resource_group_resource_id(local.subscription_id, local.resource_group_name, local.resource_type, local.resource_names)
}
```

### `extension_resource_id`

Constructs an Azure extension resource ID given the base resource ID, extension name, and additional resource names.

#### Arguments

* `resource_id` - (Required) The base resource ID of the Azure resource.
* `extension_name` - (Required) The name of the extension resource.
* `resource_names` - (Required) A list of resource names to construct the extension resource ID.

#### Attributes

* `resource_id` - The constructed Azure extension resource ID.

#### Example Usage

```hcl
locals {
  resource_id = "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/rg1/providers/Microsoft.Network/virtualNetworks/vnet1"
  extension_name = "Microsoft.Authorization/locks"
  resource_names = ["mylock"]
}

output "extension_resource_id" {
  value = extension_resource_id(local.resource_id, local.extension_name, local.resource_names)
}
```
