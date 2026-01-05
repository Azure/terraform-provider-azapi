terraform {
  required_providers {
    azapi = {
      source = "Azure/azapi"
    }
    azurerm = {
      source  = "hashicorp/azurerm"
      version = ">= 4.0"
    }
  }
}

provider "azapi" {
}

provider "azurerm" {
  features {}
  subscription_id = data.azapi_client_config.current.subscription_id
}

resource "azapi_resource" "resourceGroup" {
  type     = "Microsoft.Resources/resourceGroups@2021-04-01"
  name     = "example-contentunderstanding-rg"
  location = "WestUS3"
}

resource "azurerm_ai_services" "cognitiveAccount" {
  name                  = "examplefoundryproj"
  resource_group_name   = azapi_resource.resourceGroup.name
  location              = azapi_resource.resourceGroup.location
  sku_name              = "S0"
  custom_subdomain_name = "examplefoundryproj"

  network_acls {
    default_action = "Allow"
    ip_rules       = []
  }

  identity {
    type = "SystemAssigned"
  }
}

data "azapi_client_config" "current" {}

data "azapi_resource_list" "roleDefinitions" {
  type      = "Microsoft.Authorization/roleDefinitions@2022-04-01"
  parent_id = "/subscriptions/${data.azapi_client_config.current.subscription_id}"
  response_export_values = {
    cognitiveServicesUserRoleId = "value[?properties.roleName == 'Cognitive Services Contributor'].id | [0]"
  }
}

resource "azapi_resource" "roleAssignment" {
  type      = "Microsoft.Authorization/roleAssignments@2022-04-01"
  parent_id = azurerm_ai_services.cognitiveAccount.id
  name      = uuid()
  body = {
    properties = {
      principalId      = data.azapi_client_config.current.object_id
      roleDefinitionId = data.azapi_resource_list.roleDefinitions.output.cognitiveServicesUserRoleId
    }
  }
  lifecycle {
    ignore_changes = [name]
  }
}

resource "azurerm_cognitive_deployment" "gpt_41_mini" {
  name                   = "gpt-4.1-mini"
  cognitive_account_id   = azurerm_ai_services.cognitiveAccount.id
  version_upgrade_option = "OnceNewDefaultVersionAvailable"

  model {
    format  = "OpenAI"
    name    = "gpt-4.1-mini"
    version = "2025-04-14"
  }

  sku {
    name     = "GlobalStandard"
    capacity = 1
  }
}

resource "azurerm_cognitive_deployment" "gpt_41" {
  name                   = "gpt-4.1"
  cognitive_account_id   = azurerm_ai_services.cognitiveAccount.id
  version_upgrade_option = "OnceNewDefaultVersionAvailable"

  model {
    format  = "OpenAI"
    name    = "gpt-4.1"
    version = "2025-04-14"
  }

  sku {
    name     = "Standard"
    capacity = 1
  }
}

resource "azurerm_cognitive_deployment" "text_embedding_3_large" {
  name                   = "text-embedding-3-large"
  cognitive_account_id   = azurerm_ai_services.cognitiveAccount.id
  version_upgrade_option = "OnceNewDefaultVersionAvailable"

  model {
    format  = "OpenAI"
    name    = "text-embedding-3-large"
    version = "1"
  }

  sku {
    name     = "Standard"
    capacity = 1
  }
}

# Set Content Understanding defaults via data plane PATCH API
resource "terraform_data" "contentUnderstandingDefaults" {
  triggers_replace = [
    azurerm_ai_services.cognitiveAccount.id,
    azurerm_cognitive_deployment.gpt_41.id,
    azurerm_cognitive_deployment.gpt_41_mini.id,
    azurerm_cognitive_deployment.text_embedding_3_large.id
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
          "gpt-4.1"                = "gpt-4.1"
          "text-embedding-3-large" = "text-embedding-3-large"
        }
      } | ConvertTo-Json

      $endpoint = "${azurerm_ai_services.cognitiveAccount.endpoint}contentunderstanding/defaults?api-version=2025-11-01"
      Invoke-RestMethod -Uri $endpoint -Method Patch -Headers $headers -Body $body
    EOT
  }

  depends_on = [
    azapi_resource.roleAssignment,
    azurerm_cognitive_deployment.gpt_41,
    azurerm_cognitive_deployment.gpt_41_mini,
    azurerm_cognitive_deployment.text_embedding_3_large
  ]
}

resource "azapi_data_plane_resource" "example" {
  type      = "Microsoft.CognitiveServices/accounts/ContentUnderstanding/analyzers@2025-11-01"
  parent_id = trimprefix(azurerm_ai_services.cognitiveAccount.endpoint, "https://")
  name      = "exampleanalyzer"
  body = {
    description    = "My test analyzer"
    baseAnalyzerId = "prebuilt-document",
    models : {
      completion : "gpt-4.1",
      embedding : "text-embedding-3-large"
    },
  }
  depends_on = [
    terraform_data.contentUnderstandingDefaults,
  ]
}
