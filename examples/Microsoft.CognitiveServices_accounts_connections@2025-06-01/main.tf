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

variable "oauth2_client_id" {
  type        = string
  sensitive   = true
  description = "OAuth2 Client ID for the connection."
}

variable "oauth2_client_secret" {
  type        = string
  sensitive   = true
  description = "OAuth2 Client Secret for the connection."
}

variable "oauth2_tenant_id" {
  type        = string
  sensitive   = true
  description = "OAuth2 Tenant ID for the connection."
}

variable "oauth2_developer_token" {
  type        = string
  sensitive   = true
  description = "OAuth2 Developer Token for the connection."
}

variable "oauth2_refresh_token" {
  type        = string
  sensitive   = true
  description = "OAuth2 Refresh Token for the connection."
}

variable "oauth2_username" {
  type        = string
  sensitive   = true
  description = "OAuth2 Username for the connection."
}

variable "oauth2_password" {
  type        = string
  sensitive   = true
  description = "OAuth2 Password for the connection."
}

data "azapi_client_config" "current" {}

resource "azapi_resource" "resourceGroup" {
  type     = "Microsoft.Resources/resourceGroups@2020-06-01"
  name     = var.resource_name
  location = var.location
}

resource "azapi_resource" "userAssignedIdentity" {
  type                   = "Microsoft.ManagedIdentity/userAssignedIdentities@2024-11-30"
  name                   = var.resource_name
  location               = var.location
  parent_id              = azapi_resource.resourceGroup.id
  response_export_values = ["*"]
}

resource "azapi_resource" "account" {
  type      = "Microsoft.CognitiveServices/accounts@2025-06-01"
  parent_id = azapi_resource.resourceGroup.id
  name      = var.resource_name
  location  = var.location
  identity {
    type         = "SystemAssigned, UserAssigned"
    identity_ids = [azapi_resource.userAssignedIdentity.id]
  }
  body = {
    kind = "AIServices"
    properties = {
      allowProjectManagement = true
      allowedFqdnList = [
      ]
      apiProperties = {
      }
      disableLocalAuth              = false
      dynamicThrottlingEnabled      = false
      publicNetworkAccess           = "Enabled"
      restrictOutboundNetworkAccess = false
    }
    sku = {
      name = "S0"
      tier = "Standard"
    }
  }
  schema_validation_enabled = false
  response_export_values    = ["*"]
}

resource "azapi_resource" "account_openai" {
  type      = "Microsoft.CognitiveServices/accounts@2025-06-01"
  parent_id = azapi_resource.resourceGroup.id
  name      = "${var.resource_name}-openai"
  location  = var.location
  identity {
    type = "SystemAssigned"
  }

  body = {
    kind = "OpenAI"
    properties = {
      allowProjectManagement = true
      allowedFqdnList = [
      ]
      apiProperties = {
      }
      disableLocalAuth              = false
      dynamicThrottlingEnabled      = false
      publicNetworkAccess           = "Enabled"
      restrictOutboundNetworkAccess = false
    }
    sku = {
      name = "S0"
      tier = "Standard"
    }
  }
  schema_validation_enabled = false
  response_export_values    = ["*"]
}

resource "azapi_resource" "storageAccount" {
  type      = "Microsoft.Storage/storageAccounts@2021-09-01"
  parent_id = azapi_resource.resourceGroup.id
  name      = var.resource_name
  location  = var.location
  body = {
    kind = "StorageV2"
    properties = {
      accessTier                   = "Hot"
      allowBlobPublicAccess        = false
      allowCrossTenantReplication  = true
      allowSharedKeyAccess         = false
      defaultToOAuthAuthentication = false
      encryption = {
        keySource = "Microsoft.Storage"
        services = {
          queue = {
            keyType = "Service"
          }
          table = {
            keyType = "Service"
          }
        }
      }
      isHnsEnabled      = false
      isNfsV3Enabled    = false
      isSftpEnabled     = false
      minimumTlsVersion = "TLS1_2"
      networkAcls = {
        bypass        = "AzureServices"
        defaultAction = "Deny"
        resourceAccessRules = [
          {
            resourceId = azapi_resource.account.id
            tenantId   = data.azapi_client_config.current.tenant_id
          }
        ]
      }
      publicNetworkAccess      = "Enabled"
      supportsHttpsTrafficOnly = true
    }
    sku = {
      name = "Standard_LRS"
    }
  }
  schema_validation_enabled = false
  response_export_values    = ["*"]
}

resource "azapi_resource" "container" {
  type      = "Microsoft.Storage/storageAccounts/blobServices/containers@2024-01-01"
  parent_id = "${azapi_resource.storageAccount.id}/blobServices/default"
  name      = var.resource_name
  body = {
    properties = {
    }
  }

  schema_validation_enabled = false
  response_export_values    = ["*"]
}

# Retrieving keys
resource "azapi_resource_action" "account_keys" {
  type        = "Microsoft.CognitiveServices/accounts@2025-06-01"
  resource_id = azapi_resource.account.id
  action      = "listKeys"
  method      = "POST"

  sensitive_response_export_values = ["key1"]
}

resource "azapi_resource_action" "account_openai_keys" {
  type        = "Microsoft.CognitiveServices/accounts@2025-06-01"
  resource_id = azapi_resource.account_openai.id
  action      = "listKeys"
  method      = "POST"

  sensitive_response_export_values = ["key1", "key2"]
}

## Connections note:
# Credentials will not be returned since it's a sensitive data. if we want credentials, we can use .../{connectionName}/listsecrets
## Resources depend on each other so that they get deleted one after another instead of together.
# This helps escape a transient error that occurs when deleting all the connections together on cleanup.
resource "azapi_resource" "connection_aad" {
  type      = "Microsoft.CognitiveServices/accounts/connections@2025-06-01"
  parent_id = azapi_resource.account.id
  name      = "${var.resource_name}-aad"
  body = {
    properties = {
      authType = "AAD"
      category = "AzureBlob"
      target   = azapi_resource.storageAccount.output.properties.primaryEndpoints.blob
      metadata = {
        containerName = azapi_resource.container.name
        accountName   = azapi_resource.storageAccount.name
      }
    }
  }
  schema_validation_enabled = false
  ignore_casing             = false
  ignore_missing_property   = false
}

resource "azapi_resource" "connection_apikey" {
  type      = "Microsoft.CognitiveServices/accounts/connections@2025-06-01"
  parent_id = azapi_resource.account.id
  name      = "${var.resource_name}-apikey"
  body = {
    properties = {
      authType = "ApiKey"
      category = "AzureOpenAI"
      target   = azapi_resource.account_openai.output.properties.endpoint
      metadata = {
        ApiType    = "Azure"
        ResourceId = azapi_resource.account_openai.id
        location   = var.location
      }
    }
  }
  sensitive_body = {
    properties = {
      credentials = {
        key = azapi_resource_action.account_openai_keys.sensitive_output.key1
      }
    }
  }
  schema_validation_enabled = false
  ignore_casing             = false
  ignore_missing_property   = false
}

resource "azapi_resource" "connection_customkeys" {
  type      = "Microsoft.CognitiveServices/accounts/connections@2025-06-01"
  parent_id = azapi_resource.account.id
  name      = "${var.resource_name}-custom"
  body = {
    properties = {
      authType = "CustomKeys"
      category = "CustomKeys"
      target   = azapi_resource.account_openai.output.properties.endpoint
      metadata = {
        ApiType    = "Azure"
        ResourceId = azapi_resource.account_openai.id
        location   = var.location
      }
    }
  }
  sensitive_body = {
    properties = {
      credentials = {
        keys = {
          primaryKey   = azapi_resource_action.account_openai_keys.sensitive_output.key1
          secondaryKey = azapi_resource_action.account_openai_keys.sensitive_output.key2
        }
      }
    }
  }
  schema_validation_enabled = false
  ignore_casing             = false
  ignore_missing_property   = false
}

# This is example is based on having an external resource that uses OAuth2. 
resource "azapi_resource" "connection_oauth" {
  type      = "Microsoft.CognitiveServices/accounts/connections@2025-06-01"
  parent_id = azapi_resource.account.id
  name      = "${var.resource_name}-oauth"
  body = {
    properties = {
      authType = "OAuth2"
      category = "AzureBlob"
      target   = azapi_resource.storageAccount.output.properties.primaryEndpoints.blob
      metadata = {
        containerName = azapi_resource.container.name
        accountName   = azapi_resource.storageAccount.name
      }
    }
  }
  sensitive_body = {
    properties = {
      credentials = {
        # Not all fields are required.
        # Use the fields that are necessary in an actual use of the credentials, you don't need to use all of them, they are just placeholders for validation in this connection.
        authUrl        = "https://login.microsoftonline.com/${var.oauth2_tenant_id}/oauth2/v2.0/token"
        clientId       = var.oauth2_client_id
        clientSecret   = var.oauth2_client_secret
        tenantId       = var.oauth2_tenant_id
        developerToken = var.oauth2_developer_token
        refreshToken   = var.oauth2_refresh_token
        username       = var.oauth2_username
        password       = var.oauth2_password
      }
    }
  }
  schema_validation_enabled = false
  ignore_casing             = false
  ignore_missing_property   = false
}
