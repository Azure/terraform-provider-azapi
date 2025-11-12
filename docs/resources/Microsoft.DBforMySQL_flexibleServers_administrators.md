---
subcategory: "Microsoft.DBforMySQL - Azure Database for MySQL"
page_title: "flexibleServers/administrators"
description: |-
  Manages a Active Directory administrator on a MySQL Flexible Server.
---

# Microsoft.DBforMySQL/flexibleServers/administrators - Active Directory administrator on a MySQL Flexible Server

This article demonstrates how to use `azapi` provider to manage the Active Directory administrator on a MySQL Flexible Server resource in Azure.



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

data "azapi_client_config" "current" {}

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
  description = "The administrator login password for the MySQL flexible server"
}

resource "azapi_resource" "resourceGroup" {
  type     = "Microsoft.Resources/resourceGroups@2020-06-01"
  name     = var.resource_name
  location = var.location
}

resource "azapi_resource" "userAssignedIdentity" {
  type      = "Microsoft.ManagedIdentity/userAssignedIdentities@2023-01-31"
  parent_id = azapi_resource.resourceGroup.id
  name      = "${var.resource_name}-uai1"
  location  = var.location
}

resource "azapi_resource" "userAssignedIdentity_1" {
  type      = "Microsoft.ManagedIdentity/userAssignedIdentities@2023-01-31"
  parent_id = azapi_resource.resourceGroup.id
  name      = "${var.resource_name}-uai2"
  location  = var.location
}

resource "azapi_resource" "flexibleServer" {
  type      = "Microsoft.DBforMySQL/flexibleServers@2023-12-30"
  parent_id = azapi_resource.resourceGroup.id
  name      = "${var.resource_name}-mysql"
  location  = var.location
  identity {
    type         = "UserAssigned"
    identity_ids = [azapi_resource.userAssignedIdentity.id]
  }
  body = {
    properties = {
      administratorLogin         = "tfadmin"
      administratorLoginPassword = var.administrator_login_password
      backup = {
        backupRetentionDays = 7
        geoRedundantBackup  = "Disabled"
      }
      dataEncryption = {
        type = "SystemManaged"
      }
      highAvailability = {
        mode = "Disabled"
      }
      version = "8.0.21"
    }
    sku = {
      name = "Standard_B1ms"
      tier = "Burstable"
    }
  }
}

resource "azapi_resource" "administrator" {
  type      = "Microsoft.DBforMySQL/flexibleServers/administrators@2023-12-30"
  parent_id = azapi_resource.flexibleServer.id
  name      = "ActiveDirectory"
  body = {
    properties = {
      administratorType  = "ActiveDirectory"
      identityResourceId = azapi_resource.userAssignedIdentity.id
      login              = "sqladmin"
      sid                = data.azapi_client_config.current.object_id
      tenantId           = data.azapi_client_config.current.tenant_id
    }
  }
}


```



## Arguments Reference

The following arguments are supported:

* `type` - (Required) The type of the resource. This should be set to `Microsoft.DBforMySQL/flexibleServers/administrators@api-version`. The available api-versions for this resource are: [`2021-12-01-preview`, `2022-01-01`, `2023-06-01-preview`, `2023-06-30`, `2023-12-30`, `2024-12-01-preview`, `2024-12-30`].

* `parent_id` - (Required) The ID of the azure resource in which this resource is created. The allowed values are:  
  `/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.DBforMySQL/flexibleServers/{resourceName}`

* `name` - (Required) Specifies the name of the azure resource. Changing this forces a new resource to be created.

* `body` - (Required) Specifies the configuration of the resource. More information about the arguments in `body` can be found in the [Microsoft documentation](https://learn.microsoft.com/en-us/azure/templates/Microsoft.DBforMySQL/flexibleServers/administrators?pivots=deployment-language-terraform).

For other arguments, please refer to the [azapi_resource](https://registry.terraform.io/providers/Azure/azapi/latest/docs/resources/resource) documentation.

## Import

 ```shell
 # Azure resource can be imported using the resource id, e.g.
 terraform import azapi_resource.example /subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.DBforMySQL/flexibleServers/{resourceName}/administrators/{resourceName}
 
 # It also supports specifying API version by using the resource id with api-version as a query parameter, e.g.
 terraform import azapi_resource.example /subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.DBforMySQL/flexibleServers/{resourceName}/administrators/{resourceName}?api-version=2024-12-30
 ```
