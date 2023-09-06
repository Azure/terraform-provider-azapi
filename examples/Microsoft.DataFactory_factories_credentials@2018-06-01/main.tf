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

resource "azapi_resource" "factory" {
  type      = "Microsoft.DataFactory/factories@2018-06-01"
  parent_id = azapi_resource.resourceGroup.id
  name      = var.resource_name
  location  = var.location

  identity {
    type         = "UserAssigned"
    identity_ids = [azapi_resource.userAssignedIdentity.id]
  }
  body = jsonencode({
    properties = {
      publicNetworkAccess = "Enabled"
      repoConfiguration   = null
    }
  })
}

resource "azapi_resource" "userAssignedIdentity" {
  type      = "Microsoft.ManagedIdentity/userAssignedIdentities@2023-01-31"
  parent_id = azapi_resource.resourceGroup.id
  name      = var.resource_name
  location  = azapi_resource.resourceGroup.location
}

resource "azapi_resource" "credential" {
  type      = "Microsoft.DataFactory/factories/credentials@2018-06-01"
  parent_id = azapi_resource.factory.id
  name      = var.resource_name
  body = jsonencode({
    properties = {
      type        = "ManagedIdentity"
      annotations = ["test"]
      description = "this is a test"
      typeProperties = {
        resourceId = azapi_resource.userAssignedIdentity.id
      }
    }

  })
  ignore_casing           = false
  ignore_missing_property = false
}
