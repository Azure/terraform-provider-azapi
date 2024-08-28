locals {
  resource_id    = "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/rg1/providers/Microsoft.Network/virtualNetworks/vnet1"
  resource_type  = "Microsoft.Authorization/locks"
  resource_names = ["mylock"]
}

// it will output "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/rg1/providers/Microsoft.Network/virtualNetworks/vnet1/providers/Microsoft.Authorization/locks/mylock"
output "extension_resource_id" {
  value = provider::azapi::extension_resource_id(local.resource_id, local.resource_type, local.resource_names)
}