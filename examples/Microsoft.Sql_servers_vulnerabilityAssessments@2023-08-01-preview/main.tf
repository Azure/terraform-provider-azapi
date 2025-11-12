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
  default = "centralus"
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
  name      = var.resource_name
  location  = var.location
  body = {
    properties = {
      administratorLogin            = "mradministrator"
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
  name      = "${var.resource_name}sa"
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
      name = "Standard_GRS"
    }
  }
}

resource "azapi_resource_action" "storageAccountKeys" {
  type                   = "Microsoft.Storage/storageAccounts@2023-05-01"
  resource_id            = azapi_resource.storageAccount.id
  action                 = "listKeys"
  response_export_values = ["keys"]
}

resource "azapi_resource" "storageContainer" {
  type      = "Microsoft.Storage/storageAccounts/blobServices/containers@2023-05-01"
  parent_id = "${azapi_resource.storageAccount.id}/blobServices/default"
  name      = "vulnerability-assessment"
  body = {
    properties = {}
  }
}

resource "azapi_resource" "securityAlertPolicy" {
  type      = "Microsoft.Sql/servers/securityAlertPolicies@2023-08-01-preview"
  parent_id = azapi_resource.server.id
  name      = "default"
  body = {
    properties = {
      state = "Enabled"
    }
  }
}

resource "azapi_resource" "vulnerabilityAssessment" {
  type      = "Microsoft.Sql/servers/vulnerabilityAssessments@2023-08-01-preview"
  parent_id = azapi_resource.server.id
  name      = "default"
  body = {
    properties = {
      recurringScans = {
        emailSubscriptionAdmins = false
        emails                  = []
        isEnabled               = false
      }
      storageAccountAccessKey = azapi_resource_action.storageAccountKeys.output.keys[0].value
      storageContainerPath    = "${azapi_resource.storageAccount.output.properties.primaryEndpoints.blob}${azapi_resource.storageContainer.name}/"
    }
  }
  depends_on = [azapi_resource.securityAlertPolicy]
}

