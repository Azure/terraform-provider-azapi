---
subcategory: "Microsoft.DBforMariaDB - Azure Database for MariaDB"
page_title: "servers/configurations"
description: |-
  Manages a Database for MariaDB Servers Configurations.
---

# Microsoft.DBforMariaDB/servers/configurations - Database for MariaDB Servers Configurations

This article demonstrates how to use `azapi` provider to manage the Database for MariaDB Servers Configurations resource in Azure.

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
      name     = "GP_Gen5_2"
      tier     = "GeneralPurpose"
    }
  }
  schema_validation_enabled = false
  response_export_values    = ["*"]
}

resource "azapi_update_resource" "configuration" {
  type      = "Microsoft.DBforMariaDB/servers/configurations@2018-06-01"
  parent_id = azapi_resource.server.id
  name      = "character_set_server"
  body = {
    properties = {
      value = "LATIN1"
    }
  }
  response_export_values = ["*"]
}


```



## Arguments Reference

The following arguments are supported:

* `type` - (Required) The type of the resource. This should be set to `Microsoft.DBforMariaDB/servers/configurations@api-version`. The available api-versions for this resource are: [`2018-06-01`, `2018-06-01-preview`].

* `parent_id` - (Required) The ID of the azure resource in which this resource is created. The allowed values are:  
  `/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.DBforMariaDB/servers/{resourceName}`

* `name` - (Required) Specifies the name of the azure resource. Changing this forces a new resource to be created.

* `body` - (Required) Specifies the configuration of the resource. More information about the arguments in `body` can be found in the [Microsoft documentation](https://learn.microsoft.com/en-us/azure/templates/Microsoft.DBforMariaDB/servers/configurations?pivots=deployment-language-terraform).

For other arguments, please refer to the [azapi_resource](https://registry.terraform.io/providers/Azure/azapi/latest/docs/resources/resource) documentation.

## Import

 ```shell
 # Azure resource can be imported using the resource id, e.g.
 terraform import azapi_resource.example /subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.DBforMariaDB/servers/{resourceName}/configurations/{resourceName}
 
 # It also supports specifying API version by using the resource id with api-version as a query parameter, e.g.
 terraform import azapi_resource.example /subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.DBforMariaDB/servers/{resourceName}/configurations/{resourceName}?api-version=2018-06-01-preview
 ```
