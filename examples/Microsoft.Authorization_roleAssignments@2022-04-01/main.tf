terraform {
  required_providers {
    azapi = {
      source = "Azure/azapi"
    }
    azurerm = {
      source = "hashicorp/azurerm"
    }
  }
}

provider "azurerm" {
  features {
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

data "azurerm_role_definition" "roleAcrpull" {
  name  = "AcrPull"
  scope = azapi_resource.resourceGroup.id
}

resource "azurerm_user_assigned_identity" "uai" {
  name                = "TestUAI"
  resource_group_name = azapi_resource.resourceGroup.name
  location            = azapi_resource.resourceGroup.location
}

resource "azapi_resource" "roleAssignments" {
  type      = "Microsoft.Authorization/roleAssignments@2022-04-01"
  name      = "6faae21a-0cd6-4536-8c23-a278823d12ed"
  parent_id = azapi_resource.resourceGroup.id
  body = {
    properties = {
      principalId      = azurerm_user_assigned_identity.uai.principal_id
      principalType    = "ServicePrincipal"
      roleDefinitionId = data.azurerm_role_definition.roleAcrpull.id
    }
  }
}
