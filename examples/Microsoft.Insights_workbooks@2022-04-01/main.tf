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

resource "azapi_resource" "workbook" {
  type      = "Microsoft.Insights/workbooks@2022-04-01"
  parent_id = azapi_resource.resourceGroup.id
  name      = "be1ad266-d329-4454-b693-8287e4d3b35d"
  location  = var.location
  body = {
    kind = "shared"
    properties = {
      category       = "workbook"
      displayName    = "acctest-amw-230630032616547405"
      serializedData = "{\"fallbackResourceIds\":[\"Azure Monitor\"],\"isLocked\":false,\"items\":[{\"content\":{\"json\":\"Test2022\"},\"name\":\"text - 0\",\"type\":1}],\"version\":\"Notebook/1.0\"}"
      sourceId       = "azure monitor"
    }
  }
  schema_validation_enabled = false
  response_export_values    = ["*"]
}

