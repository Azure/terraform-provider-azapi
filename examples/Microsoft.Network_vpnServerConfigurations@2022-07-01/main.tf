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

variable "radius_server_secret" {
  type        = string
  description = "The RADIUS server secret for VPN server configuration"
  sensitive   = true
}

resource "azapi_resource" "resourceGroup" {
  type     = "Microsoft.Resources/resourceGroups@2020-06-01"
  name     = var.resource_name
  location = var.location
}

resource "azapi_resource" "vpnServerConfiguration" {
  type      = "Microsoft.Network/vpnServerConfigurations@2022-07-01"
  parent_id = azapi_resource.resourceGroup.id
  name      = var.resource_name
  location  = var.location
  body = {
    properties = {
      radiusClientRootCertificates = [
      ]
      radiusServerAddress = ""
      radiusServerRootCertificates = [
      ]
      radiusServerSecret = var.radius_server_secret
      radiusServers = [
        {
          radiusServerAddress = "10.105.1.1"
          radiusServerScore   = 15
          radiusServerSecret  = var.radius_server_secret
        },
      ]
      vpnAuthenticationTypes = [
        "Radius",
      ]
      vpnClientIpsecPolicies = [
      ]
      vpnClientRevokedCertificates = [
      ]
      vpnClientRootCertificates = [
      ]
      vpnProtocols = [
        "OpenVPN",
        "IkeV2",
      ]
    }
  }
  schema_validation_enabled = false
  response_export_values    = ["*"]
}

