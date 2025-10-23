terraform {
  required_providers {
    azapi = {
      source = "Azure/azapi"
    }
  }
}

provider "azapi" {
}

data "azapi_client_config" "current" {}

resource "azapi_resource" "resourceGroup" {
  type     = "Microsoft.Resources/resourceGroups@2021-04-01"
  name     = "example-resources"
  location = "westeurope"
}

resource "azapi_resource" "vault" {
  type      = "Microsoft.KeyVault/vaults@2023-02-01"
  parent_id = azapi_resource.resourceGroup.id
  name      = "examplekeyvault"
  location  = azapi_resource.resourceGroup.location
  body = {
    properties = {
      sku = {
        family = "A"
        name   = "standard"
      }
      tenantId                  = data.azapi_client_config.current.tenant_id
      enabledForDiskEncryption  = true
      softDeleteRetentionInDays = 7
      enablePurgeProtection     = false
      accessPolicies            = []
    }
  }
  schema_validation_enabled = false
  lifecycle {
    ignore_changes = [body.properties.accessPolicies]
  }
  response_export_values = {
    vaultUri = "properties.vaultUri"
  }
}

resource "azapi_resource_action" "add_accesspolicy" {
  type        = "Microsoft.KeyVault/vaults/accessPolicies@2023-02-01"
  resource_id = "${azapi_resource.vault.id}/accessPolicies/add"
  method      = "PUT"
  body = {
    properties = {
      accessPolicies = [{
        tenantId = data.azapi_client_config.current.tenant_id
        objectId = data.azapi_client_config.current.object_id
        permissions = {
          certificates = ["managecontacts", "getissuers", "setissuers", "deleteissuers"]
        }
      }]
    }
  }
}

resource "azapi_data_plane_resource" "example" {
  type      = "Microsoft.KeyVault/vaults/certificates/issuers@7.4"
  parent_id = trimsuffix(trimprefix(azapi_resource.vault.output.vaultUri, "https://"), "/")
  name      = "exampleissuer"
  body = {
    provider = "Test"
    credentials = {
      account_id = "keyvaultuser"
    }
    org_details = {
      admin_details = [
        {
          first_name = "John"
          last_name  = "Doe"
          email      = "admin@microsoft.com"
          phone      = "4255555555"
        }
      ]
    }
  }

  depends_on = [
    azapi_resource_action.add_accesspolicy
  ]
}
