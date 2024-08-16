---
layout: "azapi"
page_title: "AzAPI Provider: Frequently asked questions"
description: |-
  This guide will cover some frequently asked questions
---

## After applying the configuration, running terraform plan found there's still a change?

More context:
```
Terraform used the selected providers to generate the following execution plan. Resource actions are indicated with the following symbols:
  ~ update in-place

Terraform will perform the following actions:

  # azapi_resource.test will be updated in-place
  ~ resource "azapi_resource" "test" {
      ~ body                            = jsonencode(
          ~ {
              ~ properties = {
                  + accountKey  = "TOI************QqA=="
                    # (2 unchanged elements hidden)
      ~ output                          = jsonencode({}) -> (known after apply)
        tags                            = {}
        # (5 unchanged attributes hidden)
    }

Plan: 0 to add, 1 to change, 0 to destroy.
```

In some cases, when `body` refers to some sensitive properties, user may not get a clear plan as above. This happens when a property contains sensitive credential like a storage account's access key, the value won't be returned by design. 

Please add `ignore_missing_property = true` to the resource block and apply it.

Similarly, values in incorrect casing might be returned and cause a diff. Please use `ignore_casing = true` to suppress it.

## How to manage resources which can be managed by both parent resource API and its own API like `Microsoft.Network/virtualNetworks/subnets`?

For example, `Microsoft.Network/virtualNetworks/subnets` has its own API and it can also be managed by `Microsoft.Network/virtualNetworks` API.

It's recommendded to manage `Microsoft.Network/virtualNetworks/subnets` in `Microsoft.Network/virtualNetworks` instead of using its own API.

The reason is, if subnets are defined separately, when update the vnet which has no definition of its subnets, 
its request body won't contain any subnets definitions, so existing subnets will be removed.


## How to use API/properties which is not in embeded schema?

`AzApi` provider will use embedded schema to verify the inputs, to skip the validation, please add `schema_validation_enabled = false` to the resource block.


## How to manage properties which can only be set during update?

More context: Properties like CMK(customer managed key) can't be set when create this resource. Because CMK depends on key vault key which depends on this resource's identity.

It can be solved by using `azapi_update_resource` to perform a multi-steps deployment, here's an [example](https://github.com/Azure/terraform-provider-azapi/tree/main/examples/Microsoft.ServiceBus/ServiceBusNamespace-CMK/main.tf).

## What are the types for the resource group, subscription and tenant?

| Resource Type | Type | Example | Comment |
| --- | --- | --- | --- |
| resource group | `Microsoft.Resources/resourceGroups` | `/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/rg` | |
| subscription | `Microsoft.Resources/subscriptions` | `/subscriptions/00000000-0000-0000-0000-000000000000` | |
| tenant | `Microsoft.Resources/tenants` | `/tenants/00000000-0000-0000-0000-000000000000` | |
| management group | `Microsoft.Management/managementGroups` | `/providers/Microsoft.Management/managementGroups/00000000-0000-0000-0000-000000000000` | |
| generic provider | `Microsoft.Resources/providers` | `<scope>/providers/<RP>` | This is the generic provider resource type whose operations are defined by ARM, and is supported by all of Azure RPs (e.g. actions like `register`, `providerPermissions`, etc.) |
| concrete provider | `Microsoft.Resources/providers` | `<scope>/providers/Microsoft.Storage` | This is the concrete provider whose operations are only defined for this provider itself |
