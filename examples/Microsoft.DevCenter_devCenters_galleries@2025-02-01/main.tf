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

resource "azapi_resource" "resourceGroup" {
  type     = "Microsoft.Resources/resourceGroups@2020-06-01"
  name     = var.resource_name
  location = var.location
}

resource "azapi_resource" "gallery" {
  type      = "Microsoft.Compute/galleries@2022-03-03"
  parent_id = azapi_resource.resourceGroup.id
  name      = "${var.resource_name}sig"
  location  = var.location
  body = {
    properties = {
      description = ""
    }
  }
}

resource "azapi_resource" "userAssignedIdentity" {
  type      = "Microsoft.ManagedIdentity/userAssignedIdentities@2023-01-31"
  parent_id = azapi_resource.resourceGroup.id
  name      = "${var.resource_name}uai"
  location  = var.location
}

resource "azapi_resource" "devCenter" {
  type      = "Microsoft.DevCenter/devCenters@2025-02-01"
  parent_id = azapi_resource.resourceGroup.id
  name      = "${var.resource_name}dc"
  location  = var.location
  identity {
    type         = "UserAssigned"
    identity_ids = [azapi_resource.userAssignedIdentity.id]
  }
  body = {
    properties = {}
  }
}

# Assign Reader role on the gallery to the Dev Center's user-assigned identity
locals {
  reader_role_id = "acdd72a7-3385-48ef-bd42-f606fba81ae7"
}

resource "azapi_resource" "roleAssignment" {
  type      = "Microsoft.Authorization/roleAssignments@2022-04-01"
  name      = uuidv5("url", "${azapi_resource.gallery.id}/roleAssignments/${azapi_resource.userAssignedIdentity.output.properties.principalId}")
  parent_id = azapi_resource.gallery.id
  body = {
    properties = {
      roleDefinitionId = "/subscriptions/${data.azapi_client_config.current.subscription_id}/providers/Microsoft.Authorization/roleDefinitions/${local.reader_role_id}"
      principalId      = azapi_resource.userAssignedIdentity.output.properties.principalId
      principalType    = "ServicePrincipal"
    }
  }
  depends_on = [azapi_resource.userAssignedIdentity, azapi_resource.gallery]
}

resource "azapi_resource" "gallery_1" {
  type      = "Microsoft.DevCenter/devCenters/galleries@2025-02-01"
  parent_id = azapi_resource.devCenter.id
  name      = "${var.resource_name}dcg"
  body = {
    properties = {
      galleryResourceId = azapi_resource.gallery.id
    }
  }
  depends_on = [azapi_resource.roleAssignment]
}
