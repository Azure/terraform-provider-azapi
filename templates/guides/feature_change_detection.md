---
layout: "azapi"
page_title: "Feature: Detect and Manage Resource Changes"
description: |-
  This guide will cover how change detection and management works in the AzAPI provider.

---

## azapi_resource

The `azapi_resource` resource is used to manage the full lifecycle of Azure resources. 

In the `azapi` v2.5.0 release, the `azapi_resource` resource introduced an enhanced change detection feature, it will only show the changes if the remote resource state is different from the desired state defined in your configuration, and the resource must be updated to match the desired state. This feature is designed to reduce noise in the plan output and make it easier to identify actual changes that need to be applied. This feature is enabled by default, and it can be disabled by setting the `ignore_no_op_changes` attribute to `false` in the provider block:

```hcl
provider "azapi" {
  ignore_no_op_changes = false
}
```

When enabled, the `azapi_resource` resource retrieves the current state of the resource with the api-version defined in the configuration, please note that the api-version in the configuration could be different from last applied version. For example, if you changes the api-version in the configuration, and run `terraform plan`, it will make a GET request to retrieve the current state of the resource using the new api-version. Then it compares the current state with the desired state defined in your configuration and similarly the desired state in the configuration could be different from last applied version. It only checks for changes in the properties that are currently defined in the `body` attribute of the resource. Removing properties will not be detected as changes, and adding new properties that are matching the remote state will not be detected as changes, because the remote state is still matching the desired state, there's no need to update the resource.

Here are some examples of how the change detection works:

### Example 1: Updating the api-version

In the following example, we are updating the api-version of a virtual machine resource from `2021-07-01` to `2024-07-01`. The body of the resource remains unchanged. Because azure resource can be managed with different api-versions, and there are no changes in the resource's configuration, the `azapi_resource` will show no changes in the plan output.

```hcl
resource "azapi_resource" "example" {
  # Previously applied configuration
  # type    = "Microsoft.Compute/virtualMachines@2021-07-01"

  # Current configuration
  type = "Microsoft.Compute/virtualMachines@2024-07-01"

  # No changes in the body
  body = {}
}
```


### Example 2: Removing properties in the body

In this example, we are removing some properties that is not needed anymore from the body of a virtual machine resource. The `azapi_resource` will not detect this as a change because the remote state is still matching the desired state defined in your configuration. Therefore, it will show no changes in the plan output.

```hcl
resource "azapi_resource" "example" {
  type = "Microsoft.Compute/virtualMachines@2024-07-01"

  # Previously applied configuration
  # body = {
  #   someNotImportantProperty = "value"
  # }

  # Current configuration
  body = {}
}
```


### Example 3: Adding properties in the body

In this example, we are adding some properties to the body of a virtual machine resource. The `azapi_resource` will not detect this as a change because the remote state is still matching the desired state defined in your configuration. Therefore, it will show no changes in the plan output.

```hcl
resource "azapi_resource" "example" {
  type = "Microsoft.Compute/virtualMachines@2024-07-01"

  # Previously applied configuration
  # body = {}

  # Current configuration
  body = {
    # remote state: someNewProperty = "foo"
    someNewProperty = "foo"
  }
}
```

But in the following example, the added property is different from the remote state, so the `azapi_resource` will detect this as a change and show it in the plan output. This is because the remote state is not matching the desired state defined in your configuration, and the resource must be updated to match the desired state.

```hcl
resource "azapi_resource" "example" {
  type = "Microsoft.Compute/virtualMachines@2024-07-01"

  # Previously applied configuration
  # body = {}

  # Current configuration
  body = {
    # remote state: someNewProperty = "foo"
    someNewProperty = "bar"
  }
}
```

When disabled, the `azapi_resource` will make an GET request to retrieve the current state of the resource and compare it with the desired state defined in your configuration. It only checks for changes in the properties that are defined in the `body` attribute of the resource. Removing or adding properties that are not defined in the `body` attribute will also be detected as changes. If there are any differences, Terraform will show them in the plan output.

It supports options to ignore specific changes:
1. `ignore_missing_property`: This option allows you to ignore properties that are not present in the current remote state but are defined in your configuration. This is useful when you specified credentials or other properties that can not be retrieved from the remote state. This feature is enabled by default.
2. `ignore_casing`: This option allows you to ignore differences in casing for string properties. This is useful when the remote state may have different casing than what you defined in your configuration. This feature is disabled by default.
3. `ignore_null_property`: This option allows you to ignore properties that are set to `null` in your configuration. And these properties will not be included in the request body when creating or updating the resource. This is useful when you want to use conditional logic to include or exclude properties. This feature is disabled by default.

It also supports an enhanced change detection feature for array properties. If the `name` attribute of the array item is set, it will compare the items in the array based on their names to provide a more accurate change detection. If the `name` attribute is not set, it will compare the items based on their order in the array. This feature is enabled by default.

## azapi_update_resource

The `azapi_update_resource` resource is used to update existing resources. For example, if you want to use a feature that is not supported by the `azurerm` provider, you can use the `azapi_update_resource` to update the resource with the new configuration.

The `azapi_update_resource` resource supports change detection. When you run `terraform plan`, it will make an GET request to retrieve the current state of the resource and compare it with the desired state defined in your configuration. It only checks for changes in the properties that are defined in the `body` attribute of the resource. Removing or adding properties that are not defined in the `body` attribute will also be detected as changes.If there are any differences, Terraform will show them in the plan output.

It supports options to ignore specific changes:
1. `ignore_missing_property`: This option allows you to ignore properties that are not present in the current remote state but are defined in your configuration. This is useful when you specified credentials or other properties that can not be retrieved from the remote state. This feature is enabled by default.
2. `ignore_casing`: This option allows you to ignore differences in casing for string properties. This is useful when the remote state may have different casing than what you defined in your configuration. This feature is disabled by default.

It also supports an enhanced change detection feature for array properties. If the `name` attribute of the array item is set, it will compare the items in the array based on their names to provide a more accurate change detection. If the `name` attribute is not set, it will compare the items based on their order in the array. This feature is enabled by default.


## azapi_resource_action

The `azapi_resource_action` does not support change detection. It is designed to perform actions on existing resources, such as starting or stopping a virtual machine. But it could also be used to create a new resource, in which case it will not detect changes.
