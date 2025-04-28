---
subcategory: "Microsoft.DBforMySQL - Azure Database for MySQL"
page_title: "flexibleServers/firewallRules"
description: |-
  Manages a Firewall Rule for a MySQL Flexible Server.
---

# Microsoft.DBforMySQL/flexibleServers/firewallRules - Firewall Rule for a MySQL Flexible Server

This article demonstrates how to use `azapi` provider to manage the Firewall Rule for a MySQL Flexible Server resource in Azure.

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

resource "azapi_resource" "flexibleServer" {
  type      = "Microsoft.DBforMySQL/flexibleServers@2021-05-01"
  parent_id = azapi_resource.resourceGroup.id
  name      = var.resource_name
  location  = var.location
  body = {
    properties = {
      administratorLogin         = "adminTerraform"
      administratorLoginPassword = "QAZwsx123"
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
      version = "5.7"
    }
    sku = {
      name = "Standard_B1s"
      tier = "Burstable"
    }
  }
  schema_validation_enabled = false
  response_export_values    = ["*"]
}

resource "azapi_resource" "firewallRule" {
  type      = "Microsoft.DBforMySQL/flexibleServers/firewallRules@2021-05-01"
  parent_id = azapi_resource.flexibleServer.id
  name      = var.resource_name
  body = {
    properties = {
      endIpAddress   = "255.255.255.255"
      startIpAddress = "0.0.0.0"
    }
  }
  schema_validation_enabled = false
  response_export_values    = ["*"]
}


```



## Arguments Reference

The following arguments are supported:

* `type` - (Required) The type of the resource. This should be set to `Microsoft.DBforMySQL/flexibleServers/firewallRules@api-version`. The available api-versions for this resource are: [`2020-07-01-preview`, `2021-05-01`, `2021-05-01-preview`, `2021-12-01-preview`, `2022-01-01`, `2023-06-01-preview`, `2023-06-30`, `2023-12-30`].

* `parent_id` - (Required) The ID of the azure resource in which this resource is created. The allowed values are:  
  `/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.DBforMySQL/flexibleServers/{resourceName}`

* `name` - (Required) Specifies the name of the azure resource. Changing this forces a new resource to be created.

* `body` - (Required) Specifies the configuration of the resource. More information about the arguments in `body` can be found in the [Microsoft documentation](https://learn.microsoft.com/en-us/azure/templates/Microsoft.DBforMySQL/flexibleServers/firewallRules?pivots=deployment-language-terraform).

## Import

 ```shell
 # Azure resource can be imported using the resource id, e.g.
 terraform import azapi_resource.example /subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.DBforMySQL/flexibleServers/{resourceName}/firewallRules/{resourceName}
 
 # It also supports specifying API version by using the resource id with api-version as a query parameter, e.g.
 terraform import azapi_resource.example /subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.DBforMySQL/flexibleServers/{resourceName}/firewallRules/{resourceName}?api-version=2023-12-30
 ```
