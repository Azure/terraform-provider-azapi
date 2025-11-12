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
  description = "The administrator login password for the PostgreSQL flexible server"
}

resource "azapi_resource" "resourceGroup" {
  type     = "Microsoft.Resources/resourceGroups@2020-06-01"
  name     = var.resource_name
  location = var.location
}

resource "azapi_resource" "flexibleServer" {
  type      = "Microsoft.DBforPostgreSQL/flexibleServers@2024-08-01"
  parent_id = azapi_resource.resourceGroup.id
  name      = "${var.resource_name}-primary"
  location  = var.location
  body = {
    properties = {
      administratorLogin         = "psqladmin"
      administratorLoginPassword = var.administrator_login_password
      availabilityZone           = "1"
      backup = {
        geoRedundantBackup = "Disabled"
      }
      highAvailability = {
        mode = "Disabled"
      }
      network = {
        publicNetworkAccess = "Disabled"
      }
      storage = {
        autoGrow      = "Disabled"
        storageSizeGB = 32
        tier          = "P30"
      }
      version = "16"
    }
    sku = {
      name = "Standard_D2ads_v5"
      tier = "GeneralPurpose"
    }
  }
}

resource "azapi_resource" "flexibleServer_1" {
  type      = "Microsoft.DBforPostgreSQL/flexibleServers@2024-08-01"
  parent_id = azapi_resource.resourceGroup.id
  name      = "${var.resource_name}-replica"
  location  = var.location
  body = {
    properties = {
      availabilityZone = "1"
      backup = {
        geoRedundantBackup = "Disabled"
      }
      createMode = "Replica"
      highAvailability = {
        mode = "Disabled"
      }
      network = {
        publicNetworkAccess = "Disabled"
      }
      sourceServerResourceId = azapi_resource.flexibleServer.id
      storage = {
        autoGrow      = "Disabled"
        storageSizeGB = 32
        tier          = "P30"
      }
      version = "16"
    }
  }
}

resource "azapi_resource" "virtualEndpoint" {
  type      = "Microsoft.DBforPostgreSQL/flexibleServers/virtualEndpoints@2024-08-01"
  parent_id = azapi_resource.flexibleServer.id
  name      = var.resource_name
  body = {
    properties = {
      endpointType = "ReadWrite"
      members      = [azapi_resource.flexibleServer_1.name]
    }
  }
}
