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
  default = "westus"
}

variable "administrator_login_password" {
  type        = string
  sensitive   = true
  description = "The administrator login password for the MySQL flexible server"
}

resource "azapi_resource" "resourceGroup" {
  type     = "Microsoft.Resources/resourceGroups@2020-06-01"
  name     = var.resource_name
  location = var.location
}

resource "azapi_resource" "flexibleServer" {
  type      = "Microsoft.DBforMySQL/flexibleServers@2023-12-30"
  parent_id = azapi_resource.resourceGroup.id
  name      = "${var.resource_name}-mysql"
  location  = var.location
  body = {
    properties = {
      administratorLogin         = "tfadmin"
      administratorLoginPassword = var.administrator_login_password
      backup = {
        backupRetentionDays = 7
        geoRedundantBackup  = "Disabled"
      }
      dataEncryption = {
        type = "SystemManaged"
      }
      highAvailability = {
        mode = "Disabled"
      }
      version = "8.0.21"
    }
    sku = {
      name = "Standard_B1ms"
      tier = "Burstable"
    }
  }
}

resource "azapi_update_resource" "configuration" {
  type      = "Microsoft.DBforMySQL/flexibleServers/configurations@2023-12-30"
  parent_id = azapi_resource.flexibleServer.id
  name      = "character_set_server"
  body = {
    properties = {
      value = "utf8mb4"
    }
  }
}

