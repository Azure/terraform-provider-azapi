---
layout: "azapi"
page_title: "Feature: List Resources"
description: |-
  This guide will cover how to use `azapi_resource_list` data source to list resources in Azure.

---

The `azapi_resource_list` data source is used to list resources in Azure. This data source is useful when you need to get a list of resources that match a specific criteria, such as all virtual machines in a subscription or all storage accounts in a resource group.

## List Resources by Subscription

The following example demonstrates how to list all virtual machines in a subscription:

```hcl
data "azapi_resource_list" "listVirtualMachinesBySubscription" {
  type      = "Microsoft.Compute/virtualMachines@2020-06-01"
  parent_id = "/subscriptions/00000000-0000-0000-0000-000000000000"
}
```

The `output` attribute contains the result of the query, if default output is enabled(default output feature is introduced and by default enabled in version 2.1.0, you can disable it by setting `disable_default_output = true`):

```hcl
{
  value = [
    {
      id = "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/myResourceGroup/providers/Microsoft.Compute/virtualMachines/myVirtualMachine"
      name = "myVirtualMachine"
      type = "Microsoft.Compute/virtualMachines"
      location = "eastus"
      tags = {}
      properties = {
        // more properties
      }
    },
    {
      id = "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/myResourceGroup/providers/Microsoft.Compute/virtualMachines/myVirtualMachine2"
      name = "myVirtualMachine2"
      type = "Microsoft.Compute/virtualMachines"
      location = "eastus"
      tags = {}
      properties = {
        // more properties
      }
    }
  ]
}
```

## List Resources by Resource Group

The following example demonstrates how to list all virtual machines in a resource group:

```hcl
data "azapi_resource_list" "listVirtualMachinesByResourceGroup" {
  type      = "Microsoft.Compute/virtualMachines@2020-06-01"
  parent_id = "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/myResourceGroup"
}
```

## List Resources by Parent Resource

The following example demonstrates how to list all subnets in a virtual network:

```hcl
data "azapi_resource_list" "listSubnetsByVirtualNetwork" {
  type      = "Microsoft.Network/virtualNetworks/subnets@2020-06-01"
  parent_id = "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/myResourceGroup/providers/Microsoft.Network/virtualNetworks/myVirtualNetwork"
}
```

## List Extension Resources

The following example demonstrates how to list all role definitions for a storage account:

```hcl
data "azapi_resource_list" "roleDefinitions" {
  type      = "Microsoft.Authorization/roleDefinitions@2022-04-01"
  parent_id = "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/myResourceGroup/providers/Microsoft.Storage/storageAccounts/myStorageAccount"
}
```

## Filter Resources

There are several ways to filter resources using the `azapi_resource_list` data source. You can filter resources by HCL expressions and by using the `JMESPath` query language.

### Filter by HCL Expressions

With default output enabled, you can filter resources by HCL expressions. For example, to filter virtual machines by location:

```hcl
data "azapi_resource_list" "listVirtualMachinesBySubscription" {
  type      = "Microsoft.Compute/virtualMachines@2020-06-01"
  parent_id = "/subscriptions/00000000-0000-0000-0000-000000000000"
}

output "vmIds" {
  value = [for vm in data.azapi_resource_list.listVirtualMachinesBySubscription.output.value : vm.id if vm.location == "southeastasia"]
}
```

It uses the `for` expression to iterate over the list of virtual machines and filter them by location.

### Filter by JMESPath

You can also filter resources using the `JMESPath` query language. For example, to filter virtual machines by location:

```hcl
data "azapi_resource_list" "listVirtualMachinesBySubscription" {
  type      = "Microsoft.Compute/virtualMachines@2020-06-01"
  parent_id = "/subscriptions/00000000-0000-0000-0000-000000000000"
  response_export_values = {
    vmIds = "value[?location == 'southeastasia'].id"
  }
}

output "vmIds" {
  value = data.azapi_resource_list.listVirtualMachinesBySubscription.output.vmIds
}
```

It uses the `response_export_values` attribute to specify the `JMESPath` query to filter virtual machines by location. More details about using `JMESPath` query language to filter resources can be found [here](./feature_jmes_query.html).
