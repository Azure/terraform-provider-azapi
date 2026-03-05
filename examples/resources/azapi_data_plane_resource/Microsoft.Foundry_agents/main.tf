terraform {
  required_providers {
    azapi = {
      source = "Azure/azapi"
    }
  }
}

provider "azapi" {
}

variable "foundry_host" {
  type        = string
  description = "Foundry host name, e.g. <aiservices-id>.services.ai.azure.com"
}

variable "project_name" {
  type        = string
  description = "Foundry project name"
}

resource "azapi_data_plane_resource" "agent" {
  type      = "Microsoft.Foundry/agents@v1"
  parent_id = "${var.foundry_host}/api/projects/${var.project_name}"
  name      = "terraform-poc-agent"

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
