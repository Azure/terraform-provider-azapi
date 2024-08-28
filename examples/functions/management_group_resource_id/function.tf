locals {
  management_group_name = "mg1"
  resource_type         = "Microsoft.Billing/billingAccounts/billingProfiles"
  resource_names        = ["ba1", "bp1"]
}

// it will output "/providers/Microsoft.Management/managementGroups/mg1/providers/Microsoft.Billing/billingAccounts/ba1/billingProfiles/bp1"
output "management_group_resource_id" {
  value = provider::azapi::management_group_resource_id(local.management_group_name, local.resource_type, local.resource_names)
}