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

variable "github_token" {
  type        = string
  description = "The GitHub access token for source control integration"
  sensitive   = true
}

variable "github_token_secret" {
  type        = string
  description = "The GitHub token secret for source control integration"
  sensitive   = true
}

resource "azapi_update_resource" "sourcecontrol" {
  type      = "Microsoft.Web/sourcecontrols@2021-02-01"
  parent_id = "/"
  name      = "GitHub"
  body = {
    properties = {
      token       = var.github_token
      tokenSecret = var.github_token_secret
    }
  }
  response_export_values = ["*"]
}

