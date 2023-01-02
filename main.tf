terraform {
  required_providers {
    azapi = {
      version = "~> 1.0.2"
      source  = "***REMOVED***/azure/azapi"
    }
  }
}

provider "azapi" {
  subscription_id = "***REMOVED***"
  client_id = "***REMOVED***"
  tenant_id = "***REMOVED***"
  oidc_token_file_path = "token.txt"
  use_oidc = true
  # oidc_authority_host = 
}

resource "azapi_resource" "example102" {
  type      = "Microsoft.Resources/resourceGroups@2021-04-01"
  name      = "registry102"
  parent_id = "/subscriptions/***REMOVED***"

  location = "West Europe"

  body = jsonencode({})
  tags = {
    "key" = "value3"
  }
}