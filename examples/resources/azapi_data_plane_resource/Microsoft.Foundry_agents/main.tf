terraform {
  required_providers {
    azapi = {
      source = "Azure/azapi"
    }
    random = {
      source = "hashicorp/random"
    }
  }
}

provider "azapi" {
}

provider "random" {
}

data "azapi_client_config" "current" {}

resource "random_string" "unique" {
  length      = 5
  min_numeric = 5
  numeric     = true
  special     = false
  lower       = true
  upper       = false
}

variable "location" {
  type        = string
  description = "Azure region for the resources"
}

variable "resource_group_tags" {
  type        = map(string)
  description = "Optional tags to apply to the resource group (useful for policy-constrained subscriptions)."
  default     = {}
}

locals {
  foundry_project_user_role_definition_id = "/subscriptions/${data.azapi_client_config.current.subscription_id}/providers/Microsoft.Authorization/roleDefinitions/53ca6127-db72-4b80-b1b0-d745d6d5456d"
  foundry_project_user_role_assignment_id = format("%s-%s-%s-%s-%s",
    substr(md5("foundry-project-user-${random_string.unique.result}"), 0, 8),
    substr(md5("foundry-project-user-${random_string.unique.result}"), 8, 4),
    substr(md5("foundry-project-user-${random_string.unique.result}"), 12, 4),
    substr(md5("foundry-project-user-${random_string.unique.result}"), 16, 4),
    substr(md5("foundry-project-user-${random_string.unique.result}"), 20, 12)
  )
}

resource "azapi_resource" "resource_group" {
  type     = "Microsoft.Resources/resourceGroups@2021-04-01"
  name     = "acctest${random_string.unique.result}"
  location = var.location
  tags     = var.resource_group_tags
}

resource "azapi_resource" "foundry" {
  type                      = "Microsoft.CognitiveServices/accounts@2025-06-01"
  name                      = "acctest${random_string.unique.result}"
  parent_id                 = azapi_resource.resource_group.id
  location                  = var.location
  schema_validation_enabled = false

  body = {
    kind = "AIServices"
    sku = {
      name = "S0"
    }
    identity = {
      type = "SystemAssigned"
    }
    properties = {
      disableLocalAuth       = false
      allowProjectManagement = true
      customSubDomainName    = "acctest${random_string.unique.result}"
    }
  }
}

resource "azapi_resource" "foundry_deployment" {
  type      = "Microsoft.CognitiveServices/accounts/deployments@2023-05-01"
  name      = "gpt-5-mini"
  parent_id = azapi_resource.foundry.id

  body = {
    sku = {
      name     = "DataZoneStandard"
      capacity = 1
    }
    properties = {
      model = {
        format  = "OpenAI"
        name    = "gpt-5-mini"
        version = "2025-08-07"
      }
    }
  }
}

resource "azapi_resource" "foundry_project" {
  type                      = "Microsoft.CognitiveServices/accounts/projects@2025-06-01"
  name                      = "project${random_string.unique.result}"
  parent_id                 = azapi_resource.foundry.id
  location                  = var.location
  schema_validation_enabled = false

  body = {
    sku = {
      name = "S0"
    }
    identity = {
      type = "SystemAssigned"
    }
    properties = {
      displayName = "project"
      description = "Foundry project for agent example"
    }
  }
}

resource "azapi_resource" "foundry_project_user" {
  type      = "Microsoft.Authorization/roleAssignments@2022-04-01"
  name      = local.foundry_project_user_role_assignment_id
  parent_id = azapi_resource.foundry_project.id

  body = {
    properties = {
      principalId      = data.azapi_client_config.current.object_id
      principalType    = "User"
      roleDefinitionId = local.foundry_project_user_role_definition_id
    }
  }
}

resource "azapi_data_plane_resource" "agent" {
  type      = "Microsoft.Foundry/agents@v1"
  parent_id = "acctest${random_string.unique.result}.services.ai.azure.com/api/projects/${azapi_resource.foundry_project.name}"
  name      = "terraform-poc-agent"

  depends_on = [
    azapi_resource.foundry_deployment,
    azapi_resource.foundry_project_user,
  ]

  retry = {
    error_message_regex  = ["PermissionDenied", "Unauthorized", "authorization", "context deadline exceeded"]
    interval_seconds     = 30
    max_interval_seconds = 180
  }

  body = {
    name = "terraform-poc-agent"
    definition = {
      kind         = "prompt"
      model        = "gpt-5-mini"
      instructions = "You are a helpful agent created via Terraform"
    }
  }
}

output "agent_name" {
  value = azapi_data_plane_resource.agent.name
}

output "agent_resource_id" {
  value = azapi_data_plane_resource.agent.id
}

output "foundry_host" {
  value = "acctest${random_string.unique.result}.services.ai.azure.com"
}

output "project_name" {
  value = azapi_resource.foundry_project.name
}
