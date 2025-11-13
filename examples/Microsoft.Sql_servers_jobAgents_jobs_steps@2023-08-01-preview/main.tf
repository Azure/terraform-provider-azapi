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

variable "administrator_login_password" {
  type        = string
  sensitive   = true
  description = "The administrator login password for the SQL server"
}

variable "job_credential_password" {
  type        = string
  sensitive   = true
  description = "The password for the SQL job credential"
}

data "azapi_client_config" "current" {}

locals {
  maintenance_config_id = "/subscriptions/${data.azapi_client_config.current.subscription_id}/providers/Microsoft.Maintenance/publicMaintenanceConfigurations/SQL_Default"
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
      collation                  = "SQL_Latin1_General_CP1_CI_AS"
      createMode                 = "Default"
      maintenanceConfigurationId = local.maintenance_config_id
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

resource "azapi_resource" "credential" {
  type      = "Microsoft.Sql/servers/jobAgents/credentials@2023-08-01-preview"
  parent_id = azapi_resource.jobAgent.id
  name      = "${var.resource_name}-job-credential"
  body = {
    properties = {
      password = var.job_credential_password
      username = "testusername"
    }
  }
}

resource "azapi_resource" "targetGroup" {
  type      = "Microsoft.Sql/servers/jobAgents/targetGroups@2023-08-01-preview"
  parent_id = azapi_resource.jobAgent.id
  name      = "${var.resource_name}-target-group"
  body = {
    properties = {
      members = []
    }
  }
}

resource "azapi_resource" "step" {
  type      = "Microsoft.Sql/servers/jobAgents/jobs/steps@2023-08-01-preview"
  parent_id = azapi_resource.job.id
  name      = "${var.resource_name}-job-step"
  body = {
    properties = {
      action = {
        value = "IF NOT EXISTS (SELECT * FROM sys.objects WHERE [name] = N'Person')\n  CREATE TABLE Person (\n    FirstName NVARCHAR(50),\n    LastName NVARCHAR(50),\n  );\n"
      }
      credential = azapi_resource.credential.id
      executionOptions = {
        initialRetryIntervalSeconds    = 1
        maximumRetryIntervalSeconds    = 120
        retryAttempts                  = 10
        retryIntervalBackoffMultiplier = 2
        timeoutSeconds                 = 43200
      }
      stepId      = 1
      targetGroup = azapi_resource.targetGroup.id
    }
  }
}

