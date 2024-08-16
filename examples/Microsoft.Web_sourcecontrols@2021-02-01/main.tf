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
  default = "eastus"
}

resource "azapi_update_resource" "sourcecontrol" {
  type      = "Microsoft.Web/sourcecontrols@2021-02-01"
  parent_id = "/"
  name      = "GitHub"
  body = {
    properties = {
      token       = "abcdefghijklmnopqrstuvwxyz"
      tokenSecret = "abcdefghijklmnopqrstuvwxyz"
    }
  }
  response_export_values = ["*"]
}

