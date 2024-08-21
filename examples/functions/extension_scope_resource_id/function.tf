locals {
  resource_id = "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/rg1/providers/Microsoft.Network/virtualNetworks/vnet1"
  extension_name = "Microsoft.Authorization/locks"
  resource_names = ["mylock"]
}

output "extension_resource_id" {
  value = provider::azapi::extension_resource_id(local.resource_id, local.extension_name, local.resource_names)
}