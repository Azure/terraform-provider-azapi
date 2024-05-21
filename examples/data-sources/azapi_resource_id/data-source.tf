terraform {
  required_providers {
    azapi = {
      source = "Azure/azapi"
    }
  }
}

provider "azapi" {
}

data "azapi_resource_id" "account" {
  type        = "Microsoft.Automation/automationAccounts@2021-06-22"
  resource_id = "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/rg1/providers/Microsoft.Automation/automationAccounts/automationAccount1"
}

output "account_name" {
  value = data.azapi_resource_id.account.name
}

output "account_resource_group" {
  value = data.azapi_resource_id.account.resource_group_name
}

output "account_subscription" {
  value = data.azapi_resource_id.account.subscription_id
}

output "account_parent_id" {
  value = data.azapi_resource_id.account.parent_id
}

data "azapi_resource_id" "vnet" {
  type      = "Microsoft.Network/virtualNetworks@2021-02-01"
  parent_id = "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/rg1"
  name      = "vnet1"
}

output "vnet_id" {
  value = data.azapi_resource_id.vnet.id
}

output "vnet_resource_group" {
  value = data.azapi_resource_id.vnet.resource_group_name
}

output "vnet_subscription" {
  value = data.azapi_resource_id.vnet.subscription_id
}
