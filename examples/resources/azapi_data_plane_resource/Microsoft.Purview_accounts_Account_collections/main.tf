terraform {
  required_providers {
    azapi = {
      source = "Azure/azapi"
    }
  }
}

provider "azapi" {
}

resource "azapi_resource" "resourceGroup" {
  type     = "Microsoft.Resources/resourceGroups@2021-04-01"
  name     = "example-resources"
  location = "westeurope"
}

resource "azapi_resource" "account" {
  type      = "Microsoft.Purview/accounts@2021-12-01"
  parent_id = azapi_resource.resourceGroup.id
  name      = "examplepurview"
  location  = azapi_resource.resourceGroup.location
  identity {
    type         = "SystemAssigned"
    identity_ids = []
  }
  body = {
    properties = {
      publicNetworkAccess = "Enabled"
    }
  }
  response_export_values = ["*"]
}

resource "azapi_data_plane_resource" "example" {
  type      = "Microsoft.Purview/accounts/Account/collections@2019-11-01-preview"
  parent_id = "${azapi_resource.account.name}.purview.azure.com"
  name      = "defaultResourceSetRuleConfig"
  body = {
    friendlyName = "Finance"
  }
}
