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

resource "azapi_resource" "elasticSan" {
  type      = "Microsoft.ElasticSan/elasticSans@2023-01-01"
  parent_id = azapi_resource.resourceGroup.id
  name      = "${var.resource_name}-es"
  location  = var.location
  body = {
    properties = {
      baseSizeTiB             = 1
      extendedCapacitySizeTiB = 0
      sku = {
        name = "Premium_LRS"
        tier = "Premium"
      }
    }
  }
}

resource "azapi_resource" "volumeGroup" {
  type      = "Microsoft.ElasticSan/elasticSans/volumeGroups@2023-01-01"
  parent_id = azapi_resource.elasticSan.id
  name      = "${var.resource_name}-vg"
  body = {
    properties = {
      encryption = "EncryptionAtRestWithPlatformKey"
      networkAcls = {
        virtualNetworkRules = []
      }
      protocolType = "Iscsi"
    }
  }
}

resource "azapi_resource" "volume" {
  type      = "Microsoft.ElasticSan/elasticSans/volumeGroups/volumes@2023-01-01"
  parent_id = azapi_resource.volumeGroup.id
  name      = "${var.resource_name}-v"
  body = {
    properties = {
      sizeGiB = 1
    }
  }
}
