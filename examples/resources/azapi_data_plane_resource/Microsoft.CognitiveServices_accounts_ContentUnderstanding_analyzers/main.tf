terraform {
  required_providers {
    azapi = {
      source = "Azure/azapi"
    }
  }
}

provider "azapi" {
}

resource "azapi_resource" "resourceGroup" {
  type     = "Microsoft.Resources/resourceGroups@2021-04-01"
  name     = "example-resources"
  location = "eastus2"
}

resource "azapi_resource" "cognitiveAccount" {
  type      = "Microsoft.CognitiveServices/accounts@2024-10-01"
  name      = "exampleai"
  parent_id = azapi_resource.resourceGroup.id
  location  = azapi_resource.resourceGroup.location
  identity {
    type         = "SystemAssigned"
    identity_ids = []
  }
  body = {
    kind = "AIServices"
    properties = {
      publicNetworkAccess = "Enabled"
      customSubDomainName = "exampleai"
    }
    sku = {
      name = "S0"
    }
  }
  response_export_values = ["properties.endpoint"]
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
  parent_id = azapi_resource.cognitiveAccount.id
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

# Deploy GPT-4.1 model (required for Content Understanding)
resource "azapi_resource" "gptDeployment" {
  type      = "Microsoft.CognitiveServices/accounts/deployments@2024-10-01"
  name      = "gpt-4.1"
  parent_id = azapi_resource.cognitiveAccount.id
  body = {
    properties = {
      model = {
        format  = "OpenAI"
        name    = "gpt-4.1"
        version = "2025-04-14"
      }
    }
    sku = {
      name     = "GlobalStandard"
      capacity = 25
    }
  }
}

# Deploy text-embedding-3-large model (required for Content Understanding)
resource "azapi_resource" "embeddingDeployment" {
  type      = "Microsoft.CognitiveServices/accounts/deployments@2024-10-01"
  name      = "text-embedding-3-large"
  parent_id = azapi_resource.cognitiveAccount.id
  body = {
    properties = {
      model = {
        format  = "OpenAI"
        name    = "text-embedding-3-large"
        version = "1"
      }
    }
    sku = {
      name     = "Standard"
      capacity = 120
    }
  }
}

# Set Content Understanding defaults via data plane PATCH API
resource "terraform_data" "contentUnderstandingDefaults" {
  triggers_replace = [
    azapi_resource.cognitiveAccount.id,
    azapi_resource.gptDeployment.id,
    azapi_resource.embeddingDeployment.id
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
      
      $endpoint = "${azapi_resource.cognitiveAccount.output.properties.endpoint}contentunderstanding/defaults?api-version=2025-05-01-preview"
      Invoke-RestMethod -Uri $endpoint -Method Patch -Headers $headers -Body $body
    EOT
  }

  depends_on = [
    azapi_resource.roleAssignment,
    azapi_resource.gptDeployment,
    azapi_resource.embeddingDeployment
  ]
}

resource "azapi_data_plane_resource" "example" {
  type      = "Microsoft.CognitiveServices/accounts/ContentUnderstanding/analyzers@2025-05-01-preview"
  parent_id = trimprefix(azapi_resource.cognitiveAccount.output.properties.endpoint, "https://")
  name      = "exampleanalyzer"
  body = {
    description    = "My test analyzer"
    baseAnalyzerId = "prebuilt-document"
  }
  depends_on = [
    terraform_data.contentUnderstandingDefaults,
  ]
}
