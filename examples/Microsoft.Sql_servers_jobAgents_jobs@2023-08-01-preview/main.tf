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

resource "azapi_resource" "resourceGroup" {
  type     = "Microsoft.Resources/resourceGroups@2020-06-01"
  name     = var.resource_name
  location = var.location
}

resource "azapi_resource" "server" {
  type      = "Microsoft.Sql/servers@2023-08-01-preview"
  parent_id = azapi_resource.resourceGroup.id
  name      = "${var.resource_name}-server"
  location  = var.location
  body = {
    properties = {
      administratorLogin            = "4dm1n157r470r"
      administratorLoginPassword    = "4-v3ry-53cr37-p455w0rd"
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
      maintenanceConfigurationId       = "/subscriptions/${data.azapi_client_config.current.subscription_id}/providers/Microsoft.Maintenance/publicMaintenanceConfigurations/SQL_Default"
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

resource "azapi_resource" "jobAgent" {
  type      = "Microsoft.Sql/servers/jobAgents@2023-08-01-preview"
  parent_id = azapi_resource.server.id
  name      = "${var.resource_name}-job-agent"
  location  = var.location
  body = {
    properties = {
      databaseId = azapi_resource.database.id
    }
    sku = {
      name = "JA100"
    }
  }
}

resource "azapi_resource" "job" {
  type      = "Microsoft.Sql/servers/jobAgents/jobs@2023-08-01-preview"
  parent_id = azapi_resource.jobAgent.id
  name      = "${var.resource_name}-job"
  body = {
    properties = {
      description = ""
    }
  }
}

