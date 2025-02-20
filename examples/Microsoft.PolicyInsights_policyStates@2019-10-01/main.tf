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

resource "azapi_resource" "resourceGroup" {
  type     = "Microsoft.Resources/resourceGroups@2020-06-01"
  name     = var.resource_name
  location = var.location
}

resource "azapi_resource_action" "triggerEvaluation" {
  type        = "Microsoft.PolicyInsights/policyStates@2019-10-01"
  resource_id = "${azapi_resource.resourceGroup.id}/providers/Microsoft.PolicyInsights/policyStates/latest"
  action      = "triggerEvaluation"
}
