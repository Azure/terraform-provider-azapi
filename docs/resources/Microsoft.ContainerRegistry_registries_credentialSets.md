---
subcategory: "Microsoft.ContainerRegistry - Container Registry"
page_title: "registries/credentialSets"
description: |-
  Manages a Container Registry Credential Set.
---

# Microsoft.ContainerRegistry/registries/credentialSets - Container Registry Credential Set

This article demonstrates how to use `azapi` provider to manage the Container Registry Credential Set resource in Azure.



## Example Usage

### default

```hcl
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


```



## Arguments Reference

The following arguments are supported:

* `type` - (Required) The type of the resource. This should be set to `Microsoft.ContainerRegistry/registries/credentialSets@api-version`. The available api-versions for this resource are: [`2023-01-01-preview`, `2023-06-01-preview`, `2023-07-01`, `2023-08-01-preview`, `2023-11-01-preview`, `2024-11-01-preview`, `2025-03-01-preview`, `2025-04-01`, `2025-05-01-preview`].

* `parent_id` - (Required) The ID of the azure resource in which this resource is created. The allowed values are:  
  `/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.ContainerRegistry/registries/{resourceName}`

* `name` - (Required) Specifies the name of the azure resource. Changing this forces a new resource to be created.

* `body` - (Required) Specifies the configuration of the resource. More information about the arguments in `body` can be found in the [Microsoft documentation](https://learn.microsoft.com/en-us/azure/templates/Microsoft.ContainerRegistry/registries/credentialSets?pivots=deployment-language-terraform).

For other arguments, please refer to the [azapi_resource](https://registry.terraform.io/providers/Azure/azapi/latest/docs/resources/resource) documentation.

## Import

 ```shell
 # Azure resource can be imported using the resource id, e.g.
 terraform import azapi_resource.example /subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.ContainerRegistry/registries/{resourceName}/credentialSets/{resourceName}
 
 # It also supports specifying API version by using the resource id with api-version as a query parameter, e.g.
 terraform import azapi_resource.example /subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.ContainerRegistry/registries/{resourceName}/credentialSets/{resourceName}?api-version=2025-05-01-preview
 ```
