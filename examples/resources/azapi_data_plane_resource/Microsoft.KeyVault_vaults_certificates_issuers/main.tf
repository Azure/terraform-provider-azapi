terraform {
  required_providers {
    azapi = {
      source = "Azure/azapi"
    }
  }
}

provider "azapi" {
}

variable "resource_name" {
  type    = string
  default = "acctest0001"
}

variable "location" {
  type    = string
  default = "westeurope"
}

data "azapi_client_config" "current" {}

locals {
  key_vault_certificates_officer_role_id = "/subscriptions/${data.azapi_client_config.current.subscription_id}/providers/Microsoft.Authorization/roleDefinitions/a4417e6f-fecd-4de8-b567-7b0420556985"
}

resource "azapi_resource" "resourceGroup" {
  type     = "Microsoft.Resources/resourceGroups@2021-04-01"
  name     = var.resource_name
  location = var.location
}

resource "azapi_resource" "vault" {
  type      = "Microsoft.KeyVault/vaults@2026-02-01"
  parent_id = azapi_resource.resourceGroup.id
  name      = var.resource_name
  location  = azapi_resource.resourceGroup.location
  body = {
    properties = {
      sku = {
        family = "A"
        name   = "standard"
      }
      tenantId                  = data.azapi_client_config.current.tenant_id
      enableRbacAuthorization   = true
      enabledForDiskEncryption  = true
      softDeleteRetentionInDays = 7
      accessPolicies            = []
    }
  }
  response_export_values = {
    vaultUri = "properties.vaultUri"
  }
}

resource "azapi_resource" "keyVaultCertificatesOfficer" {
  type      = "Microsoft.Authorization/roleAssignments@2022-04-01"
  parent_id = azapi_resource.vault.id
  name      = uuidv5("url", "${azapi_resource.vault.id}/${data.azapi_client_config.current.object_id}/${local.key_vault_certificates_officer_role_id}")
  body = {
    properties = {
      principalId      = data.azapi_client_config.current.object_id
      roleDefinitionId = local.key_vault_certificates_officer_role_id
    }
  }
}

resource "azapi_data_plane_resource" "example" {
  type      = "Microsoft.KeyVault/vaults/certificates/issuers@7.5"
  parent_id = trimsuffix(trimprefix(azapi_resource.vault.output.vaultUri, "https://"), "/")
  name      = "${var.resource_name}-issuer"
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

  retry = {
    error_message_regex  = ["Forbidden", "Unauthorized", "authorization"]
    interval_seconds     = 10
    max_interval_seconds = 60
  }

  depends_on = [
    azapi_resource.keyVaultCertificatesOfficer
  ]
}
