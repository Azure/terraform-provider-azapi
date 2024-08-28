locals {
  parent_id     = "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/myResourceGroup"
  resource_type = "Microsoft.Network/virtualNetworks"
  name          = "myVNet"
}

// it will output "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/myResourceGroup/providers/Microsoft.Network/virtualNetworks/myVNet"
output "resource_id" {
  value = provider::azapi::build_resource_id(local.parent_id, local.resource_type, local.name)
}