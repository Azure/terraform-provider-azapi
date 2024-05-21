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
  type      = "Microsoft.Sql/servers@2015-05-01-preview"
  parent_id = azapi_resource.resourceGroup.id
  name      = var.resource_name
  location  = var.location
  body = {
    properties = {
      administratorLogin         = "acctestadmin"
      administratorLoginPassword = "t2RX8A76GrnE4EKC"
      version                    = "12.0"
    }
  }
  schema_validation_enabled = false
  response_export_values    = ["*"]
}

resource "azapi_update_resource" "securityAlertPolicy" {
  type      = "Microsoft.Sql/servers/securityAlertPolicies@2017-03-01-preview"
  parent_id = azapi_resource.server.id
  name      = "Default"
  body = {
    properties = {
      state = "Disabled"
    }
  }
  response_export_values = ["*"]
}

