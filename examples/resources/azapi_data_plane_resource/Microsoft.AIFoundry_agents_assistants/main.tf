terraform {
  required_providers {
    azapi = {
      source = "Azure/azapi"
    }
  }
}

provider "azapi" {
  endpoint = [{
    # Project endpoint base (no protocol) for AI Foundry Agents data plane
    resource_manager_endpoint = var.ai_foundry_host
    resource_manager_audience = "https://ai.azure.com"
  }]
}

variable "ai_foundry_host" {
  type        = string
  description = "AI Foundry host name, e.g. <aiservices-id>.services.ai.azure.com"
}

variable "project_name" {
  type        = string
  description = "AI Foundry project name"
}

resource "azapi_data_plane_resource" "assistant" {
  type      = "Microsoft.AIFoundry/agents/assistants@v1"
  parent_id = "${var.ai_foundry_host}/api/projects/${var.project_name}"

  # name is server-generated on create (e.g. asst_...), so it is omitted.
  body = {
    model        = "gpt-4o"
    instructions = "You are a helpful assistant created via Terraform"
    name         = "terraform-poc-agent"
  }
}

output "assistant_id" {
  value = azapi_data_plane_resource.assistant.name
}

output "assistant_resource_id" {
  value = azapi_data_plane_resource.assistant.id
}
