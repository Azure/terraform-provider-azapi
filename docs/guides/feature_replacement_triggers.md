---
layout: "azapi"
page_title: "Feature: Triggers for Resource Replacement"
description: |-
  This guide will cover how to use the triggers for resource replacement feature in the AzAPI provider. 

---

There are scenarios where you need to replace a resource when a specific condition is met. For example, you may need to replace a virtual machine when the OS type is changed from `Linux` to `Windows`. The AzAPI provider supports these scenarios with the triggers for resource replacement feature.

## replace_triggers_external_values

The `replace_triggers_external_values` attribute will trigger a resource replacement when the value changes and is not `null`.

For example, the following configuration will replace the public IP address when the `zones` value changes:

```hcl
resource "azapi_resource" "publicIPAddress" {
  type      = "Microsoft.Network/publicIPAddresses@2023-11-01"
  parent_id = "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/example"
  name      = var.name
  body = {
    properties = {
      publicIPAllocationMethod = "Dynamic"
    }
    zones = var.zones
  }
  replace_triggers_external_values = [
    var.zones,
  ]
}
```

It is recommended to use `variables` to define the values that trigger the resource replacement. This way, you can easily change the values without modifying the resource configuration.

## replace_triggers_refs

The `replace_triggers_refs` attribute accepts a list of `JMESPath` expressions that query the resource attributes. The resource will be replaced when the query result changes.

For example, the following configuration will replace the public IP address when the `zones` value changes:

```hcl
resource "azapi_resource" "publicIPAddress" {
  type      = "Microsoft.Network/publicIPAddresses@2023-11-01"
  parent_id = "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/example"
  name      = var.name
  body = {
    properties = {
      publicIPAllocationMethod = "Dynamic"
    }
    zones = var.zones
  }
  replace_triggers_refs = [
    "zones",
  ]
}
```
