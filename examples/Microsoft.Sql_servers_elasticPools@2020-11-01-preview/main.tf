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
      administratorLogin            = "4dm1n157r470r"
      administratorLoginPassword    = "4-v3ry-53cr37-p455w0rd"
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

resource "azapi_resource" "elasticPool" {
  type      = "Microsoft.Sql/servers/elasticPools@2020-11-01-preview"
  parent_id = azapi_resource.server.id
  name      = var.resource_name
  location  = var.location
  body = {
    properties = {
      maintenanceConfigurationId = data.azapi_resource_id.publicMaintenanceConfiguration.id
      maxSizeBytes               = 5.24288e+09
      perDatabaseSettings = {
        maxCapacity = 5
        minCapacity = 0
      }
      zoneRedundant = false
    }
    sku = {
      capacity = 50
      family   = ""
      name     = "BasicPool"
      tier     = "Basic"
    }
  }
  schema_validation_enabled = false
  response_export_values    = ["*"]
}

