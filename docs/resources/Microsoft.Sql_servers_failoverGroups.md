---
subcategory: "Microsoft.Sql - Azure SQL Database, Azure SQL Managed Instance, Azure Synapse Analytics"
page_title: "servers/failoverGroups"
description: |-
  Manages a Microsoft Azure SQL Failover Group.
---

# Microsoft.Sql/servers/failoverGroups - Microsoft Azure SQL Failover Group

This article demonstrates how to use `azapi` provider to manage the Microsoft Azure SQL Failover Group resource in Azure.

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

variable "secondary_location" {
  type    = string
  default = "eastus"
}

resource "azapi_resource" "resourceGroup" {
  type     = "Microsoft.Resources/resourceGroups@2020-06-01"
  name     = var.resource_name
  location = var.location
}

resource "azapi_resource" "server" {
  type      = "Microsoft.Sql/servers@2023-08-01-preview"
  parent_id = azapi_resource.resourceGroup.id
  name      = "${var.resource_name}-primary"
  location  = var.location
  body = {
    properties = {
      administratorLogin            = "mradministrator"
      administratorLoginPassword    = "thisIsDog11"
      minimalTlsVersion             = "1.2"
      publicNetworkAccess           = "Enabled"
      restrictOutboundNetworkAccess = "Disabled"
      version                       = "12.0"
    }
  }
}

resource "azapi_resource" "server_1" {
  type      = "Microsoft.Sql/servers@2023-08-01-preview"
  parent_id = azapi_resource.resourceGroup.id
  name      = "${var.resource_name}-secondary"
  location  = var.secondary_location
  body = {
    properties = {
      administratorLogin            = "mradministrator"
      administratorLoginPassword    = "thisIsDog11"
      minimalTlsVersion             = "1.2"
      publicNetworkAccess           = "Enabled"
      restrictOutboundNetworkAccess = "Disabled"
      version                       = "12.0"
    }
  }
}

resource "azapi_resource" "database" {
  type      = "Microsoft.Sql/servers/databases@2023-08-01-preview"
  parent_id = azapi_resource.server.id
  name      = "${var.resource_name}-db"
  location  = var.location
  body = {
    properties = {
      autoPauseDelay                   = 0
      collation                        = "SQL_Latin1_General_CP1_CI_AS"
      createMode                       = "Default"
      elasticPoolId                    = ""
      encryptionProtectorAutoRotation  = false
      highAvailabilityReplicaCount     = 0
      isLedgerOn                       = false
      licenseType                      = ""
      maxSizeBytes                     = 214748364800
      minCapacity                      = 0
      readScale                        = "Disabled"
      requestedBackupStorageRedundancy = "Geo"
      sampleName                       = ""
      secondaryType                    = ""
      zoneRedundant                    = false
    }
    sku = {
      name = "S1"
    }
  }
}

resource "azapi_resource" "failoverGroup" {
  type      = "Microsoft.Sql/servers/failoverGroups@2023-08-01-preview"
  parent_id = azapi_resource.server.id
  name      = "${var.resource_name}-fg"
  body = {
    properties = {
      databases = [azapi_resource.database.id]
      partnerServers = [{
        id = azapi_resource.server_1.id
      }]
      readOnlyEndpoint = {
        failoverPolicy = "Disabled"
      }
      readWriteEndpoint = {
        failoverPolicy                         = "Automatic"
        failoverWithDataLossGracePeriodMinutes = 60
      }
    }
  }
}


```



## Arguments Reference

The following arguments are supported:

* `type` - (Required) The type of the resource. This should be set to `Microsoft.Sql/servers/failoverGroups@api-version`. The available api-versions for this resource are: [`2015-05-01-preview`, `2020-02-02-preview`, `2020-08-01-preview`, `2020-11-01-preview`, `2021-02-01-preview`, `2021-05-01-preview`, `2021-08-01-preview`, `2021-11-01`, `2021-11-01-preview`, `2022-02-01-preview`, `2022-05-01-preview`, `2022-08-01-preview`, `2022-11-01-preview`, `2023-02-01-preview`, `2023-05-01-preview`, `2023-08-01`, `2023-08-01-preview`, `2024-05-01-preview`, `2024-11-01-preview`].

* `parent_id` - (Required) The ID of the azure resource in which this resource is created. The allowed values are:  
  `/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Sql/servers/{resourceName}`

* `name` - (Required) Specifies the name of the azure resource. Changing this forces a new resource to be created.

* `body` - (Required) Specifies the configuration of the resource. More information about the arguments in `body` can be found in the [Microsoft documentation](https://learn.microsoft.com/en-us/azure/templates/Microsoft.Sql/servers/failoverGroups?pivots=deployment-language-terraform).

For other arguments, please refer to the [azapi_resource](https://registry.terraform.io/providers/Azure/azapi/latest/docs/resources/resource) documentation.

## Import

 ```shell
 # Azure resource can be imported using the resource id, e.g.
 terraform import azapi_resource.example /subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Sql/servers/{resourceName}/failoverGroups/{resourceName}
 
 # It also supports specifying API version by using the resource id with api-version as a query parameter, e.g.
 terraform import azapi_resource.example /subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Sql/servers/{resourceName}/failoverGroups/{resourceName}?api-version=2024-11-01-preview
 ```
