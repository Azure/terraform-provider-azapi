terraform {
  required_providers {
    azapi = {
      source = "Azure/azapi"
    }
  }
}

variable "resource_name" {
  type    = string
  default = "acctest0001"
}

variable "location" {
  type    = string
  default = "westus"
}

data "azapi_client_config" "current" {}

resource "azapi_resource" "resourceGroup" {
  type     = "Microsoft.Resources/resourceGroups@2020-06-01"
  name     = var.resource_name
  location = var.location
}

resource "azapi_resource" "service" {
  type      = "Microsoft.ApiManagement/service@2022-08-01"
  parent_id = azapi_resource.resourceGroup.id
  name      = "${var.resource_name}-am"
  location  = var.location
  body = {
    properties = {
      certificates = []
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
}

resource "azapi_resource" "api" {
  type      = "Microsoft.ApiManagement/service/apis@2022-08-01"
  parent_id = azapi_resource.service.id
  name      = "${var.resource_name}-api;rev=1"
  body = {
    properties = {
      apiRevisionDescription = ""
      apiType                = "http"
      apiVersionDescription  = ""
      authenticationSettings = {}
      description            = "What is my purpose? You parse butter."
      displayName            = "Butter Parser"
      path                   = "butter-parser"
      protocols              = ["http", "https"]
      serviceUrl             = "https://example.com/foo/bar"
      subscriptionKeyParameterNames = {
        header = "X-Butter-Robot-API-Key"
        query  = "location"
      }
      subscriptionRequired = true
      type                 = "http"
    }
  }
}

resource "azapi_resource" "operation" {
  type      = "Microsoft.ApiManagement/service/apis/operations@2022-08-01"
  parent_id = azapi_resource.api.id
  name      = "${var.resource_name}-operation"
  body = {
    properties = {
      description        = ""
      displayName        = "DELETE Resource"
      method             = "DELETE"
      responses          = []
      templateParameters = []
      urlTemplate        = "/resource"
    }
  }
}
