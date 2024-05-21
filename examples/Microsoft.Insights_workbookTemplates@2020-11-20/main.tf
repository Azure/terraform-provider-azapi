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

resource "azapi_resource" "workbookTemplate" {
  type      = "Microsoft.Insights/workbookTemplates@2020-11-20"
  parent_id = azapi_resource.resourceGroup.id
  name      = var.resource_name
  location  = var.location
  body = {
    properties = {
      galleries = [
        {
          category     = "workbook"
          name         = "test"
          order        = 0
          resourceType = "Azure Monitor"
          type         = "workbook"
        },
      ]
      priority = 0
      templateData = {
        "$schema" = "https://github.com/Microsoft/Application-Insights-Workbooks/blob/master/schema/workbook.json"
        items = [
          {
            content = {
              json = "## New workbook\n---\n\nWelcome to your new workbook."
            }
            name = "text - 2"
            type = 1
          },
        ]
        styleSettings = {
        }
        version = "Notebook/1.0"
      }
    }

  }
  schema_validation_enabled = false
  response_export_values    = ["*"]
}

