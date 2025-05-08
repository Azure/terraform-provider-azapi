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

variable "pat" {
  type        = string
  sensitive   = true
  description = "GitHub Personal Access Token"
}

resource "azapi_resource" "resourceGroup" {
  type     = "Microsoft.Resources/resourceGroups@2020-06-01"
  name     = var.resource_name
  location = var.location
}

resource "azapi_resource" "automationAccount" {
  type      = "Microsoft.Automation/automationAccounts@2023-11-01"
  parent_id = azapi_resource.resourceGroup.id
  name      = var.resource_name
  location  = var.location
  identity {
    type         = "SystemAssigned"
    identity_ids = []
  }

  body = {
    properties = {
      encryption = {
        keySource = "Microsoft.Automation"
      }
      publicNetworkAccess = true
      sku = {
        name = "Basic"
      }
    }
  }
}

resource "azapi_resource" "sourceControl" {
  type      = "Microsoft.Automation/automationAccounts/sourceControls@2023-11-01"
  name      = var.resource_name
  parent_id = azapi_resource.automationAccount.id

  body = {
    properties = {
      repoUrl        = "https://github.com/Azure-Samples/acr-build-helloworld-node.git"
      branch         = "master"
      sourceType     = "GitHub"
      folderPath     = "/"
      autoSync       = false
      publishRunbook = false

      securityToken = {
        tokenType   = "PersonalAccessToken"
        accessToken = var.pat
      }
    }
  }
}
