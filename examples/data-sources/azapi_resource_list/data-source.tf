terraform {
  required_providers {
    azapi = {
      source = "Azure/azapi"
    }
  }
}

provider "azapi" {
}

data "azapi_resource_list" "listBySubscription" {
  type      = "Microsoft.Automation/automationAccounts@2021-06-22"
  parent_id = "/subscriptions/00000000-0000-0000-0000-000000000000"
  response_export_values = {
    "values" = "value[].{name: name, publicNetworkAccess: properties.publicNetworkAccess}"
    "names"  = "value[].name"
  }
}

data "azapi_resource_list" "listByResourceGroup" {
  type      = "Microsoft.Automation/automationAccounts@2021-06-22"
  parent_id = "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/rg1"
  response_export_values = {
    "names" = "value[].name"
  }
}

data "azapi_resource_list" "listSubnetsByVnet" {
  type                   = "Microsoft.Network/virtualNetworks/subnets@2021-02-01"
  parent_id              = "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/rg1/providers/Microsoft.Network/virtualNetworks/vnet1"
  response_export_values = ["*"]
}
