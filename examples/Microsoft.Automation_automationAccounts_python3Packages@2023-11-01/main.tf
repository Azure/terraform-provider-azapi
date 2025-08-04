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

resource "azapi_resource" "automationAccount" {
  type      = "Microsoft.Automation/automationAccounts@2023-11-01"
  parent_id = azapi_resource.resourceGroup.id
  name      = var.resource_name
  location  = var.location
  body = {
    properties = {
      disableLocalAuth = false
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

resource "azapi_resource" "python3Package" {
  type      = "Microsoft.Automation/automationAccounts/python3Packages@2023-11-01"
  parent_id = azapi_resource.automationAccount.id
  name      = var.resource_name
  body = {
    properties = {
      contentLink = {
        uri     = "https://files.pythonhosted.org/packages/py3/r/requests/requests-2.31.0-py3-none-any.whl"
        version = "2.31.0"
      }
    }
  }
  tags = {
    key = "foo"
  }
}

