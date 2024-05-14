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

resource "azapi_resource" "virtualWan" {
  type      = "Microsoft.Network/virtualWans@2022-07-01"
  parent_id = azapi_resource.resourceGroup.id
  name      = var.resource_name
  location  = var.location
  body = {
    properties = {
      allowBranchToBranchTraffic     = true
      disableVpnEncryption           = false
      office365LocalBreakoutCategory = "None"
      type                           = "Standard"
    }
  }
  schema_validation_enabled = false
  response_export_values    = ["*"]
}

resource "azapi_resource" "vpnSite" {
  type      = "Microsoft.Network/vpnSites@2022-07-01"
  parent_id = azapi_resource.resourceGroup.id
  name      = var.resource_name
  location  = var.location
  body = {
    properties = {
      addressSpace = {
        addressPrefixes = [
          "10.0.1.0/24",
        ]
      }
      virtualWan = {
        id = azapi_resource.virtualWan.id
      }
      vpnSiteLinks = [
        {
          name = "link1"
          properties = {
            fqdn      = ""
            ipAddress = "10.0.1.1"
            linkProperties = {
              linkProviderName = ""
              linkSpeedInMbps  = 0
            }
          }
        },
        {
          name = "link2"
          properties = {
            fqdn      = ""
            ipAddress = "10.0.1.2"
            linkProperties = {
              linkProviderName = ""
              linkSpeedInMbps  = 0
            }
          }
        },
      ]
    }
  }
  schema_validation_enabled = false
  response_export_values    = ["*"]
}

