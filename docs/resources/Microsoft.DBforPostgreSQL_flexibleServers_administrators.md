---
subcategory: "Microsoft.DBforPostgreSQL - Azure Database for PostgreSQL"
page_title: "flexibleServers/administrators"
description: |-
  Manages a Active Directory administrator on a PostgreSQL Flexible server.
---

# Microsoft.DBforPostgreSQL/flexibleServers/administrators - Active Directory administrator on a PostgreSQL Flexible server

This article demonstrates how to use `azapi` provider to manage the Active Directory administrator on a PostgreSQL Flexible server resource in Azure.



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

variable "administrator_login" {
  type        = string
  description = "The administrator login for the PostgreSQL flexible server"
}

variable "administrator_login_password" {
  type        = string
  description = "The administrator login password for the PostgreSQL flexible server"
  sensitive   = true
}

data "azurerm_client_config" "current" {
}

resource "azapi_resource" "resourceGroup" {
  type     = "Microsoft.Resources/resourceGroups@2020-06-01"
  name     = var.resource_name
  location = var.location
}

resource "azapi_resource" "flexibleServer" {
  type      = "Microsoft.DBforPostgreSQL/flexibleServers@2022-12-01"
  parent_id = azapi_resource.resourceGroup.id
  name      = var.resource_name
  location  = var.location
  body = {
    properties = {
      administratorLogin         = var.administrator_login
      administratorLoginPassword = var.administrator_login_password
      authConfig = {
        activeDirectoryAuth = "Enabled"
        passwordAuth        = "Enabled"
        tenantId            = data.azurerm_client_config.current.tenant_id
      }
      availabilityZone = "2"
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

resource "azapi_resource" "administrator" {
  type      = "Microsoft.DBforPostgreSQL/flexibleServers/administrators@2022-12-01"
  parent_id = azapi_resource.flexibleServer.id
  name      = data.azurerm_client_config.current.object_id
  body = {
    properties = {
      principalType = "ServicePrincipal"
      tenantId      = data.azurerm_client_config.current.tenant_id
    }
  }
  schema_validation_enabled = false
  response_export_values    = ["*"]
}


```



## Arguments Reference

The following arguments are supported:

* `type` - (Required) The type of the resource. This should be set to `Microsoft.DBforPostgreSQL/flexibleServers/administrators@api-version`. The available api-versions for this resource are: [`2022-03-08-preview`, `2022-12-01`, `2023-03-01-preview`, `2023-06-01-preview`, `2023-12-01-preview`, `2024-03-01-preview`, `2024-08-01`, `2024-11-01-preview`, `2025-01-01-preview`, `2025-06-01-preview`].

* `parent_id` - (Required) The ID of the azure resource in which this resource is created. The allowed values are:  
  `/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.DBforPostgreSQL/flexibleServers/{resourceName}`

* `name` - (Required) Specifies the name of the azure resource. Changing this forces a new resource to be created.

* `body` - (Required) Specifies the configuration of the resource. More information about the arguments in `body` can be found in the [Microsoft documentation](https://learn.microsoft.com/en-us/azure/templates/Microsoft.DBforPostgreSQL/flexibleServers/administrators?pivots=deployment-language-terraform).

For other arguments, please refer to the [azapi_resource](https://registry.terraform.io/providers/Azure/azapi/latest/docs/resources/resource) documentation.

## Import

 ```shell
 # Azure resource can be imported using the resource id, e.g.
 terraform import azapi_resource.example /subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.DBforPostgreSQL/flexibleServers/{resourceName}/administrators/{resourceName}
 
 # It also supports specifying API version by using the resource id with api-version as a query parameter, e.g.
 terraform import azapi_resource.example /subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.DBforPostgreSQL/flexibleServers/{resourceName}/administrators/{resourceName}?api-version=2025-06-01-preview
 ```
