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

data "azurerm_client_config" "current" {
}

resource "azapi_resource" "resourceGroup" {
  type     = "Microsoft.Resources/resourceGroups@2020-06-01"
  name     = var.resource_name
  location = var.location
}

resource "azapi_resource" "server" {
  type      = "Microsoft.Sql/servers@2021-02-01-preview"
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

data "azapi_resource_id" "publicMaintenanceConfiguration" {
  type      = "Microsoft.Maintenance/publicMaintenanceConfigurations@2023-04-01"
  parent_id = "/subscriptions/${data.azurerm_client_config.current.subscription_id}"
  name      = "SQL_Default"
}

resource "azapi_resource" "database" {
  type      = "Microsoft.Sql/servers/databases@2021-02-01-preview"
  parent_id = azapi_resource.server.id
  name      = var.resource_name
  location  = var.location
  body = {
    properties = {
      autoPauseDelay                   = 0
      createMode                       = "Default"
      elasticPoolId                    = ""
      highAvailabilityReplicaCount     = 0
      isLedgerOn                       = false
      licenseType                      = "LicenseIncluded"
      maintenanceConfigurationId       = data.azapi_resource_id.publicMaintenanceConfiguration.id
      minCapacity                      = 0
      readScale                        = "Disabled"
      requestedBackupStorageRedundancy = "Geo"
      zoneRedundant                    = false
    }
  }
  schema_validation_enabled = false
  response_export_values    = ["*"]
}

resource "azapi_update_resource" "securityAlertPolicy" {
  type      = "Microsoft.Sql/servers/databases/securityAlertPolicies@2020-11-01-preview"
  parent_id = azapi_resource.database.id
  name      = "default"
  body = {
    properties = {
      state = "Disabled"
    }
  }
  response_export_values = ["*"]
}

