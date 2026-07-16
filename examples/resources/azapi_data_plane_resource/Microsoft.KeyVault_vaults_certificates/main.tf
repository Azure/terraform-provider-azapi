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

resource "azapi_resource_action" "add_accesspolicy_certificate" {
  type        = "Microsoft.KeyVault/vaults/accessPolicies@2023-02-01"
  resource_id = "${azapi_resource.vault.id}/accessPolicies/add"
  method      = "PUT"
  body = {
    properties = {
      accessPolicies = [{
        tenantId = data.azapi_client_config.current.tenant_id
        objectId = data.azapi_client_config.current.object_id
        permissions = {
          certificates = [
            "Get", "List", "Update", "Create", "Import", "Delete", "Recover", "Backup", "Restore", "ManageContacts", "ManageIssuers", "GetIssuers", "ListIssuers", "SetIssuers", "DeleteIssuers", "Purge"
          ]
        }
      }]
    }
  }
}

resource "azapi_data_plane_resource" "certificate" {
  type      = "Microsoft.KeyVault/vaults/certificates@7.4"
  parent_id = trimsuffix(trimprefix(azapi_resource.vault.output.vaultUri, "https://"), "/")
  name      = "examplecertificate"
  body = {
    policy = {
      issuer = {
        name = "Self"
      }
      key_props = {
        exportable = true
        kty        = "RSA"
        key_size   = 2048
        reuse_key  = false
      }
      secret_props = {
        contentType = "application/x-pkcs12"
      }
      x509_props = {
        subject         = "CN=contoso.com"
        validity_months = 12
      }
    }
  }

  depends_on = [
    azapi_resource_action.add_accesspolicy_certificate
  ]
}
