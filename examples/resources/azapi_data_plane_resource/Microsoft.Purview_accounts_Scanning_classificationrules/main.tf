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
  type      = "Microsoft.Purview/accounts/Scanning/classificationrules@2022-07-01-preview"
  parent_id = replace(azapi_resource.account.output.properties.endpoints.scan, "https://", "")
  name      = "exampleclassificationrule"
  body = {
    kind = "Custom"
    properties = {
      description        = "Custom classification rule for bank account numbers"
      classificationName = "MICROSOFT.FINANCIAL.AUSTRALIA.BANK_ACCOUNT_NUMBER"
      columnPatterns = [
        {
          pattern = "^data$"
          kind    = "Regex"
        }
      ]
      dataPatterns = [
        {
          pattern = "^[0-9]{2}-[0-9]{4}-[0-9]{6}-[0-9]{3}$"
          kind    = "Regex"
        }
      ]
      minimumPercentageMatch = 60
      ruleStatus             = "Enabled"
    }
  }
}
