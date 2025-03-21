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

variable "pat" {
  type        = string
  sensitive   = true
  description = "GitHub Personal Access Token"
}

resource "azapi_resource" "resourceGroup" {
  type     = "Microsoft.Resources/resourceGroups@2020-06-01"
  name     = var.resource_name
  location = var.location
}

resource "azapi_resource" "automationAccount" {
  type      = "Microsoft.Automation/automationAccounts@2023-11-01"
  parent_id = azapi_resource.resourceGroup.id
  name      = var.resource_name
  location  = var.location
  identity {
    type         = "SystemAssigned"
    identity_ids = []
  }

  body = {
    properties = {
      encryption = {
        keySource = "Microsoft.Automation"
      }
      publicNetworkAccess = true
      sku = {
        name = "Basic"
      }
    }
  }
}

resource "azapi_resource" "sourceControl" {
  type      = "Microsoft.Automation/automationAccounts/sourceControls@2023-11-01"
  name      = var.resource_name
  parent_id = azapi_resource.automationAccount.id

  body = {
    properties = {
      repoUrl        = "https://github.com/Azure-Samples/acr-build-helloworld-node.git"
      branch         = "master"
      sourceType     = "GitHub"
      folderPath     = "/"
      autoSync       = false
      publishRunbook = false

      securityToken = {
        tokenType   = "PersonalAccessToken"
        accessToken = var.pat
      }
    }
  }
}

data "azapi_resource_list" "roleDefinitions" {
  type      = "Microsoft.Authorization/roleDefinitions@2022-04-01"
  parent_id = azapi_resource.automationAccount.id
  response_export_values = {
    contributorRoleDefinitionId = "value[?properties.roleName == 'Contributor'].id | [0]"
  }
}

resource "azapi_resource" "roleAssignment" {
  type      = "Microsoft.Authorization/roleAssignments@2022-04-01"
  parent_id = azapi_resource.automationAccount.id
  name      = uuid()
  body = {
    properties = {
      principalId      = azapi_resource.automationAccount.identity[0].principal_id
      roleDefinitionId = data.azapi_resource_list.roleDefinitions.output.contributorRoleDefinitionId
    }
  }
  lifecycle {
    ignore_changes = [name]
  }
}


# use azapi_resource_action to create a sourceControlSyncJob because the API doesn't support DELETE
resource "azapi_resource_action" "sourceControlSyncJob" {
  type        = "Microsoft.Automation/automationAccounts/sourceControls/sourceControlSyncJobs@2023-11-01"
  resource_id = provider::azapi::build_resource_id(azapi_resource.sourceControl.id, "Microsoft.Automation/automationAccounts/sourceControls/sourceControlSyncJobs", uuid())
  method      = "PUT"

  body = {
    properties = {
      commitId = ""
    }
  }
  depends_on = [azapi_resource.roleAssignment]
  lifecycle {
    ignore_changes = [resource_id]
  }
}
