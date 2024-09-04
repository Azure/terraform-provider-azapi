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

resource "azapi_resource" "deploymentScript" {
  type      = "Microsoft.Resources/deploymentScripts@2020-10-01"
  parent_id = azapi_resource.resourceGroup.id
  name      = var.resource_name
  location  = var.location
  body = {
    kind = "AzurePowerShell"
    properties = {
      azPowerShellVersion  = "8.3"
      cleanupPreference    = "Always"
      environmentVariables = null
      retentionInterval    = "P1D"
      scriptContent        = "\t\t$output = 'Hello'\n\t\tWrite-Output $output\n\t\t$DeploymentScriptOutputs = @{}\n\t\t$DeploymentScriptOutputs['text'] = $output\n"
      supportingScriptUris = null
      timeout              = "P1D"
    }
  }
  schema_validation_enabled = false
  response_export_values    = ["*"]
}

