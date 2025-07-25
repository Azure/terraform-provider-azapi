---
subcategory: "Microsoft.DBforPostgreSQL - Azure Database for PostgreSQL"
page_title: "flexibleServers/configurations"
description: |-
  Manages a Sets a PostgreSQL Configuration value on a Azure PostgreSQL Flexible Server.
---

# Microsoft.DBforPostgreSQL/flexibleServers/configurations - Sets a PostgreSQL Configuration value on a Azure PostgreSQL Flexible Server

This article demonstrates how to use `azapi` provider to manage the Sets a PostgreSQL Configuration value on a Azure PostgreSQL Flexible Server resource in Azure.

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
  default = "eastus"
}

variable "postgresql_administrator_password" {
  type        = string
  description = "The administrator password for the PostgreSQL flexible server"
  sensitive   = true
}

resource "azapi_resource" "resourceGroup" {
  type     = "Microsoft.Resources/resourceGroups@2020-06-01"
  name     = var.resource_name
  location = var.location
}

resource "azapi_resource" "flexibleServer" {
  type      = "Microsoft.DBforPostgreSQL/flexibleServers@2023-06-01-preview"
  parent_id = azapi_resource.resourceGroup.id
  name      = var.resource_name
  location  = var.location
  body = {
    identity = {
      type                   = "None",
      userAssignedIdentities = null
    },
    properties = {
      administratorLogin         = "adminTerraform"
      administratorLoginPassword = var.postgresql_administrator_password
      availabilityZone           = "2"
      backup = {
        geoRedundantBackup = "Disabled"
      }
      highAvailability = {
        mode = "Disabled"
      }
      network = {
      }
      storage = {
        storageSizeGB = 32
      }
      version = "12"
    }
    sku = {
      name = "Standard_D2s_v3"
      tier = "GeneralPurpose"
    }
  }
  schema_validation_enabled = false
  response_export_values    = ["*"]
}

resource "azapi_update_resource" "pgbouncerEnabled" {
  type      = "Microsoft.DBforPostgreSQL/flexibleServers/configurations@2022-12-01"
  parent_id = azapi_resource.flexibleServer.id
  name      = "pgbouncer.enabled"
  body = {
    properties = {
      value  = "true"
      source = "user-override"
    }
  }
}

resource "azapi_update_resource" "pgbouncerDefaultPoolSize" {
  depends_on = [azapi_update_resource.pgbouncerEnabled]
  type       = "Microsoft.DBforPostgreSQL/flexibleServers/configurations@2022-12-01"
  parent_id  = azapi_resource.flexibleServer.id
  name       = "pgbouncer.default_pool_size"
  body = {
    properties = {
      value  = "40"
      source = "user-override"
    }
  }
}

```



## Arguments Reference

The following arguments are supported:

* `type` - (Required) The type of the resource. This should be set to `Microsoft.DBforPostgreSQL/flexibleServers/configurations@api-version`. The available api-versions for this resource are: [`2020-02-14-preview`, `2021-06-01`, `2021-06-01-preview`, `2022-01-20-preview`, `2022-03-08-preview`, `2022-12-01`, `2023-03-01-preview`, `2023-06-01-preview`, `2023-12-01-preview`, `2024-03-01-preview`, `2024-08-01`, `2024-11-01-preview`, `2025-01-01-preview`].

* `parent_id` - (Required) The ID of the azure resource in which this resource is created. The allowed values are:  
  `/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.DBforPostgreSQL/flexibleServers/{resourceName}`

* `name` - (Required) Specifies the name of the azure resource. Changing this forces a new resource to be created.

* `body` - (Required) Specifies the configuration of the resource. More information about the arguments in `body` can be found in the [Microsoft documentation](https://learn.microsoft.com/en-us/azure/templates/Microsoft.DBforPostgreSQL/flexibleServers/configurations?pivots=deployment-language-terraform).

For other arguments, please refer to the [azapi_resource](https://registry.terraform.io/providers/Azure/azapi/latest/docs/resources/resource) documentation.

## Import

 ```shell
 # Azure resource can be imported using the resource id, e.g.
 terraform import azapi_resource.example /subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.DBforPostgreSQL/flexibleServers/{resourceName}/configurations/{resourceName}
 
 # It also supports specifying API version by using the resource id with api-version as a query parameter, e.g.
 terraform import azapi_resource.example /subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.DBforPostgreSQL/flexibleServers/{resourceName}/configurations/{resourceName}?api-version=2025-01-01-preview
 ```
