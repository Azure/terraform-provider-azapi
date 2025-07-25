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

variable "postgresql_administrator_password" {
  type        = string
  description = "The administrator password for the PostgreSQL flexible server"
  sensitive   = true
}

resource "azapi_resource" "resourceGroup" {
  type     = "Microsoft.Resources/resourceGroups@2020-06-01"
  name     = var.resource_name
  location = var.location
}

resource "azapi_resource" "flexibleServer" {
  type      = "Microsoft.DBforPostgreSQL/flexibleServers@2023-06-01-preview"
  parent_id = azapi_resource.resourceGroup.id
  name      = var.resource_name
  location  = var.location
  body = {
    identity = {
      type                   = "None",
      userAssignedIdentities = null
    },
    properties = {
      administratorLogin         = "adminTerraform"
      administratorLoginPassword = var.postgresql_administrator_password
      availabilityZone           = "2"
      backup = {
        geoRedundantBackup = "Disabled"
      }
      highAvailability = {
        mode = "Disabled"
      }
      network = {
      }
      storage = {
        storageSizeGB = 32
      }
      version = "12"
    }
    sku = {
      name = "Standard_D2s_v3"
      tier = "GeneralPurpose"
    }
  }
  schema_validation_enabled = false
  response_export_values    = ["*"]
}

resource "azapi_update_resource" "pgbouncerEnabled" {
  type      = "Microsoft.DBforPostgreSQL/flexibleServers/configurations@2022-12-01"
  parent_id = azapi_resource.flexibleServer.id
  name      = "pgbouncer.enabled"
  body = {
    properties = {
      value  = "true"
      source = "user-override"
    }
  }
}

resource "azapi_update_resource" "pgbouncerDefaultPoolSize" {
  depends_on = [azapi_update_resource.pgbouncerEnabled]
  type       = "Microsoft.DBforPostgreSQL/flexibleServers/configurations@2022-12-01"
  parent_id  = azapi_resource.flexibleServer.id
  name       = "pgbouncer.default_pool_size"
  body = {
    properties = {
      value  = "40"
      source = "user-override"
    }
  }
}
