locals {
  resource_group_id = "/subscriptions/12345678-1234-1234-1234-123456789012/resourceGroups/my-rg"
  resource_type     = "Microsoft.Storage/storageAccounts"
}

// it will output "z22rah77jfqry"
output "unique_name" {
  value = provider::azapi::unique_string(local.resource_group_id, local.resource_type)
}