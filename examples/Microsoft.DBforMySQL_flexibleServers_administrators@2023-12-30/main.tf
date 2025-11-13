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

data "azapi_client_config" "current" {}

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

resource "azapi_resource" "userAssignedIdentity" {
  type      = "Microsoft.ManagedIdentity/userAssignedIdentities@2023-01-31"
  parent_id = azapi_resource.resourceGroup.id
  name      = "${var.resource_name}-uai1"
  location  = var.location
}

resource "azapi_resource" "userAssignedIdentity_1" {
  type      = "Microsoft.ManagedIdentity/userAssignedIdentities@2023-01-31"
  parent_id = azapi_resource.resourceGroup.id
  name      = "${var.resource_name}-uai2"
  location  = var.location
}

resource "azapi_resource" "flexibleServer" {
  type      = "Microsoft.DBforMySQL/flexibleServers@2023-12-30"
  parent_id = azapi_resource.resourceGroup.id
  name      = "${var.resource_name}-mysql"
  location  = var.location
  identity {
    type         = "UserAssigned"
    identity_ids = [azapi_resource.userAssignedIdentity.id]
  }
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

resource "azapi_resource" "administrator" {
  type      = "Microsoft.DBforMySQL/flexibleServers/administrators@2023-12-30"
  parent_id = azapi_resource.flexibleServer.id
  name      = "ActiveDirectory"
  body = {
    properties = {
      administratorType  = "ActiveDirectory"
      identityResourceId = azapi_resource.userAssignedIdentity.id
      login              = "sqladmin"
      sid                = data.azapi_client_config.current.object_id
      tenantId           = data.azapi_client_config.current.tenant_id
    }
  }
}

