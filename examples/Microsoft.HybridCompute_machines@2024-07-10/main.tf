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

resource "azapi_resource" "machine" {
  type      = "Microsoft.HybridCompute/machines@2024-07-10"
  parent_id = azapi_resource.resourceGroup.id
  name      = "${var.resource_name}hcm"
  location  = var.location
  body = {
    kind = "SCVMM"
  }
}

