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
  default = "westus"
}

resource "azapi_resource" "resourceGroup" {
  type     = "Microsoft.Resources/resourceGroups@2020-06-01"
  name     = var.resource_name
  location = var.location
}

resource "azapi_resource" "storageAccount" {
  type      = "Microsoft.Storage/storageAccounts@2023-05-01"
  parent_id = azapi_resource.resourceGroup.id
  name      = var.resource_name
  location  = var.location
  body = {
    kind = "StorageV2"
    properties = {
      accessTier                   = "Hot"
      allowBlobPublicAccess        = false
      allowCrossTenantReplication  = false
      allowSharedKeyAccess         = true
      defaultToOAuthAuthentication = false
      dnsEndpointType              = "Standard"
      encryption = {
        keySource = "Microsoft.Storage"
        services = {
          blob = {
            enabled = true
            keyType = "Account"
          }
          file = {
            enabled = true
            keyType = "Account"
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

resource "azapi_resource" "export" {
  type      = "Microsoft.CostManagement/exports@2023-08-01"
  parent_id = azapi_resource.resourceGroup.id
  name      = var.resource_name
  body = {
    properties = {
      definition = {
        timeframe = "TheLastMonth"
        type      = "Usage"
        dataSet = {
          granularity = "Daily"
        }
      }
      deliveryInfo = {
        destination = {
          container      = "exports"
          resourceId     = azapi_resource.storageAccount.id
          rootFolderPath = "ad-hoc"
        }
      }
      format = "Csv"
      schedule = {
        recurrence = "Monthly"
        recurrencePeriod = {
          from = timeadd(timestamp(), "24h")
          to   = timeadd(timestamp(), "744h")
        }
        status = "Active"
      }
    }
  }
}
