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

variable "secondary_location" {
  type    = string
  default = "eastus"
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
  name      = "${var.resource_name}-primary"
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

resource "azapi_resource" "server_1" {
  type      = "Microsoft.Sql/servers@2023-08-01-preview"
  parent_id = azapi_resource.resourceGroup.id
  name      = "${var.resource_name}-secondary"
  location  = var.secondary_location
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

resource "azapi_resource" "database" {
  type      = "Microsoft.Sql/servers/databases@2023-08-01-preview"
  parent_id = azapi_resource.server.id
  name      = "${var.resource_name}-db"
  location  = var.location
  body = {
    properties = {
      autoPauseDelay                   = 0
      collation                        = "SQL_Latin1_General_CP1_CI_AS"
      createMode                       = "Default"
      elasticPoolId                    = ""
      encryptionProtectorAutoRotation  = false
      highAvailabilityReplicaCount     = 0
      isLedgerOn                       = false
      licenseType                      = ""
      maxSizeBytes                     = 214748364800
      minCapacity                      = 0
      readScale                        = "Disabled"
      requestedBackupStorageRedundancy = "Geo"
      sampleName                       = ""
      secondaryType                    = ""
      zoneRedundant                    = false
    }
    sku = {
      name = "S1"
    }
  }
}

resource "azapi_resource" "failoverGroup" {
  type      = "Microsoft.Sql/servers/failoverGroups@2023-08-01-preview"
  parent_id = azapi_resource.server.id
  name      = "${var.resource_name}-fg"
  body = {
    properties = {
      databases = [azapi_resource.database.id]
      partnerServers = [{
        id = azapi_resource.server_1.id
      }]
      readOnlyEndpoint = {
        failoverPolicy = "Disabled"
      }
      readWriteEndpoint = {
        failoverPolicy                         = "Automatic"
        failoverWithDataLossGracePeriodMinutes = 60
      }
    }
  }
}

