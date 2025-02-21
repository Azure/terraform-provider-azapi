terraform {
  required_providers {
    azapi = {
      source = "Azure/azapi"
    }
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
  default = "eastus"
}

resource "azapi_resource" "resourceGroup" {
  type     = "Microsoft.Resources/resourceGroups@2020-06-01"
  name     = var.resource_name
  location = var.location
}

resource "azapi_resource" "server" {
  type      = "Microsoft.Sql/servers@2022-05-01-preview"
  parent_id = azapi_resource.resourceGroup.id
  name      = var.resource_name
  location  = var.location
  body = {
    properties = {
      administratorLogin            = "mradministrator"
      administratorLoginPassword    = "thisIsDog11"
      minimalTlsVersion             = "1.2"
      publicNetworkAccess           = "Enabled"
      restrictOutboundNetworkAccess = "Disabled"
      version                       = "12.0"
    }
  }
  schema_validation_enabled = false
  response_export_values    = ["*"]
}

resource "azapi_update_resource" "auditingSettings" {
  type      = "Microsoft.Sql/servers/auditingSettings@2022-05-01-preview"
  name      = "default"
  parent_id = azapi_resource.server.id
  body = {
    properties = {
      auditActionsAndGroups = [
        "FAILED_DATABASE_AUTHENTICATION_GROUP",
        "SUCCESSFUL_DATABASE_AUTHENTICATION_GROUP"
      ]
      isAzureMonitorTargetEnabled = true
      state                       = "Enabled"
    }
  }
}
