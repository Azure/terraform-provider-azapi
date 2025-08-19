---
layout: "azapi"
page_title: "AzAPI Provider: Choosing the Right Resource Type"
description: |-
  This guide provides clear instructions on when to use specific AzAPI resource types like azapi_resource, azapi_resource_action, and azapi_update_resource.
---

# Choosing the Right AzAPI Resource Type

The AzAPI provider offers several resource types that serve different purposes in managing Azure resources. Understanding when to use each type is crucial for effective infrastructure management. This guide provides clear instructions on when to use `azapi_resource`, `azapi_resource_action`, and `azapi_update_resource`.

## Overview of AzAPI Resource Types

| Resource Type | Purpose | Lifecycle Management | Use Case |
|---------------|---------|---------------------|----------|
| `azapi_resource` | Full resource lifecycle management | Complete (Create, Read, Update, Delete) | Primary resource management |
| `azapi_update_resource` | Modify properties of existing resources | Update-only (no delete) | Adding/modifying properties without full control |
| `azapi_resource_action` | Perform actions on existing resources | Action-based (no lifecycle) | Operations like start/stop, reset, etc. |

## When to Use `azapi_resource`

Use `azapi_resource` when you need **complete lifecycle management** of an Azure resource.

### Best for:
- Creating new Azure resources from scratch
- Managing the full configuration of a resource
- Resources that need to be created, updated, and destroyed as part of your Terraform lifecycle
- Resources not yet supported by the AzureRM provider
- Resources requiring cutting-edge API features

### Example Scenarios:
```terraform
# Creating a new Azure Container Registry
resource "azapi_resource" "container_registry" {
  type      = "Microsoft.ContainerRegistry/registries@2023-01-01-preview"
  name      = "myregistry"
  parent_id = azurerm_resource_group.example.id
  location  = "East US"

  body = {
    properties = {
      adminUserEnabled = true
      sku = {
        name = "Premium"
      }
    }
  }
}

# Managing a resource with preview features
resource "azapi_resource" "cognitive_service" {
  type      = "Microsoft.CognitiveServices/accounts@2023-10-01-preview"
  name      = "mycognitiveservice"
  parent_id = azurerm_resource_group.example.id
  location  = "East US"

  body = {
    properties = {
      # Preview feature not yet in AzureRM
      networkAcls = {
        defaultAction = "Deny"
        ipRules = [
          {
            value = "203.0.113.0/24"
          }
        ]
      }
    }
  }
}
```

## When to Use `azapi_update_resource`

Use `azapi_update_resource` when you need to **modify specific properties** of an existing resource without taking full control of its lifecycle.

### Best for:
- Adding tags to existing resources
- Modifying specific properties that aren't managed elsewhere
- Configuring properties that can only be set after resource creation
- Working alongside AzureRM provider resources
- Properties that require multi-step configuration processes

### Example Scenarios:
```terraform
# Add custom tags to an existing resource
resource "azapi_update_resource" "add_tags" {
  type        = "Microsoft.Storage/storageAccounts@2023-01-01"
  resource_id = azurerm_storage_account.example.id

  body = {
    tags = {
      Environment = "Production"
      CostCenter  = "Engineering"
      CustomTag   = "CustomValue"
    }
  }
}

# Configure Customer Managed Keys (CMK) after resource creation
resource "azapi_update_resource" "enable_cmk" {
  type        = "Microsoft.ServiceBus/namespaces@2021-11-01"
  resource_id = azurerm_servicebus_namespace.example.id

  body = {
    properties = {
      encryption = {
        keySource = "Microsoft.KeyVault"
        keyVaultProperties = [
          {
            keyName     = azurerm_key_vault_key.example.name
            keyVaultUri = azurerm_key_vault.example.vault_uri
          }
        ]
      }
    }
  }

  depends_on = [
    azurerm_servicebus_namespace.example,
    azurerm_key_vault_key.example
  ]
}

# Update diagnostic settings
resource "azapi_update_resource" "diagnostic_settings" {
  type        = "Microsoft.Network/networkSecurityGroups@2023-02-01"
  resource_id = azurerm_network_security_group.example.id

  body = {
    properties = {
      diagnosticsSettings = {
        logs = [
          {
            category = "NetworkSecurityGroupEvent"
            enabled  = true
          }
        ]
      }
    }
  }
}
```

### Key Characteristics:
- **Partial management**: Only manages the properties you specify
- **No deletion impact**: When deleted, the properties remain unchanged on the resource
- **Complementary**: Works well alongside other Terraform resources

### When NOT to use `azapi_update_resource`:
- You need full control over the resource lifecycle
- You're creating a new resource
- You need to perform actions rather than modify properties

## When to Use `azapi_resource_action`

Use `azapi_resource_action` when you need to **perform operations** on existing Azure resources without managing their lifecycle.

### Best for:
- Start/stop operations on virtual machines, app services, etc.
- Restart or reset operations
- Triggering specific actions like backups, rotations, or refreshes
- One-time operations that modify resource state
- Operations that should be repeated when Terraform runs

### Example Scenarios:
```terraform
# Start/stop a virtual machine based on a variable
resource "azapi_resource_action" "vm_start_stop" {
  type        = "Microsoft.Compute/virtualMachines@2023-03-01"
  resource_id = azurerm_linux_virtual_machine.example.id
  action      = var.vm_should_be_running ? "start" : "deallocate"
  method      = "POST"

  depends_on = [azurerm_linux_virtual_machine.example]
}

# Restart an app service when configuration changes
resource "azapi_resource_action" "app_service_restart" {
  type        = "Microsoft.Web/sites@2022-03-01"
  resource_id = azurerm_linux_web_app.example.id
  action      = "restart"
  method      = "POST"

  # Trigger restart when app settings change
  response_export_values = ["*"]

  depends_on = [azurerm_linux_web_app.example]
}

# Perform a backup operation
resource "azapi_resource_action" "storage_backup" {
  type        = "Microsoft.Storage/storageAccounts@2023-01-01"
  resource_id = azurerm_storage_account.example.id
  action      = "backup"
  method      = "POST"

  body = {
    backupType = "full"
  }
}
```

### Key Characteristics:
- **No resource deletion**: When you delete an `azapi_resource_action`, it doesn't affect the target Azure resource
- **Repeatable**: Actions are performed each time Terraform applies (unless `when` is specified)
- **Stateless**: No state is maintained about the resource itself, only about the action

### When NOT to use `azapi_resource_action`:
- You need to create or delete a resource
- You need to manage the full configuration of a resource
- You need to track and maintain specific property values


## Common Patterns and Best Practices

### Pattern 1: Multi-Step Resource Configuration
When resources require configuration that can only be applied after creation:

```terraform
# 1. Create the Service Bus namespace with system-assigned identity
resource "azapi_resource" "servicebus_namespace" {
  type      = "Microsoft.ServiceBus/namespaces@2021-11-01"
  name      = "my-servicebus-namespace"
  parent_id = azurerm_resource_group.example.id
  location  = "East US"

  body = {
    sku = {
      name = "Premium"
      tier = "Premium"
    }
    identity = {
      type = "SystemAssigned"
    }
    properties = {
      # Initial configuration without CMK
    }
  }
}

# 2. Create Key Vault and key (depends on the namespace identity)
resource "azurerm_key_vault" "example" {
  name                = "my-keyvault"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  tenant_id           = data.azurerm_client_config.current.tenant_id
  sku_name            = "premium"

  access_policy {
    tenant_id = data.azurerm_client_config.current.tenant_id
    object_id = azapi_resource.servicebus_namespace.identity[0].principal_id

    key_permissions = [
      "Get", "WrapKey", "UnwrapKey"
    ]
  }
}

resource "azurerm_key_vault_key" "example" {
  name         = "servicebus-cmk-key"
  key_vault_id = azurerm_key_vault.example.id
  key_type     = "RSA"
  key_size     = 2048

  key_opts = [
    "decrypt", "encrypt", "sign", "unwrapKey", "verify", "wrapKey"
  ]

  depends_on = [azapi_resource.servicebus_namespace]
}

# 3. Configure Customer Managed Keys (can only be done after namespace and key exist)
resource "azapi_update_resource" "enable_cmk" {
  type        = "Microsoft.ServiceBus/namespaces@2021-11-01"
  resource_id = azapi_resource.servicebus_namespace.id

  body = {
    properties = {
      encryption = {
        keySource = "Microsoft.KeyVault"
        keyVaultProperties = [
          {
            keyName     = azurerm_key_vault_key.example.name
            keyVaultUri = azurerm_key_vault.example.vault_uri
            keyVersion  = azurerm_key_vault_key.example.version
          }
        ]
      }
    }
  }

  depends_on = [
    azapi_resource.servicebus_namespace,
    azurerm_key_vault_key.example
  ]
}
```

### Pattern 2: Conditional Operations
Use `azapi_resource_action` for conditional operations:

```terraform
resource "azapi_resource_action" "conditional_restart" {
  count = var.restart_required ? 1 : 0

  type        = "Microsoft.Web/sites@2022-03-01"
  resource_id = azurerm_linux_web_app.example.id
  action      = "restart"
  method      = "POST"
}
```

### Pattern 3: Working with AzureRM Provider
Enhance AzureRM resources with additional capabilities:

```terraform
# Create with AzureRM for standard properties
resource "azurerm_storage_account" "example" {
  name                     = "mystorageaccount"
  resource_group_name      = azurerm_resource_group.example.name
  location                 = azurerm_resource_group.example.location
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

# Add preview features with azapi_update_resource
resource "azapi_update_resource" "storage_enhanced" {
  type        = "Microsoft.Storage/storageAccounts@2023-01-01"
  resource_id = azurerm_storage_account.example.id

  body = {
    properties = {
      # Preview feature not yet in AzureRM
      dnsEndpointType = "AzureDnsZone"
    }
  }
}
```

## Summary

Choose your AzAPI resource type based on your specific needs:

- **`azapi_resource`**: Full lifecycle management of Azure resources
- **`azapi_update_resource`**: Modifying specific properties of existing resources
- **`azapi_resource_action`**: Performing operations on existing resources

Each type serves a distinct purpose in the Azure infrastructure management ecosystem, and understanding these differences will help you build more maintainable and effective Terraform configurations.
