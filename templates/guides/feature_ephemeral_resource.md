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

```
ephemeral.azapi_resource_action.listKeys.output.primary_access_key
```

More information about the ephemeral resource can be found in the [Terraform documentation](https://developer.hashicorp.com/terraform/language/resources/ephemeral).

## Write-Only Attributes

Write-only arguments let you securely pass temporary values to Terraform's managed resources during an operation without persisting those values to state or plan files.

The AzAPI provider supports `sensitive_body` to define attributes that are not stored in the state file, ensuring sensitive information remains secure.

### Example

The below example demonstrates how to use the `sensitive_body` attribute to create a storage insight configuration.

The `properties.storageAccount.key` acceptes the storage account key, which is sensitive information. To avoid storing this key in the state file, it was retrieved using the `azapi_resource_action` ephemeral resource, and passed to the `sensitive_body` attribute.

The `sensitive_body` and `body` attributes will be merged into the request body.


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
  sensitive_body = {
    properties = {
      storageAccount = {
        key = ephemeral.azapi_resource_action.listKeys.output.primary_access_key
      }
    }
  }
}
```

The `sensitive_body` attribute detects the changes based on the SHA256 hash of the value. If the value changes, the resource will be updated. However, ephemeral resources like `random_password` are always considered changed, to avoid unnecessary updates, `azapi_resource` supports the `sensitive_body_version` attribute. This attribute can be used to manually control the version of the sensitive body, allowing you to avoid unnecessary updates.

### Example with `sensitive_body_version`

In the following example, we create an Ubuntu virtual machine using the `azapi_resource` resource. The `sensitive_body` attribute is used to pass the admin username and password securely, while the `sensitive_body_version` attribute is used to control the versioning of these sensitive attributes.

The `sensitive_body_version` attribute is a map where the key is the path to the sensitive attribute, and the value is the version of that attribute. The key is a string in the format of `path.to.property[index].subproperty`, where `index` is the index of the item in an array. This allows you to specify different versions for different sensitive attributes. When the version changes, the sensitive attribute will be included in the request body, otherwise it will be omitted from the request body. 

Please note that if the `senstive_body_version` is provided, the SHA256 hash of the value will not be calculated, and the resource will not be updated unless the version changes.

```hcl
resource "azapi_resource" "ubuntuVM" {
  type      = "Microsoft.Compute/virtualMachines@2020-12-01"
  parent_id = azapi_resource.resourceGroup.id
  name      = "name"
  location  = "location"
  body = {
    properties = {
      hardwareProfile = {
        vmSize = "Standard_A2_v2"
      }
    }
  }
  sensitive_body = {
    properties = {
      osProfile = {
        adminPassword = random_password.password.result
        adminUsername = random_password.username.result
      }
    }
  }
  sensitive_body_version = {
    "properties.osProfile.adminPassword" = "v2"
    "properties.osProfile.adminUsername" = "v1"
  }
}
```
