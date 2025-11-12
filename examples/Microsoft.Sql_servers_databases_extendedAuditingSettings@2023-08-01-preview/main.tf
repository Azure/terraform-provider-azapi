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
  default = "westus"
}

variable "administrator_login_password" {
  type        = string
  sensitive   = true
  description = "The administrator login password for the SQL server"
}

resource "azapi_resource" "resourceGroup" {
  type     = "Microsoft.Resources/resourceGroups@2020-06-01"
  name     = var.resource_name
  location = var.location
}

resource "azapi_resource" "server" {
  type      = "Microsoft.Sql/servers@2023-08-01-preview"
  parent_id = azapi_resource.resourceGroup.id
  name      = "${var.resource_name}-sqlserver"
  location  = var.location
  body = {
    properties = {
      administratorLogin            = "missadministrator"
      administratorLoginPassword    = var.administrator_login_password
      minimalTlsVersion             = "1.2"
      publicNetworkAccess           = "Enabled"
      restrictOutboundNetworkAccess = "Disabled"
      version                       = "12.0"
    }
  }
}

resource "azapi_resource" "storageAccount" {
  type      = "Microsoft.Storage/storageAccounts@2023-05-01"
  parent_id = azapi_resource.resourceGroup.id
  name      = "${replace(var.resource_name, "-", "")}sa"
  location  = var.location
  body = {
    kind = "StorageV2"
    properties = {
      accessTier                   = "Hot"
      allowBlobPublicAccess        = true
      allowCrossTenantReplication  = false
      allowSharedKeyAccess         = true
      defaultToOAuthAuthentication = false
      dnsEndpointType              = "Standard"
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
      isHnsEnabled       = false
      isLocalUserEnabled = true
      isNfsV3Enabled     = false
      isSftpEnabled      = false
      minimumTlsVersion  = "TLS1_2"
      networkAcls = {
        bypass              = "AzureServices"
        defaultAction       = "Allow"
        ipRules             = []
        resourceAccessRules = []
        virtualNetworkRules = []
      }
      publicNetworkAccess      = "Enabled"
      supportsHttpsTrafficOnly = true
    }
    sku = {
      name = "Standard_LRS"
    }
  }
}

resource "azapi_resource" "database" {
  type      = "Microsoft.Sql/servers/databases@2023-08-01-preview"
  parent_id = azapi_resource.server.id
  name      = "${var.resource_name}-db"
  location  = var.location
  body = {
    properties = {
      createMode                       = "Default"
      requestedBackupStorageRedundancy = "Geo"
    }
  }
}

resource "azapi_resource_action" "storageAccountKeys" {
  type                   = "Microsoft.Storage/storageAccounts@2023-05-01"
  resource_id            = azapi_resource.storageAccount.id
  action                 = "listKeys"
  method                 = "POST"
  response_export_values = ["*"]
}

resource "azapi_resource" "extendedAuditingSetting" {
  type      = "Microsoft.Sql/servers/databases/extendedAuditingSettings@2023-08-01-preview"
  parent_id = azapi_resource.database.id
  name      = "default"
  body = {
    properties = {
      isAzureMonitorTargetEnabled = true
      isStorageSecondaryKeyInUse  = false
      retentionDays               = 0
      state                       = "Enabled"
      storageAccountAccessKey     = azapi_resource_action.storageAccountKeys.output.keys[0].value
      storageEndpoint             = azapi_resource.storageAccount.output.properties.primaryEndpoints.blob
    }
  }
}

