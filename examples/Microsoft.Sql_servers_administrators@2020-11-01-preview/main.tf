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

variable "administrator_login_password" {
  type        = string
  description = "The administrator login password for the SQL server"
  sensitive   = true
}

data "azurerm_client_config" "current" {
}

resource "azapi_resource" "resourceGroup" {
  type     = "Microsoft.Resources/resourceGroups@2020-06-01"
  name     = var.resource_name
  location = var.location
}

resource "azapi_resource" "server" {
  type      = "Microsoft.Sql/servers@2015-05-01-preview"
  parent_id = azapi_resource.resourceGroup.id
  name      = var.resource_name
  location  = var.location
  body = {
    properties = {
      administratorLogin         = "mradministrator"
      administratorLoginPassword = var.administrator_login_password
      version                    = "12.0"
    }
  }
  schema_validation_enabled = false
  response_export_values    = ["*"]
}

resource "azapi_resource" "administrator" {
  type      = "Microsoft.Sql/servers/administrators@2020-11-01-preview"
  parent_id = azapi_resource.server.id
  name      = "ActiveDirectory"
  body = {
    properties = {
      administratorType = "ActiveDirectory"
      login             = "sqladmin"
      sid               = data.azurerm_client_config.current.client_id
      tenantId          = data.azurerm_client_config.current.tenant_id
    }
  }
  schema_validation_enabled = false
  response_export_values    = ["*"]
}

