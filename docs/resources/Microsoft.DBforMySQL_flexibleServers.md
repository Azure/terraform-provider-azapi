---
subcategory: "Microsoft.DBforMySQL - Azure Database for MySQL"
page_title: "flexibleServers"
description: |-
  Manages a MySQL Flexible Server.
---

# Microsoft.DBforMySQL/flexibleServers - MySQL Flexible Server

This article demonstrates how to use `azapi` provider to manage the MySQL Flexible Server resource in Azure.

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
  description = "The administrator login name for the MySQL flexible server"
}

variable "administrator_login_password" {
  type        = string
  description = "The administrator login password for the MySQL flexible server"
  sensitive   = true
}

resource "azapi_resource" "resourceGroup" {
  type     = "Microsoft.Resources/resourceGroups@2020-06-01"
  name     = var.resource_name
  location = var.location
}

resource "azapi_resource" "flexibleServer" {
  type      = "Microsoft.DBforMySQL/flexibleServers@2021-05-01"
  parent_id = azapi_resource.resourceGroup.id
  name      = var.resource_name
  location  = var.location
  body = {
    properties = {
      administratorLogin         = var.administrator_login
      administratorLoginPassword = var.administrator_login_password
      backup = {
        backupRetentionDays = 7
        geoRedundantBackup  = "Disabled"
      }
      createMode = ""
      dataEncryption = {
        type = "SystemManaged"
      }
      highAvailability = {
        mode = "Disabled"
      }
      network = {
      }
    }
    sku = {
      name = "Standard_B1s"
      tier = "Burstable"
    }
  }
  schema_validation_enabled = false
  response_export_values    = ["*"]
}


```



## Arguments Reference

The following arguments are supported:

* `type` - (Required) The type of the resource. This should be set to `Microsoft.DBforMySQL/flexibleServers@api-version`. The available api-versions for this resource are: [`2020-07-01-preview`, `2021-05-01`, `2021-05-01-preview`, `2021-12-01-preview`, `2022-01-01`, `2022-09-30-preview`, `2023-06-01-preview`, `2023-06-30`, `2023-10-01-preview`, `2023-12-01-preview`, `2023-12-30`, `2024-02-01-preview`, `2024-06-01-preview`, `2024-10-01-preview`, `2024-12-01-preview`].

* `parent_id` - (Required) The ID of the azure resource in which this resource is created. The allowed values are:  
  `/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}`

* `name` - (Required) Specifies the name of the azure resource. Changing this forces a new resource to be created.

* `body` - (Required) Specifies the configuration of the resource. More information about the arguments in `body` can be found in the [Microsoft documentation](https://learn.microsoft.com/en-us/azure/templates/Microsoft.DBforMySQL/flexibleServers?pivots=deployment-language-terraform).

For other arguments, please refer to the [azapi_resource](https://registry.terraform.io/providers/Azure/azapi/latest/docs/resources/resource) documentation.

## Import

 ```shell
 # Azure resource can be imported using the resource id, e.g.
 terraform import azapi_resource.example /subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.DBforMySQL/flexibleServers/{resourceName}
 
 # It also supports specifying API version by using the resource id with api-version as a query parameter, e.g.
 terraform import azapi_resource.example /subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.DBforMySQL/flexibleServers/{resourceName}?api-version=2024-12-01-preview
 ```
