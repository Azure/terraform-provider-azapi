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

locals {
  content_understanding_owner_role_definition_id = "/subscriptions/${data.azapi_client_config.current.subscription_id}/providers/Microsoft.Authorization/roleDefinitions/4b42bd01-da42-4c92-9b07-15ea5bd6a602"
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
}

variable "resource_group_tags" {
  type        = map(string)
  description = "Optional tags to apply to the resource group (useful for policy-constrained subscriptions)."
  default     = {}
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
  response_export_values = ["properties.endpoint"]
}

resource "azapi_resource" "contentUnderstandingOwnerRoleAssignment" {
  type      = "Microsoft.Authorization/roleAssignments@2022-04-01"
  parent_id = azapi_resource.foundry.id
  name      = uuid()
  body = {
    properties = {
      principalId      = data.azapi_client_config.current.object_id
      roleDefinitionId = local.content_understanding_owner_role_definition_id
    }
  }
}

resource "azapi_resource" "foundry_gpt5_deployment" {
  type      = "Microsoft.CognitiveServices/accounts/deployments@2023-05-01"
  name      = "gpt-5"
  parent_id = azapi_resource.foundry.id

  body = {
    sku = {
      name     = "DataZoneStandard"
      capacity = 1
    }
    properties = {
      model = {
        format  = "OpenAI"
        name    = "gpt-5"
        version = "2025-08-07"
      }
    }
  }
}

resource "azapi_resource" "foundry_embedding_deployment" {
  type      = "Microsoft.CognitiveServices/accounts/deployments@2023-05-01"
  name      = "text-embedding-3-large"
  parent_id = azapi_resource.foundry.id

  body = {
    sku = {
      name     = "Standard"
      capacity = 1
    }
    properties = {
      model = {
        format  = "OpenAI"
        name    = "text-embedding-3-large"
        version = "1"
      }
    }
  }
}

resource "terraform_data" "contentUnderstandingDefaults" {
  triggers_replace = [
    azapi_resource.foundry.id,
    azapi_resource.foundry_gpt5_deployment.id,
    azapi_resource.foundry_embedding_deployment.id,
  ]

  provisioner "local-exec" {
    interpreter = ["pwsh", "-Command"]
    command     = <<-EOT
      $token = (az account get-access-token --resource https://cognitiveservices.azure.com --query accessToken -o tsv)
      $headers = @{
        "Authorization" = "Bearer $token"
        "Content-Type"  = "application/merge-patch+json"
      }
      $body = @{
        modelDeployments = @{
          "gpt-5"                  = "gpt-5"
          "text-embedding-3-large" = "text-embedding-3-large"
        }
      } | ConvertTo-Json

      $endpoint = "${azapi_resource.foundry.output.properties.endpoint}contentunderstanding/defaults?api-version=2025-11-01"
      Invoke-RestMethod -Uri $endpoint -Method Patch -Headers $headers -Body $body
    EOT
  }

  depends_on = [
    azapi_resource.contentUnderstandingOwnerRoleAssignment,
    azapi_resource.foundry_gpt5_deployment,
    azapi_resource.foundry_embedding_deployment,
  ]
}

resource "azapi_data_plane_resource" "example" {
  type      = "Microsoft.CognitiveServices/accounts/ContentUnderstanding/analyzers@2025-11-01"
  parent_id = trimprefix(azapi_resource.foundry.output.properties.endpoint, "https://")
  name      = "exampleanalyzer"

  body = {
    description    = "My test analyzer"
    baseAnalyzerId = "prebuilt-document"
    models = {
      completion = "gpt-5"
      embedding  = "text-embedding-3-large"
    }
  }

  depends_on = [terraform_data.contentUnderstandingDefaults]
}