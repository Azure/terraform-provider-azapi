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

resource "azapi_resource" "ApplicationGatewayWebApplicationFirewallPolicy" {
  type      = "Microsoft.Network/ApplicationGatewayWebApplicationFirewallPolicies@2022-07-01"
  parent_id = azapi_resource.resourceGroup.id
  name      = var.resource_name
  location  = var.location
  body = {
    properties = {
      customRules = [
      ]
      managedRules = {
        exclusions = [
        ]
        managedRuleSets = [
          {
            ruleGroupOverrides = [
            ]
            ruleSetType    = "OWASP"
            ruleSetVersion = "3.1"
          },
        ]
      }
      policySettings = {
        fileUploadLimitInMb    = 100
        maxRequestBodySizeInKb = 128
        mode                   = "Detection"
        requestBodyCheck       = true
        state                  = "Enabled"
      }
    }
  }
  schema_validation_enabled = false
  response_export_values    = ["*"]
}

