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

resource "azapi_resource" "mobileNetwork" {
  type      = "Microsoft.MobileNetwork/mobileNetworks@2022-11-01"
  parent_id = azapi_resource.resourceGroup.id
  name      = var.resource_name
  location  = var.location
  body = {
    properties = {
      publicLandMobileNetworkIdentifier = {
        mcc = "001"
        mnc = "01"
      }
    }
  }
  schema_validation_enabled = false
  response_export_values    = ["*"]
}

resource "azapi_resource" "dataBoxEdgeDevice" {
  type      = "Microsoft.DataBoxEdge/dataBoxEdgeDevices@2022-03-01"
  parent_id = azapi_resource.resourceGroup.id
  name      = var.resource_name
  location  = var.location
  body = {
    sku = {
      name = "EdgeP_Base"
      tier = "Standard"
    }
  }
  schema_validation_enabled = false
  response_export_values    = ["*"]
}

resource "azapi_resource" "site" {
  type      = "Microsoft.MobileNetwork/mobileNetworks/sites@2022-11-01"
  parent_id = azapi_resource.mobileNetwork.id
  name      = var.resource_name
  location  = var.location
  body = {
    properties = {
    }
  }
  schema_validation_enabled = false
  response_export_values    = ["*"]
}

resource "azapi_resource" "packetCoreControlPlane" {
  type      = "Microsoft.MobileNetwork/packetCoreControlPlanes@2022-11-01"
  parent_id = azapi_resource.resourceGroup.id
  name      = var.resource_name
  location  = var.location
  body = {
    properties = {
      controlPlaneAccessInterface = {
      }
      localDiagnosticsAccess = {
        authenticationType = "AAD"
      }
      platform = {
        azureStackEdgeDevice = {
          id = azapi_resource.dataBoxEdgeDevice.id
        }
        type = "AKS-HCI"
      }
      sites = [
        {
          id = azapi_resource.site.id
        },
      ]
      sku   = "G0"
      ueMtu = 1440
    }
  }
  schema_validation_enabled = false
  response_export_values    = ["*"]
}

