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
  default = "eastus"
}

resource "azapi_resource" "resourceGroup" {
  type     = "Microsoft.Resources/resourceGroups@2020-06-01"
  name     = var.resource_name
  location = var.location
}

resource "azapi_resource" "account" {
  type      = "Microsoft.CognitiveServices/accounts@2022-10-01"
  parent_id = azapi_resource.resourceGroup.id
  name      = var.resource_name
  location  = var.location

  body = {
    identity = {
      type                   = "None"
      userAssignedIdentities = null
    }
    kind = "OpenAI"
    properties = {
      disableLocalAuth              = false
      dynamicThrottlingEnabled      = false
      publicNetworkAccess           = "Enabled"
      restrictOutboundNetworkAccess = false
    }
    sku = {
      name = "S0"
    }
  }
  schema_validation_enabled = false
  response_export_values    = ["*"]
}

resource "azapi_resource" "deployment" {
  type      = "Microsoft.CognitiveServices/accounts/deployments@2023-05-01"
  name      = "testdep"
  parent_id = azapi_resource.account.id
  body = {
    properties = {
      model = {
        format = "OpenAI"
        name   = "text-embedding-ada-002"
      }
    }
  }
}
