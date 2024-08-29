locals {
  resource_id   = "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/myResourceGroup/providers/Microsoft.Network/virtualNetworks/myVNet"
  resource_type = "Microsoft.Network/virtualNetworks"
}

// it will output below object
# {
#   "id" = "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/myResourceGroup/providers/Microsoft.Network/virtualNetworks/myVNet"
#   "name" = "myVNet"
#   "parent_id" = "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/myResourceGroup"
#   "parts" = tomap({
#     "providers" = "Microsoft.Network"
#     "resourceGroups" = "myResourceGroup"
#     "subscriptions" = "00000000-0000-0000-0000-000000000000"
#     "virtualNetworks" = "myVNet"
#   })
#   "provider_namespace" = "Microsoft.Network"
#   "resource_group_name" = "myResourceGroup"
#   "subscription_id" = "00000000-0000-0000-0000-000000000000"
#   "type" = "Microsoft.Network/virtualNetworks"
# }
output "parsed_resource_id" {
  value = provider::azapi::parse_resource_id(local.resource_type, local.resource_id)
}