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

resource "azapi_resource" "component" {
  type      = "Microsoft.Insights/components@2020-02-02"
  parent_id = azapi_resource.resourceGroup.id
  name      = var.resource_name
  location  = var.location
  body = {
    kind = "web"
    properties = {
      Application_Type                = "web"
      DisableIpMasking                = false
      DisableLocalAuth                = false
      ForceCustomerStorageForProfiler = false
      RetentionInDays                 = 90
      SamplingPercentage              = 100
      publicNetworkAccessForIngestion = "Enabled"
      publicNetworkAccessForQuery     = "Enabled"
    }
  }
  schema_validation_enabled = false
  response_export_values    = ["*"]
}

data "azapi_resource_id" "analyticsItem" {
  type      = "microsoft.insights/components/analyticsItems@2015-05-01"
  parent_id = azapi_resource.component.id
  name      = "item"
}

resource "azapi_resource_action" "analyticsItem" {
  type        = "microsoft.insights/components/analyticsItems@2015-05-01"
  resource_id = data.azapi_resource_id.analyticsItem.id
  method      = "PUT"
  body = {
    Content = "requests #test"
    Name    = "testquery"
    Scope   = "shared"
    Type    = "query"
  }
}

