locals {
  resource_type  = "Microsoft.Billing/billingAccounts/billingProfiles"
  resource_names = ["ba1", "bp1"]
}

// it will output "/providers/Microsoft.Billing/billingAccounts/ba1/billingProfiles/bp1"
output "tenant_resource_id" {
  value = provider::azapi::tenant_resource_id(local.resource_type, local.resource_names)
}