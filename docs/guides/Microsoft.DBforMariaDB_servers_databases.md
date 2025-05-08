---
subcategory: "Microsoft.DBforMariaDB - Azure Database for MariaDB"
page_title: "servers/databases"
description: |-
  Manages a Database for MariaDB Servers Databases.
---

# Microsoft.DBforMariaDB/servers/databases - Database for MariaDB Servers Databases

This article demonstrates how to use `azapi` provider to manage the Database for MariaDB Servers Databases resource in Azure.

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

resource "azapi_resource" "resourceGroup" {
  type     = "Microsoft.Resources/resourceGroups@2020-06-01"
  name     = var.resource_name
  location = var.location
}

resource "azapi_resource" "server" {
  type      = "Microsoft.DBforMariaDB/servers@2018-06-01"
  parent_id = azapi_resource.resourceGroup.id
  name      = var.resource_name
  location  = var.location
  body = {
    properties = {
      administratorLogin         = "acctestun"
      administratorLoginPassword = "H@Sh1CoR3!"
      createMode                 = "Default"
      minimalTlsVersion          = "TLS1_2"
      publicNetworkAccess        = "Enabled"
      sslEnforcement             = "Enabled"
      storageProfile = {
        backupRetentionDays = 7
        storageAutogrow     = "Enabled"
        storageMB           = 51200
      }
      version = "10.2"
    }
    sku = {
      capacity = 2
      family   = "Gen5"
      name     = "B_Gen5_2"
      tier     = "Basic"
    }
  }
  schema_validation_enabled = false
  response_export_values    = ["*"]
}

resource "azapi_resource" "database" {
  type      = "Microsoft.DBforMariaDB/servers/databases@2018-06-01"
  parent_id = azapi_resource.server.id
  name      = var.resource_name
  body = {
    properties = {
      charset   = "utf8"
      collation = "utf8_general_ci"
    }
  }
  schema_validation_enabled = false
  response_export_values    = ["*"]
}


```



## Arguments Reference

The following arguments are supported:

* `type` - (Required) The type of the resource. This should be set to `Microsoft.DBforMariaDB/servers/databases@api-version`. The available api-versions for this resource are: [`2018-06-01`, `2018-06-01-preview`].

* `parent_id` - (Required) The ID of the azure resource in which this resource is created. The allowed values are:  
  `/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.DBforMariaDB/servers/{resourceName}`

* `name` - (Required) Specifies the name of the azure resource. Changing this forces a new resource to be created.

* `body` - (Required) Specifies the configuration of the resource. More information about the arguments in `body` can be found in the [Microsoft documentation](https://learn.microsoft.com/en-us/azure/templates/Microsoft.DBforMariaDB/servers/databases?pivots=deployment-language-terraform).

## Import

 ```shell
 # Azure resource can be imported using the resource id, e.g.
 terraform import azapi_resource.example /subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.DBforMariaDB/servers/{resourceName}/databases/{resourceName}
 
 # It also supports specifying API version by using the resource id with api-version as a query parameter, e.g.
 terraform import azapi_resource.example /subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.DBforMariaDB/servers/{resourceName}/databases/{resourceName}?api-version=2018-06-01-preview
 ```
