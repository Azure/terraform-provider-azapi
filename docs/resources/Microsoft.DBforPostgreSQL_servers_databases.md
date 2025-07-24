---
subcategory: "Microsoft.DBforPostgreSQL - Azure Database for PostgreSQL"
page_title: "servers/databases"
description: |-
  Manages a PostgreSQL Database within a PostgreSQL Server.
---

# Microsoft.DBforPostgreSQL/servers/databases - PostgreSQL Database within a PostgreSQL Server

This article demonstrates how to use `azapi` provider to manage the PostgreSQL Database within a PostgreSQL Server resource in Azure.

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
  default = "westeurope"
}

variable "administrator_login" {
  type        = string
  description = "The administrator login for the PostgreSQL server"
}

variable "administrator_login_password" {
  type        = string
  description = "The administrator login password for the PostgreSQL server"
  sensitive   = true
}

resource "azapi_resource" "resourceGroup" {
  type     = "Microsoft.Resources/resourceGroups@2020-06-01"
  name     = var.resource_name
  location = var.location
}

resource "azapi_resource" "server" {
  type      = "Microsoft.DBforPostgreSQL/servers@2017-12-01"
  parent_id = azapi_resource.resourceGroup.id
  name      = var.resource_name
  location  = var.location
  body = {
    properties = {
      administratorLogin         = var.administrator_login
      administratorLoginPassword = var.administrator_login_password
      createMode                 = "Default"
      infrastructureEncryption   = "Disabled"
      minimalTlsVersion          = "TLS1_2"
      publicNetworkAccess        = "Enabled"
      sslEnforcement             = "Enabled"
      storageProfile = {
        backupRetentionDays = 7
        storageAutogrow     = "Enabled"
        storageMB           = 51200
      }
      version = "9.6"
    }
    sku = {
      capacity = 2
      family   = "Gen5"
      name     = "GP_Gen5_2"
      tier     = "GeneralPurpose"
    }
  }
  schema_validation_enabled = false
  response_export_values    = ["*"]
}

resource "azapi_resource" "database" {
  type      = "Microsoft.DBforPostgreSQL/servers/databases@2017-12-01"
  parent_id = azapi_resource.server.id
  name      = var.resource_name
  body = {
    properties = {
      charset   = "UTF8"
      collation = "English_United States.1252"
    }
  }
  schema_validation_enabled = false
  response_export_values    = ["*"]
}


```



## Arguments Reference

The following arguments are supported:

* `type` - (Required) The type of the resource. This should be set to `Microsoft.DBforPostgreSQL/servers/databases@api-version`. The available api-versions for this resource are: [`2017-12-01`, `2017-12-01-preview`].

* `parent_id` - (Required) The ID of the azure resource in which this resource is created. The allowed values are:  
  `/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.DBforPostgreSQL/servers/{resourceName}`

* `name` - (Required) Specifies the name of the azure resource. Changing this forces a new resource to be created.

* `body` - (Required) Specifies the configuration of the resource. More information about the arguments in `body` can be found in the [Microsoft documentation](https://learn.microsoft.com/en-us/azure/templates/Microsoft.DBforPostgreSQL/servers/databases?pivots=deployment-language-terraform).

For other arguments, please refer to the [azapi_resource](https://registry.terraform.io/providers/Azure/azapi/latest/docs/resources/resource) documentation.

## Import

 ```shell
 # Azure resource can be imported using the resource id, e.g.
 terraform import azapi_resource.example /subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.DBforPostgreSQL/servers/{resourceName}/databases/{resourceName}
 
 # It also supports specifying API version by using the resource id with api-version as a query parameter, e.g.
 terraform import azapi_resource.example /subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.DBforPostgreSQL/servers/{resourceName}/databases/{resourceName}?api-version=2017-12-01-preview
 ```
