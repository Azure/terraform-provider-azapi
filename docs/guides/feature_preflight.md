---
layout: "azapi"
page_title: "Feature: Preflight Validation"
description: |-
  This guide will cover how to use the Preflight Validation feature in the AzAPI provider. Preflight Validation allows you to validate the configuration of your resources before applying changes.

---

Preflight validation is a feature of the AzAPI provider that allows you to validate the configuration of your resources before applying changes. This feature is useful for catching errors early in the development process and ensuring that your resources are configured correctly.

Validation is performed by sending a request to the [Azure Resource Manager (ARM) Resource Validator API](https://learn.microsoft.com/rest/api/resources/resource-validator/validate-resources?view=rest-resources-2020-10-01), which checks the configuration of your resources against the Azure Resource Manager's rules and policies. If there are any errors, Terraform will display an error message with details about the issue.

This guide will cover how to use the Preflight Validation feature in the AzAPI provider.

## Prerequisites

- [Terraform AzAPI provider](https://registry.terraform.io/providers/azure/azapi) version 2.0.1 or later

Enable the Preflight Validation feature by setting the `enable_preflight` attribute to `true` in the provider block:

```hcl
provider "azapi" {
  enable_preflight = true
}
```

- The Preflight validation API requires write permissions. Make sure that the service principal or user running Terraform has the necessary permissions to perform write operations on the resources you are validating.

## Preflight Validation

When you run `terraform plan`, the AzAPI provider will validate the configuration of your resources before applying changes. If there are any errors, Terraform will display an error message with details about the issue.

For example, if you try to create a storage account with a name that is already in use:

```hcl
resource "azapi_resource" "storageAccount" {
  type      = "Microsoft.Storage/storageAccounts@2021-06-01"
  parent_id = azapi_resource.resourceGroup.id
  name      = "example"
  location  = "westus"
  body = {
    "sku" = {
      "name" = "Standard_LRS"
    }
    "kind" = "StorageV2"
  }
}
```

When you run `terraform plan`, you will see an error message like this:

```shell
╷
│ Error: Preflight Validation: Invalid configuration
│ 
│   with azapi_resource.storageAccount,
│   on main.tf line 8, in resource "azapi_resource" "storageAccount":
│    8: resource "azapi_resource" "storageAccount" {
│ 
│ POST https://management.azure.com/providers/Microsoft.Resources/validateResources
│ --------------------------------------------------------------------------------
│ RESPONSE 400: 400 Bad Request
│ ERROR CODE: ResourceValidationFailed
│ --------------------------------------------------------------------------------
│ {
│   "error": {
│     "code": "ResourceValidationFailed",
│     "message": "Resource validation failed, correlation id: '8258ded6-68bb-45da-e2ab-1ff991519381', see details for more information.",
│     "details": [
│       {
│         "code": "PreflightValidationCheckFailed",
│         "message": "Preflight validation failed. Please refer to the details for the specific errors.",
│         "details": [
│           {
│             "code": "StorageAccountAlreadyTaken",
│             "target": "example",
│             "message": "The storage account named example is already taken."
│           }
│         ]
│       }
│     ]
│   }
│ }
```

### Enhanced Validation

The preflight validation provides enhanced validation to catch configuration errors that are not caught by embedded schema validation. For example, if you try to create a virtual network with an invalid CIDR(it should be `/16` instead of `/160`):

```hcl
resource "azapi_resource" "virtualNetwork" {
  type      = "Microsoft.Network/virtualNetworks@2019-11-01"
  parent_id = azapi_resource.resourceGroup.id
  name      = "acctestvn"
  location  = "westus"
  body = {
    properties = {
      addressSpace = {
        addressPrefixes = [
          "10.0.0.0/160"
        ]
      }
    }
  }
}
```

When you run `terraform plan`, you will see an error message like this:

```shell
╷
│ Error: Preflight Validation: Invalid configuration
│ 
│   with azapi_resource.virtualNetwork,
│   on main.tf line 8, in resource "azapi_resource" "virtualNetwork":
│    8: resource "azapi_resource" "virtualNetwork" {
│ 
│ POST https://management.azure.com/providers/Microsoft.Resources/validateResources
│ --------------------------------------------------------------------------------
│ RESPONSE 400: 400 Bad Request
│ ERROR CODE: ResourceValidationFailed
│ --------------------------------------------------------------------------------
│ {
│   "error": {
│     "code": "ResourceValidationFailed",
│     "message": "Resource validation failed, correlation id: 'fbdbf8d1-c491-1f4e-ceb6-ccf47981ea04', see details for more information.",
│     "details": [
│       {
│         "code": "InvalidAddressPrefixFormat",
│         "target": "/subscriptions/000000/resourceGroups/jvongchd/providers/Microsoft.Network/virtualNetworks/acctestvn",
│         "message": "Address prefix 10.0.0.0/160 of resource /subscriptions/000000/resourceGroups/jvongchd/providers/Microsoft.Network/virtualNetworks/acctestvn is not formatted correctly. It should follow CIDR notation, for example 10.0.0.0/24.",
│         "details": []
│       }
│     ]
│   }
│ }
│ --------------------------------------------------------------------------------
│ 
╵
```

### Check Policy Restrictions

The preflight validation also checks for policy restrictions. For example, if you try to create a storage account with `allowBlobPublicAccess` set to `true` and the policy does not allow it:

```hcl
resource "azapi_resource" "storageaccount" {
  type      = "Microsoft.Storage/storageAccounts@2021-02-01"
  parent_id = azapi_resource.resourceGroup.id
  name      = "acctestsa1227"
  location  = "westus"
  body = {
    kind = "StorageV2"
    properties = {
      allowBlobPublicAccess = true // This is not allowed by the policy
    }
    sku = {
      name = "Premium_LRS"
    }
  }
}
```

When you run `terraform plan`, you will see an error message like this:

```shell
╷
│ Error: Preflight Validation: Invalid configuration
│ 
│   with azapi_resource.storageaccount,
│   on main.tf line 11, in resource "azapi_resource" "storageaccount":
│   11: resource "azapi_resource" "storageaccount" {
│ 
│ POST https://management.azure.com/providers/Microsoft.Resources/validateResources
│ --------------------------------------------------------------------------------
│ RESPONSE 400: 400 Bad Request
│ ERROR CODE: ResourceValidationFailed
│ --------------------------------------------------------------------------------
│ {
│   "error": {
│     "code": "ResourceValidationFailed",
│     "message": "Resource validation failed, correlation id: 'd6f63f55-5c98-494c-f4ad-ad5ef3c73482', see details for more information.",
│     "details": [
│       {
│         "code": "RequestDisallowedByPolicy",
│         "target": "acctestsa1227",
│         "message": "Resource 'acctestsa1227' was disallowed by policy. Policy identifiers: '[{\"policyAssignment\":{\"name\":\"Deny Storage Account Creation with Anonymous Access\",\"id\":\"/subscriptions/000000/providers/Microsoft.Authorization/policyAssignments/5712913d870246df83e718b5\"},\"policyDefinition\":{\"name\":\"Deny Storage Account Creation with Anonymous Access\",\"id\":\"/subscriptions/000000/providers/Microsoft.Authorization/policyDefinitions/example\"}}]'.",
│         "additionalInfo": [
│           // Policy details
│         ]
│       }
│     ]
│   }
│ }
│ --------------------------------------------------------------------------------
│ 
```

### Unknown Values

During a `terraform plan`, some values may not yet be known — for example, when a field references an attribute of a resource that has not been created yet. The preflight validation handles unknown values differently depending on where they appear in the resource body.

**Unknown values inside the `properties` object** are replaced with a placeholder expression (`[length('foo')]`) before the request is sent to the preflight API. This allows validation to proceed even when some property values are not yet known.

**Unknown values outside the `properties` object** (for example, in top-level fields such as `kind`, `sku`, or `location`) cannot be safely substituted without risking false validation results. When the provider detects an unknown value in any top-level field other than `properties`, it skips preflight validation entirely for that resource and emits a warning instead.

For example, the following configuration references `azapi_resource.resourceGroup.id` as the value of `kind`, which is an unknown value outside the `properties` bag at plan time. Preflight validation will be skipped for this resource:

```hcl
resource "azapi_resource" "storageAccount" {
  type      = "Microsoft.Storage/storageAccounts@2023-05-01"
  parent_id = azapi_resource.resourceGroup.id
  name      = azapi_resource.resourceGroup.id
  location  = "westus"
  body = {
    kind = azapi_resource.resourceGroup.id  # unknown value outside properties — validation skipped
    properties = {
      accessTier = "Cold"
    }
    sku = {
      name = "Standard_LRS"
    }
  }
}
```

In contrast, an unknown value inside `properties` (such as `accessTier = azapi_resource.resourceGroup.id`) will be replaced with a placeholder and validation will still run:

```hcl
resource "azapi_resource" "storageAccount" {
  type      = "Microsoft.Storage/storageAccounts@2023-05-01"
  parent_id = azapi_resource.resourceGroup.id
  name      = azapi_resource.resourceGroup.id
  location  = "westus"
  body = {
    kind = "StorageV2"
    properties = {
      accessTier = azapi_resource.resourceGroup.id  # unknown value inside properties — replaced with placeholder
    }
    sku = {
      name = "Standard_LRS"
    }
  }
}
```

To determine if preflight validation was skipped due to unknown values, look for a warning message in the Terraform plan output similar to the following:

```shell
[WARN]  azapi: Skipping preflight validation for resource Microsoft.Storage/storageAccounts@2023-05-01 because the body is invalid: unknown value found outside the properties bag
```
