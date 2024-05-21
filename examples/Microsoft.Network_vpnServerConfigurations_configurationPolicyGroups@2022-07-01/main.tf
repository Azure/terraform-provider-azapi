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
      radiusServerSecret = ""
      radiusServers = [
        {
          radiusServerAddress = "10.105.1.1"
          radiusServerScore   = 15
          radiusServerSecret  = "vindicators-the-return-of-worldender"
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

resource "azapi_resource" "configurationPolicyGroup" {
  type      = "Microsoft.Network/vpnServerConfigurations/configurationPolicyGroups@2022-07-01"
  parent_id = azapi_resource.vpnServerConfiguration.id
  name      = var.resource_name
  body = {
    properties = {
      isDefault = false
      policyMembers = [
        {
          attributeType  = "RadiusAzureGroupId"
          attributeValue = "6ad1bd08"
          name           = "policy1"
        },
      ]
      priority = 0
    }
  }
  schema_validation_enabled = false
  response_export_values    = ["*"]
}

