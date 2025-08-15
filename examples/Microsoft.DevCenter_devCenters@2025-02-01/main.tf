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

resource "azapi_resource" "resourceGroup" {
  type     = "Microsoft.Resources/resourceGroups@2020-06-01"
  name     = var.resource_name
  location = var.location
}

resource "azapi_resource" "userAssignedIdentity" {
  type      = "Microsoft.ManagedIdentity/userAssignedIdentities@2023-01-31"
  parent_id = azapi_resource.resourceGroup.id
  name      = "${var.resource_name}-uai"
  location  = var.location
  body      = {}
}

resource "azapi_resource" "devCenter" {
  type      = "Microsoft.DevCenter/devCenters@2025-02-01"
  parent_id = azapi_resource.resourceGroup.id
  name      = var.resource_name
  location  = var.location
  identity {
    type         = "SystemAssigned, UserAssigned"
    identity_ids = [azapi_resource.userAssignedIdentity.id]
  }
  body = {}
  tags = {
    environment = "terraform-acctests"
    some_key    = "some-value"
  }
}

