---
subcategory: "Microsoft.Sql - Azure SQL Database, Azure SQL Managed Instance, Azure Synapse Analytics"
page_title: "servers/jobAgents"
description: |-
  Manages a Elastic Job Agent.
---

# Microsoft.Sql/servers/jobAgents - Elastic Job Agent

This article demonstrates how to use `azapi` provider to manage the Elastic Job Agent resource in Azure.

## Example Usage

### default

```hcl
terraform {
  required_providers {
    azapi = {
      source = "Azure/azapi"
    }
    azurerm = {
      source = "hashicorp/azurerm"
    }
  }
}

provider "azurerm" {
  features {
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

variable "sql_administrator_password" {
  type        = string
  description = "The administrator password for the SQL server"
  sensitive   = true
}

data "azurerm_client_config" "current" {
}

resource "azapi_resource" "resourceGroup" {
  type     = "Microsoft.Resources/resourceGroups@2020-06-01"
  name     = var.resource_name
  location = var.location
}

resource "azapi_resource" "server" {
  type      = "Microsoft.Sql/servers@2021-02-01-preview"
  parent_id = azapi_resource.resourceGroup.id
  name      = var.resource_name
  location  = var.location
  body = {
    properties = {
      administratorLogin            = "4dministr4t0r"
      administratorLoginPassword    = var.sql_administrator_password
      minimalTlsVersion             = "1.2"
      publicNetworkAccess           = "Enabled"
      restrictOutboundNetworkAccess = "Disabled"
      version                       = "12.0"
    }
  }
  schema_validation_enabled = false
  response_export_values    = ["*"]
}

data "azapi_resource_id" "publicMaintenanceConfiguration" {
  type      = "Microsoft.Maintenance/publicMaintenanceConfigurations@2023-04-01"
  parent_id = "/subscriptions/${data.azurerm_client_config.current.subscription_id}"
  name      = "SQL_Default"
}

resource "azapi_resource" "database" {
  type      = "Microsoft.Sql/servers/databases@2021-02-01-preview"
  parent_id = azapi_resource.server.id
  name      = var.resource_name
  location  = var.location
  body = {
    properties = {
      autoPauseDelay                   = 0
      collation                        = "SQL_Latin1_General_CP1_CI_AS"
      createMode                       = "Default"
      elasticPoolId                    = ""
      highAvailabilityReplicaCount     = 0
      isLedgerOn                       = false
      maintenanceConfigurationId       = data.azapi_resource_id.publicMaintenanceConfiguration.id
      minCapacity                      = 0
      readScale                        = "Disabled"
      requestedBackupStorageRedundancy = "Geo"
      zoneRedundant                    = false
    }
  }
  schema_validation_enabled = false
  response_export_values    = ["*"]
}

resource "azapi_resource" "jobAgent" {
  type      = "Microsoft.Sql/servers/jobAgents@2020-11-01-preview"
  parent_id = azapi_resource.server.id
  name      = var.resource_name
  location  = var.location
  body = {
    properties = {
      databaseId = azapi_resource.database.id
    }
  }
  schema_validation_enabled = false
  response_export_values    = ["*"]
}


```



## Arguments Reference

The following arguments are supported:

* `type` - (Required) The type of the resource. This should be set to `Microsoft.Sql/servers/jobAgents@api-version`. The available api-versions for this resource are: [`2017-03-01-preview`, `2020-02-02-preview`, `2020-08-01-preview`, `2020-11-01-preview`, `2021-02-01-preview`, `2021-05-01-preview`, `2021-08-01-preview`, `2021-11-01`, `2021-11-01-preview`, `2022-02-01-preview`, `2022-05-01-preview`, `2022-08-01-preview`, `2022-11-01-preview`, `2023-02-01-preview`, `2023-05-01-preview`, `2023-08-01`, `2023-08-01-preview`, `2024-05-01-preview`, `2024-11-01-preview`].

* `parent_id` - (Required) The ID of the azure resource in which this resource is created. The allowed values are:  
  `/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Sql/servers/{resourceName}`

* `name` - (Required) Specifies the name of the azure resource. Changing this forces a new resource to be created.

* `body` - (Required) Specifies the configuration of the resource. More information about the arguments in `body` can be found in the [Microsoft documentation](https://learn.microsoft.com/en-us/azure/templates/Microsoft.Sql/servers/jobAgents?pivots=deployment-language-terraform).

For other arguments, please refer to the [azapi_resource](https://registry.terraform.io/providers/Azure/azapi/latest/docs/resources/resource) documentation.

## Import

 ```shell
 # Azure resource can be imported using the resource id, e.g.
 terraform import azapi_resource.example /subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Sql/servers/{resourceName}/jobAgents/{resourceName}
 
 # It also supports specifying API version by using the resource id with api-version as a query parameter, e.g.
 terraform import azapi_resource.example /subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Sql/servers/{resourceName}/jobAgents/{resourceName}?api-version=2024-11-01-preview
 ```
