locals {
  resource_id = "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/myResourceGroup/providers/Microsoft.Network/virtualNetworks/myVNet"
  resource_type = "Microsoft.Network/virtualNetworks"
}

output "parsed_resource_id" {
  value = provider::azapi::parse_resource_id(local.resource_type, local.resource_id)
}