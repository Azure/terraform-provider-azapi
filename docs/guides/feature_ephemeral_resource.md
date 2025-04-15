---
layout: "azapi"
page_title: "Feature: Handle Sensitive Data with Ephemeral Resource and Write-Only Attributes"
description: |-
  This guide will cover how to use the ephemeral resource and write-only attributes to handle sensitive data in the AzAPI provider.

---

The AzAPI provider supports the use of ephemeral resources and write-only attributes to handle sensitive data. This feature is particularly useful when you need to manage sensitive information, such as passwords or secrets, without exposing them in your Terraform state file or logs.

This article will cover how to use the ephemeral resource and write-only attributes in the AzAPI provider.

## Prerequisites

- [Terraform AzAPI provider](https://registry.terraform.io/providers/azure/azapi) version 2.4.0 or later
- [Terraform CLI](https://www.terraform.io/downloads.html) version 1.11.0 or later

## Ephemeral Resource

The `azapi_resource_action` ephemeral resource allows you to create a resource that is not persisted in the Terraform state file. This is useful for retrieving sensitive data without storing it in the state file.

### Example

The below example demonstrates how to use the `azapi_resource_action` resource to retrieve the storage account access key.

```hcl
ephemeral "azapi_resource_action" "listKeys" {
  type        = "Microsoft.Storage/storageAccounts@2024-01-01"
  resource_id = azapi_resource.storageAccount.id
  action      = "listKeys"
  response_export_values = {
    primary_access_key = "keys[0].value"
  }
}
```

To use the property `primary_access_key`, you can use the following code:

```hcl
ephemeral.azapi_resource_action.listKeys.output.primary_access_key
```

More information about the ephemeral resource can be found in the [Terraform documentation](https://developer.hashicorp.com/terraform/language/resources/ephemeral).

## Write-Only Attributes

Write-only arguments let you securely pass temporary values to Terraform's managed resources during an operation without persisting those values to state or plan files.

The AzAPI provider supports `write_only_body` to define attributes that are not stored in the state file, ensuring sensitive information remains secure.

### Example

The below example demonstrates how to use the `write_only_body` attribute to create a storage insight configuration.

The `properties.storageAccount.key` acceptes the storage account key, which is sensitive information. To avoid storing this key in the state file, it was retrieved using the `azapi_resource_action` ephemeral resource, and passed to the `write_only_body` attribute.

The `write_only_body` and `body` attributes will be merged into the request body.


```hcl
resource "azapi_resource" "storageInsightConfig" {
  type      = "Microsoft.OperationalInsights/workspaces/storageInsightConfigs@2023-09-01"
  parent_id = azapi_resource.workspace.id
  name      = "storageInsightConfig"
  body = {
    properties = {
      storageAccount = {
        id = azapi_resource.storageAccount.id
      }
    }
  }
  write_only_body = {
    properties = {
      storageAccount = {
        key = ephemeral.azapi_resource_action.listKeys.output.primary_access_key
      }
    }
  }
}
```