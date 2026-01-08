---
layout: "azapi"
page_title: "Feature: Generate Configuration with List Resources"
description: |-
  This guide demonstrates how to use the list resource protocol to discover and generate Terraform configurations from existing Azure resources.

---

# Generate Configuration with List Resources

The `list` block in Terraform allows you to discover and generate configurations from existing Azure resources. This is particularly useful for importing existing infrastructure into Terraform or for understanding the current state of your resources.

## Overview

The list resource protocol enables you to:
- Discover existing Azure resources in your subscriptions and resource groups
- Automatically generate Terraform configuration code from discovered resources
- Import existing infrastructure into Terraform management

## List Resources of a Specific Type

To list all resources of a specific type in a resource group and generate their configurations:

### Example 1: List Storage Accounts in a Resource Group

```hcl
list "azapi_resource" "storage_accounts" {
  provider = azapi
  config {
    type      = "Microsoft.Storage/storageAccounts@2021-04-01"
    parent_id = "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/myResourceGroup"
  }
}
```

This configuration discovers all storage accounts within the specified resource group. The `type` parameter specifies the Azure resource type and API version, while `parent_id` identifies the resource group containing the storage accounts.

### Example 2: List Key Vaults in a Resource Group

```hcl
list "azapi_resource" "key_vaults" {
  provider = azapi
  config {
    type      = "Microsoft.KeyVault/vaults@2023-02-01"
    parent_id = "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/myResourceGroup"
  }
}
```

This configuration discovers all Key Vaults within the specified resource group. Replace the subscription ID and resource group name with your actual values.

### Example 3: List Subnets in a Virtual Network

```hcl
list "azapi_resource" "subnets" {
  provider = azapi
  config {
    type      = "Microsoft.Network/virtualNetworks/subnets@2023-04-01"
    parent_id = "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/myResourceGroup/providers/Microsoft.Network/virtualNetworks/myVnet"
  }
}
```

This configuration discovers all subnets within a specific virtual network. Note that for child resources like subnets, the `parent_id` must be the full resource ID of the parent virtual network, not just the resource group.

### Example 4: List Diagnostic Settings for a Resource

```hcl
list "azapi_resource" "diagnostic_settings" {
  provider = azapi
  config {
    type      = "Microsoft.Insights/diagnosticSettings@2021-05-01-preview"
    parent_id = "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/myResourceGroup/providers/Microsoft.Storage/storageAccounts/mystorageaccount"
  }
}
```

This configuration discovers all diagnostic settings for a specific storage account. Diagnostic settings are extension resources that can be attached to other Azure resources, so the `parent_id` must be the complete resource ID of the target resource.

## List All Resources in a Resource Group

To list all resources in a resource group regardless of their type:

```hcl
list "azapi_resource" "all_resources" {
  provider = azapi
  config {
    parent_id = "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/myResourceGroup"
  }
}
```

When the `type` attribute is omitted, the list block will enumerate all resources in the resource group using the Azure Resource Manager API endpoint: `/subscriptions/{sub}/resourceGroups/{rg}/resources`.

**Important Notes:**
- The `parent_id` must be a resource group ID when `type` is omitted

## Using Terraform Query Commands

Once you have defined a `list` block, you can use Terraform's query commands to interact with the discovered resources.

### Query Resources

To see a list of all discovered resources:

```bash
terraform query
```

This command will display all resources discovered by the `list` blocks in your configuration, showing their IDs, types, and other basic information.

**Example output:**
```
list.azapi_resource.storage_accounts   id=/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/myResourceGroup/providers/Microsoft.Storage/storageAccounts/mystorageaccount,type=Microsoft.Storage/storageAccounts@2021-04-01                Microsoft.Storage/storageAccounts@2021-04-01 - mystorageaccount
list.azapi_resource.storage_accounts   id=/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/myResourceGroup/providers/Microsoft.Storage/storageAccounts/anotherstorage,type=Microsoft.Storage/storageAccounts@2021-04-01                  Microsoft.Storage/storageAccounts@2021-04-01 - anotherstorage
```

### Generate Configuration Files

To automatically generate Terraform configuration files from the discovered resources:

```bash
terraform query -generate-config-out=./generated.tf
```

This command will:
1. Query all resources discovered by the `list` blocks
2. Generate complete Terraform configuration code for each resource
3. Write the generated configurations to `./generated.tf`

The generated configuration will include all resource attributes and can be used as a starting point for managing these resources with Terraform.

**Example generated configuration:**
```hcl
import {
  to       = azapi_resource.mystorageaccount
  provider = azapi
  identity = {
    id   = "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/myResourceGroup/providers/Microsoft.Storage/storageAccounts/mystorageaccount"
    type = "Microsoft.Storage/storageAccounts@2021-04-01"
  }
}

resource "azapi_resource" "mystorageaccount" {
  type      = "Microsoft.Storage/storageAccounts@2021-04-01"
  parent_id = "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/myResourceGroup"
  name      = "mystorageaccount"
  location  = "eastus"

  body = {
    kind = "StorageV2"
    sku = {
      name = "Standard_LRS"
    }
    properties = {
      accessTier               = "Hot"
      allowBlobPublicAccess    = false
      minimumTlsVersion        = "TLS1_2"
      supportsHttpsTrafficOnly = true
      // ... other properties
    }
  }
}
```

## Complete Workflow Example

Here's a complete workflow for importing existing Azure infrastructure:

### Step 1: Create List Configuration

Create a file named `main.tfquery.hcl` to discover resources:

```hcl
# main.tfquery.hcl

list "azapi_resource" "production_vms" {
  provider = azapi
  config {
    type      = "Microsoft.Compute/virtualMachines@2023-03-01"
    parent_id = "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/production-rg"
  }
}

list "azapi_resource" "production_storage" {
  provider = azapi
  config {
    type      = "Microsoft.Storage/storageAccounts@2021-04-01"
    parent_id = "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/production-rg"
  }
}

list "azapi_resource" "production_databases" {
  provider = azapi
  config {
    type      = "Microsoft.Sql/servers/databases@2021-11-01"
    parent_id = "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/production-rg/providers/Microsoft.Sql/servers/myserver"
  }
}
```

### Step 2: Initialize Terraform

```bash
terraform init
```

### Step 3: Query Resources

View the discovered resources:

```bash
terraform query
```

Review the output to verify that all expected resources are discovered.

### Step 4: Generate Configuration

Generate Terraform configuration files:

```bash
terraform query -generate-config-out=./imported-resources.tf
```

### Step 5: Review Generated Configuration

Open `imported-resources.tf` and review the generated code:

- Verify all resource attributes are correct
- Remove any sensitive data or default values
- Adjust resource names to follow your naming conventions
- Add comments and documentation
- Organize resources logically

### Step 6: Apply Configuration

After reviewing and customizing the generated configuration, run a plan to see what will be imported:

```bash
terraform plan
```

You will see output indicating the resources to be imported:

```
Plan: 16 to import, 0 to add, 0 to change, 0 to destroy.
```

Then apply the configuration to complete the import:

```bash
terraform apply
```
