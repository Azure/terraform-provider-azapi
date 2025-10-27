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
  default = "westus2"
}

resource "azapi_resource" "resourceGroup" {
  type     = "Microsoft.Resources/resourceGroups@2020-06-01"
  name     = var.resource_name
  location = var.location
}

resource "azapi_resource" "account" {
  type      = "Microsoft.CognitiveServices/accounts@2025-06-01"
  parent_id = azapi_resource.resourceGroup.id
  name      = var.resource_name
  location  = var.location
  identity {
    type = "SystemAssigned"
  }

  body = {
    kind = "AIServices"
    properties = {
      allowProjectManagement        = true
      customSubDomainName           = "cog-${var.resource_name}"
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

resource "azapi_resource" "project" {
  type      = "Microsoft.CognitiveServices/accounts/projects@2025-06-01"
  parent_id = azapi_resource.account.id
  name      = var.resource_name
  location  = var.location
  identity {
    type = "SystemAssigned"
  }

  body = {
    properties = {
      displayName = "project-${var.resource_name}"
      description = "test project"
    }
  }

  tags = {
    foo = "bar"
  }
  schema_validation_enabled = false
  response_export_values    = ["*"]
}
