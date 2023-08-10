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
  type                      = "Microsoft.Resources/resourceGroups@2020-06-01"
  name                      = var.resource_name
  location                  = var.location
}

resource "azapi_resource" "resourceProvider" {
  type      = "Microsoft.CustomProviders/resourceProviders@2018-09-01-preview"
  parent_id = azapi_resource.resourceGroup.id
  name      = var.resource_name
  location  = var.location
  body = jsonencode({
    properties = {
      resourceTypes = [
        {
          endpoint    = "https://example.com/"
          name        = "dEf1"
          routingType = "Proxy"
        },
      ]
    }
  })
  schema_validation_enabled = false
  response_export_values    = ["*"]
}

