---
subcategory: "Microsoft.DBforPostgreSQL - Azure Database for PostgreSQL"
page_title: "flexibleServers/virtualEndpoints"
description: |-
  Manages a Virtual Endpoint on a PostgreSQL Flexible Server.
---

# Microsoft.DBforPostgreSQL/flexibleServers/virtualEndpoints - Virtual Endpoint on a PostgreSQL Flexible Server

This article demonstrates how to use `azapi` provider to manage the Virtual Endpoint on a PostgreSQL Flexible Server resource in Azure.



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

resource "azapi_resource" "resourceGroup" {
  type     = "Microsoft.Resources/resourceGroups@2020-06-01"
  name     = var.resource_name
  location = var.location
}

resource "azapi_resource" "flexibleServer" {
  type      = "Microsoft.DBforPostgreSQL/flexibleServers@2024-08-01"
  parent_id = azapi_resource.resourceGroup.id
  name      = "${var.resource_name}-primary"
  location  = var.location
  body = {
    properties = {
      administratorLogin         = "psqladmin"
      administratorLoginPassword = "H@Sh1CoR3!"
      availabilityZone           = "1"
      backup = {
        geoRedundantBackup = "Disabled"
      }
      highAvailability = {
        mode = "Disabled"
      }
      network = {
        publicNetworkAccess = "Disabled"
      }
      storage = {
        autoGrow      = "Disabled"
        storageSizeGB = 32
        tier          = "P30"
      }
      version = "16"
    }
    sku = {
      name = "Standard_D2ads_v5"
      tier = "GeneralPurpose"
    }
  }
}

resource "azapi_resource" "flexibleServer_1" {
  type      = "Microsoft.DBforPostgreSQL/flexibleServers@2024-08-01"
  parent_id = azapi_resource.resourceGroup.id
  name      = "${var.resource_name}-replica"
  location  = var.location
  body = {
    properties = {
      availabilityZone = "1"
      backup = {
        geoRedundantBackup = "Disabled"
      }
      createMode = "Replica"
      highAvailability = {
        mode = "Disabled"
      }
      network = {
        publicNetworkAccess = "Disabled"
      }
      sourceServerResourceId = azapi_resource.flexibleServer.id
      storage = {
        autoGrow      = "Disabled"
        storageSizeGB = 32
        tier          = "P30"
      }
      version = "16"
    }
  }
}

resource "azapi_resource" "virtualEndpoint" {
  type      = "Microsoft.DBforPostgreSQL/flexibleServers/virtualEndpoints@2024-08-01"
  parent_id = azapi_resource.flexibleServer.id
  name      = var.resource_name
  body = {
    properties = {
      endpointType = "ReadWrite"
      members      = [azapi_resource.flexibleServer_1.name]
    }
  }
}

```



## Arguments Reference

The following arguments are supported:

* `type` - (Required) The type of the resource. This should be set to `Microsoft.DBforPostgreSQL/flexibleServers/virtualEndpoints@api-version`. The available api-versions for this resource are: [`2023-06-01-preview`, `2023-12-01-preview`, `2024-03-01-preview`, `2024-08-01`, `2024-11-01-preview`, `2025-01-01-preview`, `2025-06-01-preview`].

* `parent_id` - (Required) The ID of the azure resource in which this resource is created. The allowed values are:  
  `/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.DBforPostgreSQL/flexibleServers/{resourceName}`

* `name` - (Required) Specifies the name of the azure resource. Changing this forces a new resource to be created.

* `body` - (Required) Specifies the configuration of the resource. More information about the arguments in `body` can be found in the [Microsoft documentation](https://learn.microsoft.com/en-us/azure/templates/Microsoft.DBforPostgreSQL/flexibleServers/virtualEndpoints?pivots=deployment-language-terraform).

For other arguments, please refer to the [azapi_resource](https://registry.terraform.io/providers/Azure/azapi/latest/docs/resources/resource) documentation.

## Import

 ```shell
 # Azure resource can be imported using the resource id, e.g.
 terraform import azapi_resource.example /subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.DBforPostgreSQL/flexibleServers/{resourceName}/virtualEndpoints/{resourceName}
 
 # It also supports specifying API version by using the resource id with api-version as a query parameter, e.g.
 terraform import azapi_resource.example /subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.DBforPostgreSQL/flexibleServers/{resourceName}/virtualEndpoints/{resourceName}?api-version=2025-06-01-preview
 ```
