terraform {
  required_version = ">= 1.14.0"

  required_providers {
    azapi = {
      source = "Azure/azapi"
    }
    azurerm = {
      source = "hashicorp/azurerm"
    }
  }
}

provider "azapi" {
}

provider "azurerm" {
  features {}
}

data "azurerm_client_config" "current" {}

resource "azurerm_resource_group" "example" {
  name     = "example-foundry-toolbox"
  location = "westus3"
}

resource "azurerm_cognitive_account" "example" {
  name                       = "example-foundry-toolbox"
  location                   = azurerm_resource_group.example.location
  resource_group_name        = azurerm_resource_group.example.name
  kind                       = "AIServices"
  sku_name                   = "S0"
  project_management_enabled = true
  custom_subdomain_name      = "example-foundry-toolbox"
  local_auth_enabled         = false

  identity {
    type = "SystemAssigned"
  }
}

resource "azurerm_cognitive_account_project" "example" {
  name                 = "example-project"
  cognitive_account_id = azurerm_cognitive_account.example.id
  location             = azurerm_resource_group.example.location
  display_name         = "Toolbox example"
  description          = "Foundry project for the toolbox example"

  identity {
    type = "SystemAssigned"
  }
}

resource "azurerm_role_assignment" "foundry_user" {
  scope              = azurerm_cognitive_account_project.example.id
  role_definition_id = "/subscriptions/${data.azurerm_client_config.current.subscription_id}/providers/Microsoft.Authorization/roleDefinitions/53ca6127-db72-4b80-b1b0-d745d6d5456d"
  principal_id       = data.azurerm_client_config.current.object_id
}

resource "azapi_data_plane_resource" "toolbox" {
  type      = "Microsoft.Foundry/toolboxes@v1"
  parent_id = azurerm_cognitive_account_project.example.endpoints["AI Foundry API"]
  name      = "example-toolbox"

  depends_on = [azurerm_role_assignment.foundry_user]

  body = {
    description = "Toolbox managed by Terraform"
    tools = [
      {
        type        = "web_search"
        description = "Search the web for current information"
      }
    ]
  }

  response_export_values = ["versions.latest.version"]

  lifecycle {
    action_trigger {
      events  = [after_create, after_update]
      actions = [action.azapi_data_plane_action.activate_toolbox]
    }
  }
}

action "azapi_data_plane_action" "activate_toolbox" {
  config {
    type        = "Microsoft.Foundry/toolboxes@v1"
    resource_id = azapi_data_plane_resource.toolbox.id
    method      = "PATCH"

    body = {
      default_version = azapi_data_plane_resource.toolbox.output.versions.latest.version
    }
  }
}
