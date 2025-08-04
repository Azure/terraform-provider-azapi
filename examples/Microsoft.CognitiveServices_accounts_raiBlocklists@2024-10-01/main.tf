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

resource "azapi_resource" "account" {
  type      = "Microsoft.CognitiveServices/accounts@2024-10-01"
  parent_id = azapi_resource.resourceGroup.id
  name      = "${var.resource_name}-ca"
  location  = var.location
  body = {
    kind = "OpenAI"
    properties = {
      allowedFqdnList               = []
      apiProperties                 = {}
      customSubDomainName           = ""
      disableLocalAuth              = false
      dynamicThrottlingEnabled      = false
      publicNetworkAccess           = "Enabled"
      restrictOutboundNetworkAccess = false
    }
    sku = {
      name = "S0"
    }
  }
}

resource "azapi_resource" "raiBlocklist" {
  type      = "Microsoft.CognitiveServices/accounts/raiBlocklists@2024-10-01"
  parent_id = azapi_resource.account.id
  name      = "${var.resource_name}-crb"
  body = {
    properties = {
      description = "Acceptance test data new azurerm resource"
    }
  }
}

