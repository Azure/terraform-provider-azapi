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

variable "recommendation_id" {
  type = string
}

data "azapi_client_config" "current" {}

resource "azapi_resource" "suppression" {
  type      = "Microsoft.Advisor/recommendations/suppressions@2023-01-01"
  parent_id = "${data.azapi_client_config.current.subscription_resource_id}/providers/Microsoft.Advisor/recommendations/${var.recommendation_id}"
  name      = var.resource_name
  body = {
    properties = {
      suppressionId = ""
      ttl           = "00:30:00"
    }
  }
}
