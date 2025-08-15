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

resource "azapi_resource" "workspace" {
  type      = "Microsoft.OperationalInsights/workspaces@2022-10-01"
  parent_id = azapi_resource.resourceGroup.id
  name      = var.resource_name
  location  = var.location
  body = {
    properties = {
      sku = {
        name = "PerGB2018"
      }
    }
  }
}

resource "azapi_resource" "solution" {
  type      = "Microsoft.OperationsManagement/solutions@2015-11-01-preview"
  parent_id = azapi_resource.resourceGroup.id
  name      = "ContainerInsights(${var.resource_name})"
  location  = var.location
  body = {
    plan = {
      name          = "ContainerInsights(${var.resource_name})"
      product       = "OMSGallery/ContainerInsights"
      promotionCode = ""
      publisher     = "Microsoft"
    }
    properties = {
      workspaceResourceId = azapi_resource.workspace.id
    }
  }
  tags = {
    Environment = "Test"
  }
}
