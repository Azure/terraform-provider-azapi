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

resource "azapi_resource" "emailService" {
  type      = "Microsoft.Communication/emailServices@2023-04-01-preview"
  parent_id = azapi_resource.resourceGroup.id
  name      = var.resource_name
  location  = "global"
  body = {
    properties = {
      dataLocation = "United States"
    }
  }
  schema_validation_enabled = false
  response_export_values    = ["*"]
}

resource "azapi_resource" "domain" {
  type      = "Microsoft.Communication/emailServices/domains@2023-04-01-preview"
  name      = "example.com"
  location  = "global"
  parent_id = azapi_resource.emailService.id
  tags = {
    env = "Test"
  }
  body = {
    properties = {
      domainManagement       = "CustomerManaged"
      userEngagementTracking = "Disabled"
    }
  }
}

resource "azapi_resource" "senderUsername" {
  type      = "Microsoft.Communication/emailServices/domains/senderUsernames@2023-04-01-preview"
  name      = "TestSenderUserName"
  parent_id = azapi_resource.domain.id
  body = {
    properties = {
      displayName = "TestDisplayName"
      username    = "TestSenderUserName"
    }
  }
}
