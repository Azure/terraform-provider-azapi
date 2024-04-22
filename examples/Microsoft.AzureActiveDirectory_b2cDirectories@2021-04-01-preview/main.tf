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
  default = "acctest0003"
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

data "azapi_resource_id" "b2cDirectory" {
  type      = "Microsoft.AzureActiveDirectory/b2cDirectories@2021-04-01-preview"
  parent_id = azapi_resource.resourceGroup.id
  name      = "${var.resource_name}.onmicrosoft.com"
}

resource "azapi_resource_action" "b2cDirectory" {
  type        = "Microsoft.AzureActiveDirectory/b2cDirectories@2021-04-01-preview"
  resource_id = data.azapi_resource_id.b2cDirectory.id
  method      = "PUT"
  body = {
    location = "United States"
    properties = {
      createTenantProperties = {
        countryCode = "US"
        displayName = var.resource_name
      }
    }
    sku = {
      name = "PremiumP1"
      tier = "A0"
    }

  }
}
