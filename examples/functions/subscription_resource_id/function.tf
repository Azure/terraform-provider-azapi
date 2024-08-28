locals {
  subscription_id = "00000000-0000-0000-0000-000000000000"
  resource_type   = "Microsoft.Resources/resourceGroups"
  resource_names  = ["rg1"]
}

// it will output "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/rg1"
output "subscription_resource_id" {
  value = provider::azapi::subscription_resource_id(local.subscription_id, local.resource_type, local.resource_names)
}