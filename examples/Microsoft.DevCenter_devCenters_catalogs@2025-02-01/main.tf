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

resource "azapi_resource" "resourceGroup" {
  type     = "Microsoft.Resources/resourceGroups@2020-06-01"
  name     = var.resource_name
  location = var.location
}

resource "azapi_resource" "devCenter" {
  type      = "Microsoft.DevCenter/devCenters@2025-02-01"
  parent_id = azapi_resource.resourceGroup.id
  name      = "${substr(var.resource_name, 0, 22)}-dc"
  location  = var.location
  identity {
    type = "SystemAssigned"
  }
  body = {
    properties = {}
  }
}

resource "azapi_resource" "catalog" {
  type      = "Microsoft.DevCenter/devCenters/catalogs@2025-02-01"
  parent_id = azapi_resource.devCenter.id
  name      = "${substr(var.resource_name, 0, 17)}-catalog"
  body = {
    properties = {
      adoGit = {
        branch           = "main"
        path             = "/template"
        secretIdentifier = "https://amlim-kv.vault.azure.net/secrets/ado/6279752c2bdd4a38a3e79d958cc36a75"
        uri              = "https://amlim@dev.azure.com/amlim/testCatalog/_git/testCatalog"
      }
    }
  }
}

