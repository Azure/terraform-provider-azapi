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
  default = "westeurope"
}

resource "azapi_resource" "resourceGroup" {
  type     = "Microsoft.Resources/resourceGroups@2020-06-01"
  name     = var.resource_name
  location = var.location
}

resource "azapi_resource" "server" {
  type      = "Microsoft.DBforPostgreSQL/servers@2017-12-01"
  parent_id = azapi_resource.resourceGroup.id
  name      = var.resource_name
  location  = var.location
  body = {
    properties = {
      administratorLogin         = "acctestun"
      administratorLoginPassword = "H@Sh1CoR3!"
      createMode                 = "Default"
      infrastructureEncryption   = "Disabled"
      minimalTlsVersion          = "TLS1_2"
      publicNetworkAccess        = "Enabled"
      sslEnforcement             = "Enabled"
      storageProfile = {
        backupRetentionDays = 7
        storageAutogrow     = "Enabled"
        storageMB           = 51200
      }
      version = "9.6"
    }
    sku = {
      capacity = 2
      family   = "Gen5"
      name     = "GP_Gen5_2"
      tier     = "GeneralPurpose"
    }
  }
  schema_validation_enabled = false
  response_export_values    = ["*"]
}

resource "azapi_resource" "firewallRule" {
  type      = "Microsoft.DBforPostgreSQL/servers/firewallRules@2017-12-01"
  parent_id = azapi_resource.server.id
  name      = var.resource_name
  body = {
    properties = {
      endIpAddress   = "255.255.255.255"
      startIpAddress = "0.0.0.0"
    }
  }
  schema_validation_enabled = false
  response_export_values    = ["*"]
}

