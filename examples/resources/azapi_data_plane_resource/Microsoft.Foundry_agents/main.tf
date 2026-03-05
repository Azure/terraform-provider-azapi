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
  default     = "westus3"
}

variable "resource_group_tags" {
  type        = map(string)
  description = "Optional tags to apply to the resource group (useful for policy-constrained subscriptions)."
  default     = {}
}

resource "azapi_resource" "resource_group" {
  type     = "Microsoft.Resources/resourceGroups@2021-04-01"
  name     = "rg-foundry-${random_string.unique.result}"
  location = var.location
  tags     = var.resource_group_tags
}

resource "azapi_resource" "foundry" {
  type                      = "Microsoft.CognitiveServices/accounts@2025-06-01"
  name                      = "aifoundry${random_string.unique.result}"
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
      customSubDomainName    = "aifoundry${random_string.unique.result}"
    }
  }
}

resource "azapi_resource" "foundry_deployment_gpt_4o" {
  type      = "Microsoft.CognitiveServices/accounts/deployments@2023-05-01"
  name      = "gpt-4o"
  parent_id = azapi_resource.foundry.id

  body = {
    sku = {
      name     = "GlobalStandard"
      capacity = 1
    }
    properties = {
      model = {
        format  = "OpenAI"
        name    = "gpt-4o"
        version = "2024-11-20"
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

resource "azapi_data_plane_resource" "agent" {
  type      = "Microsoft.Foundry/agents@v1"
  parent_id = "aifoundry${random_string.unique.result}.services.ai.azure.com/api/projects/${azapi_resource.foundry_project.name}"
  name      = "terraform-poc-agent"

  depends_on = [
    azapi_resource.foundry_deployment_gpt_4o,
    azapi_resource.foundry_project,
  ]

  body = {
    name = "terraform-poc-agent"
    definition = {
      kind         = "prompt"
      model        = "gpt-4o"
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
  value = "aifoundry${random_string.unique.result}.services.ai.azure.com"
}

output "project_name" {
  value = azapi_resource.foundry_project.name
}
