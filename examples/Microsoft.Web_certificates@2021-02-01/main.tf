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

resource "azapi_resource" "certificate" {
  type      = "Microsoft.Web/certificates@2021-02-01"
  name      = var.resource_name
  parent_id = azapi_resource.resourceGroup.id
  location  = var.location
  body = {
    properties = {
      pfxBlob  = filebase64("testdata/app_service_certificate.pfx")
      password = "terraform"
    }
  }
}
