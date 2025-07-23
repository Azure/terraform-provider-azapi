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

variable "oauth_client_id" {
  type        = string
  description = "The OAuth client ID for the authorization server"
}

variable "oauth_client_secret" {
  type        = string
  description = "The OAuth client secret for the authorization server"
  sensitive   = true
}

resource "azapi_resource" "resourceGroup" {
  type     = "Microsoft.Resources/resourceGroups@2020-06-01"
  name     = var.resource_name
  location = var.location
}

resource "azapi_resource" "service" {
  type      = "Microsoft.ApiManagement/service@2021-08-01"
  parent_id = azapi_resource.resourceGroup.id
  name      = var.resource_name
  location  = var.location
  body = {
    properties = {
      certificates = [
      ]
      customProperties = {
        "Microsoft.WindowsAzure.ApiManagement.Gateway.Security.Backend.Protocols.Ssl30" = "false"
        "Microsoft.WindowsAzure.ApiManagement.Gateway.Security.Backend.Protocols.Tls10" = "false"
        "Microsoft.WindowsAzure.ApiManagement.Gateway.Security.Backend.Protocols.Tls11" = "false"
        "Microsoft.WindowsAzure.ApiManagement.Gateway.Security.Protocols.Tls10"         = "false"
        "Microsoft.WindowsAzure.ApiManagement.Gateway.Security.Protocols.Tls11"         = "false"
      }
      disableGateway      = false
      publicNetworkAccess = "Enabled"
      publisherEmail      = "pub1@email.com"
      publisherName       = "pub1"
      virtualNetworkType  = "None"
    }
    sku = {
      capacity = 0
      name     = "Consumption"
    }
  }
  schema_validation_enabled = false
  response_export_values    = ["*"]
}

resource "azapi_resource" "authorizationServer" {
  type      = "Microsoft.ApiManagement/service/authorizationServers@2021-08-01"
  parent_id = azapi_resource.service.id
  name      = var.resource_name
  body = {
    properties = {
      authorizationEndpoint = "https://azacceptance.hashicorptest.com/client/authorize"
      authorizationMethods = [
        "GET",
      ]
      clientAuthenticationMethod = [
      ]
      clientId                   = var.oauth_client_id
      clientRegistrationEndpoint = "https://azacceptance.hashicorptest.com/client/register"
      clientSecret               = var.oauth_client_secret
      defaultScope               = ""
      description                = ""
      displayName                = "Test Group"
      grantTypes = [
        "implicit",
      ]
      resourceOwnerPassword = ""
      resourceOwnerUsername = ""
      supportState          = false
      tokenBodyParameters = [
      ]
    }
  }
  schema_validation_enabled = false
  response_export_values    = ["*"]
}

