---
subcategory: "Microsoft.Sql - Azure SQL Database, Azure SQL Managed Instance, Azure Synapse Analytics"
page_title: "servers/jobAgents/jobs/steps"
description: |-
  Manages a Elastic Job Step.
---

# Microsoft.Sql/servers/jobAgents/jobs/steps - Elastic Job Step

This article demonstrates how to use `azapi` provider to manage the Elastic Job Step resource in Azure.



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


```



## Arguments Reference

The following arguments are supported:

* `type` - (Required) The type of the resource. This should be set to `Microsoft.Sql/servers/jobAgents/jobs/steps@api-version`. The available api-versions for this resource are: [`2017-03-01-preview`, `2020-02-02-preview`, `2020-08-01-preview`, `2020-11-01-preview`, `2021-02-01-preview`, `2021-05-01-preview`, `2021-08-01-preview`, `2021-11-01`, `2021-11-01-preview`, `2022-02-01-preview`, `2022-05-01-preview`, `2022-08-01-preview`, `2022-11-01-preview`, `2023-02-01-preview`, `2023-05-01-preview`, `2023-08-01`, `2023-08-01-preview`, `2024-05-01-preview`, `2024-11-01-preview`].

* `parent_id` - (Required) The ID of the azure resource in which this resource is created. The allowed values are:  
  `/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Sql/servers/{resourceName}/jobAgents/{resourceName}/jobs/{resourceName}`

* `name` - (Required) Specifies the name of the azure resource. Changing this forces a new resource to be created.

* `body` - (Required) Specifies the configuration of the resource. More information about the arguments in `body` can be found in the [Microsoft documentation](https://learn.microsoft.com/en-us/azure/templates/Microsoft.Sql/servers/jobAgents/jobs/steps?pivots=deployment-language-terraform).

For other arguments, please refer to the [azapi_resource](https://registry.terraform.io/providers/Azure/azapi/latest/docs/resources/resource) documentation.

## Import

 ```shell
 # Azure resource can be imported using the resource id, e.g.
 terraform import azapi_resource.example /subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Sql/servers/{resourceName}/jobAgents/{resourceName}/jobs/{resourceName}/steps/{resourceName}
 
 # It also supports specifying API version by using the resource id with api-version as a query parameter, e.g.
 terraform import azapi_resource.example /subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Sql/servers/{resourceName}/jobAgents/{resourceName}/jobs/{resourceName}/steps/{resourceName}?api-version=2024-11-01-preview
 ```
