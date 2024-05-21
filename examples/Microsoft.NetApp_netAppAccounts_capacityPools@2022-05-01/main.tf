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
  default = "centralus"
}

resource "azapi_resource" "resourceGroup" {
  type     = "Microsoft.Resources/resourceGroups@2020-06-01"
  name     = var.resource_name
  location = var.location
}

resource "azapi_resource" "netAppAccount" {
  type      = "Microsoft.NetApp/netAppAccounts@2022-05-01"
  parent_id = azapi_resource.resourceGroup.id
  name      = var.resource_name
  location  = var.location
  body = {
    properties = {
      activeDirectories = [
      ]
    }
    tags = {
      SkipASMAzSecPack = "true"
    }
  }
  schema_validation_enabled = false
  response_export_values    = ["*"]
}

resource "azapi_resource" "capacityPool" {
  type      = "Microsoft.NetApp/netAppAccounts/capacityPools@2022-05-01"
  parent_id = azapi_resource.netAppAccount.id
  name      = var.resource_name
  location  = var.location
  body = {
    properties = {
      serviceLevel = "Standard"
      size         = 4.398046511104e+12
    }
    tags = {
      SkipASMAzSecPack = "true"
    }
  }
  schema_validation_enabled = false
  response_export_values    = ["*"]
}

