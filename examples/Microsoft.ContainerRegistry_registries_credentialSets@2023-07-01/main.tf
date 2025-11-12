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

data "azapi_client_config" "current" {}

variable "resource_name" {
  type    = string
  default = "acctest0001"
}

variable "location" {
  type    = string
  default = "eastus"
}

variable "credential_username" {
  type        = string
  sensitive   = true
  description = "The username for the container registry credential"
  default     = "testuser"
}

variable "credential_password" {
  type        = string
  sensitive   = true
  description = "The password for the container registry credential"
}

resource "azapi_resource" "resourceGroup" {
  type     = "Microsoft.Resources/resourceGroups@2020-06-01"
  name     = var.resource_name
  location = var.location
}

resource "azapi_resource" "registry" {
  type      = "Microsoft.ContainerRegistry/registries@2023-11-01-preview"
  parent_id = azapi_resource.resourceGroup.id
  name      = var.resource_name
  location  = var.location
  body = {
    properties = {
      adminUserEnabled         = false
      anonymousPullEnabled     = false
      dataEndpointEnabled      = false
      networkRuleBypassOptions = "AzureServices"
      policies = {
        exportPolicy = {
          status = "enabled"
        }
        quarantinePolicy = {
          status = "disabled"
        }
        retentionPolicy = {}
        trustPolicy     = {}
      }
      publicNetworkAccess = "Enabled"
      zoneRedundancy      = "Disabled"
    }
    sku = {
      name = "Basic"
    }
  }
}

resource "azapi_resource" "vault" {
  type      = "Microsoft.KeyVault/vaults@2023-02-01"
  parent_id = azapi_resource.resourceGroup.id
  name      = "${var.resource_name}vault"
  location  = var.location
  body = {
    properties = {
      accessPolicies = [{
        objectId = data.azapi_client_config.current.object_id
        permissions = {
          certificates = []
          keys         = []
          secrets      = ["Get", "Set", "Delete", "Purge"]
          storage      = []
        }
        tenantId = data.azapi_client_config.current.tenant_id
      }]
      createMode                   = "default"
      enableRbacAuthorization      = false
      enableSoftDelete             = true
      enabledForDeployment         = false
      enabledForDiskEncryption     = false
      enabledForTemplateDeployment = false
      publicNetworkAccess          = "Enabled"
      sku = {
        family = "A"
        name   = "standard"
      }
      softDeleteRetentionInDays = 7
      tenantId                  = data.azapi_client_config.current.tenant_id
    }
  }
}

resource "azapi_resource" "usernameSecret" {
  type      = "Microsoft.KeyVault/vaults/secrets@2023-02-01"
  parent_id = azapi_resource.vault.id
  name      = "username"
  body = {
    properties = {
      value = var.credential_username
    }
  }
}

resource "azapi_resource" "passwordSecret" {
  type      = "Microsoft.KeyVault/vaults/secrets@2023-02-01"
  parent_id = azapi_resource.vault.id
  name      = "password"
  body = {
    properties = {
      value = var.credential_password
    }
  }
}

resource "azapi_resource" "credentialSet" {
  type      = "Microsoft.ContainerRegistry/registries/credentialSets@2023-07-01"
  parent_id = azapi_resource.registry.id
  name      = "${var.resource_name}-acr-credential-set"
  identity {
    type         = "SystemAssigned"
    identity_ids = []
  }
  body = {
    properties = {
      authCredentials = [{
        name                     = "Credential1"
        passwordSecretIdentifier = "https://${var.resource_name}vault.vault.azure.net/secrets/password"
        usernameSecretIdentifier = "https://${var.resource_name}vault.vault.azure.net/secrets/username"
      }]
      loginServer = "docker.io"
    }
  }
}

