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

resource "azapi_resource" "flexibleServer" {
  type      = "Microsoft.DBforMySQL/flexibleServers@2021-05-01"
  parent_id = azapi_resource.resourceGroup.id
  name      = var.resource_name
  location  = var.location
  body = {
    properties = {
      administratorLogin         = "adminTerraform"
      administratorLoginPassword = "QAZwsx123"
      backup = {
        backupRetentionDays = 7
        geoRedundantBackup  = "Disabled"
      }
      createMode = ""
      dataEncryption = {
        type = "SystemManaged"
      }
      highAvailability = {
        mode = "Disabled"
      }
      network = {
      }
      version = ""
    }
    sku = {
      name = "Standard_B1s"
      tier = "Burstable"
    }
  }
  schema_validation_enabled = false
  response_export_values    = ["*"]
}

resource "azapi_resource" "database" {
  type      = "Microsoft.DBforMySQL/flexibleServers/databases@2021-05-01"
  parent_id = azapi_resource.flexibleServer.id
  name      = var.resource_name
  body = {
    properties = {
      charset   = "utf8"
      collation = "utf8_unicode_ci"
    }
  }
  schema_validation_enabled = false
  response_export_values    = ["*"]
}

