terraform {
  required_providers {
    azapi = {
      source = "Azure/azapi"
    }
    azurerm = {
      source = "hashicorp/azurerm"
    }
  }
}

provider "azurerm" {
  features {
  }
}

provider "azapi" {
  skip_provider_registration = false
}

variable "resource_name" {
  type    = string
  default = "acctest0001"
}

variable "location" {
  type    = string
  default = "westeurope"
}

variable "ledger_certificate" {
  type        = string
  description = "The PEM-encoded certificate for the confidential ledger administrator"
  sensitive   = true
}

data "azurerm_client_config" "current" {
}

resource "azapi_resource" "resourceGroup" {
  type     = "Microsoft.Resources/resourceGroups@2020-06-01"
  name     = var.resource_name
  location = var.location
}

resource "azapi_resource" "ledger" {
  type      = "Microsoft.ConfidentialLedger/ledgers@2022-05-13"
  parent_id = azapi_resource.resourceGroup.id
  name      = var.resource_name
  location  = var.location
  body = {
    properties = {
      aadBasedSecurityPrincipals = [
        {
          ledgerRoleName = "Administrator"
          principalId    = data.azurerm_client_config.current.object_id
          tenantId       = data.azurerm_client_config.current.tenant_id
        },
      ]
      certBasedSecurityPrincipals = [
        {
          cert           = var.ledger_certificate
          ledgerRoleName = "Administrator"
        },
      ]
      ledgerType = "Private"
    }
  }
  schema_validation_enabled = false
  response_export_values    = ["*"]
}

