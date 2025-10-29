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

resource "azapi_resource" "resourceGroup" {
  type     = "Microsoft.Resources/resourceGroups@2020-06-01"
  name     = var.resource_name
  location = var.location
}

resource "azapi_resource" "vault" {
  type      = "Microsoft.KeyVault/vaults@2023-02-01"
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
      enableSoftDelete          = true
      softDeleteRetentionInDays = 7
      enablePurgeProtection     = true
      accessPolicies            = []
    }
  }
  schema_validation_enabled = false
  lifecycle {
    ignore_changes = [body.properties.accessPolicies]
  }
  response_export_values = {
    "vaultUri" = "properties.vaultUri"
  }
}

resource "azapi_resource_action" "add_accesspolicy_secret" {
  type        = "Microsoft.KeyVault/vaults/accessPolicies@2023-02-01"
  resource_id = "${azapi_resource.vault.id}/accessPolicies/add"
  method      = "PUT"
  body = {
    properties = {
      accessPolicies = [{
        tenantId = data.azapi_client_config.current.tenant_id
        objectId = data.azapi_client_config.current.object_id
        permissions = {
          secrets = [
            "Get", "List", "Set", "Delete", "Recover", "Backup", "Restore", "Purge"
          ]
        }
      }]
    }
  }
}

resource "azapi_data_plane_resource" "secret" {
  type      = "Microsoft.KeyVault/vaults/secrets@7.4"
  parent_id = trimsuffix(trimprefix(azapi_resource.vault.output.vaultUri, "https://"), "/")
  name      = var.resource_name
  body = {
    value = "my-secret-value"
    attributes = {
      enabled = true
    }
  }

  depends_on = [
    azapi_resource_action.add_accesspolicy_secret
  ]
}
