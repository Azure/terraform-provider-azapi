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

resource "azapi_resource" "hostPool" {
  type      = "Microsoft.DesktopVirtualization/hostPools@2023-09-05"
  parent_id = azapi_resource.resourceGroup.id
  name      = var.resource_name
  location  = var.location
  body = {
    properties = {
      hostPoolType          = "Pooled"
      loadBalancerType      = "BreadthFirst"
      maxSessionLimit       = 999999
      preferredAppGroupType = "Desktop"
      publicNetworkAccess   = "Enabled"
      startVMOnConnect      = false
      validationEnvironment = false
    }
  }
  schema_validation_enabled = false
  response_export_values    = ["*"]
}

resource "azapi_resource" "applicationGroup" {
  type      = "Microsoft.DesktopVirtualization/applicationGroups@2023-09-05"
  parent_id = azapi_resource.resourceGroup.id
  name      = var.resource_name
  location  = var.location
  body = {
    properties = {
      applicationGroupType = "RemoteApp"
      hostPoolArmPath      = azapi_resource.hostPool.id
    }
  }
  schema_validation_enabled = false
  ignore_casing             = true
  response_export_values    = ["*"]
}

resource "azapi_resource" "application" {
  type      = "Microsoft.DesktopVirtualization/applicationGroups/applications@2023-09-05"
  parent_id = azapi_resource.applicationGroup.id
  name      = var.resource_name
  location  = var.location
  body = {
    properties = {
      commandLineSetting = "DoNotAllow"
      filePath           = "C:\\Program Files\\Google\\Chrome\\Application\\chrome.exe"
      showInPortal       = false
    }
  }
  schema_validation_enabled = false
  response_export_values    = ["*"]
}
