locals {
  subscription_id     = "00000000-0000-0000-0000-000000000000"
  resource_group_name = "rg1"
  resource_type       = "Microsoft.Network/virtualNetworks/subnets"
  resource_names      = ["vnet1", "subnet1"]
}

// it will output "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/rg1/providers/Microsoft.Network/virtualNetworks/vnet1/subnets/subnet1"
output "resource_group_resource_id" {
  value = provider::azapi::resource_group_resource_id(local.subscription_id, local.resource_group_name, local.resource_type, local.resource_names)
}