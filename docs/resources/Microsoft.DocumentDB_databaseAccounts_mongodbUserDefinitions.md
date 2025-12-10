---
subcategory: "Microsoft.DocumentDB - Azure Cosmos DB"
page_title: "databaseAccounts/mongodbUserDefinitions"
description: |-
  Manages a Cosmos DB Mongo User Definition.
---

# Microsoft.DocumentDB/databaseAccounts/mongodbUserDefinitions - Cosmos DB Mongo User Definition

This article demonstrates how to use `azapi` provider to manage the Cosmos DB Mongo User Definition resource in Azure.



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

variable "mongodb_user_password" {
  type        = string
  description = "The password for the MongoDB user"
  sensitive   = true
}

resource "azapi_resource" "resourceGroup" {
  type     = "Microsoft.Resources/resourceGroups@2020-06-01"
  name     = var.resource_name
  location = var.location
}

resource "azapi_resource" "databaseAccount" {
  type      = "Microsoft.DocumentDB/databaseAccounts@2024-08-15"
  parent_id = azapi_resource.resourceGroup.id
  name      = "${var.resource_name}-acct"
  location  = var.location
  body = {
    kind = "MongoDB"
    properties = {
      backupPolicy = null
      capabilities = [{
        name = "EnableMongoRoleBasedAccessControl"
        }, {
        name = "EnableMongo"
      }]
      consistencyPolicy = {
        defaultConsistencyLevel = "Strong"
        maxIntervalInSeconds    = 5
        maxStalenessPrefix      = 100
      }
      databaseAccountOfferType           = "Standard"
      disableKeyBasedMetadataWriteAccess = false
      disableLocalAuth                   = false
      enableAnalyticalStorage            = false
      enableAutomaticFailover            = false
      enableBurstCapacity                = false
      enableFreeTier                     = false
      enableMultipleWriteLocations       = false
      enablePartitionMerge               = false
      ipRules                            = []
      isVirtualNetworkFilterEnabled      = false
      locations = [{
        failoverPriority = 0
        isZoneRedundant  = false
        locationName     = var.location
      }]
      minimalTlsVersion           = "Tls12"
      networkAclBypass            = "None"
      networkAclBypassResourceIds = []
      publicNetworkAccess         = "Enabled"
      virtualNetworkRules         = []
    }
  }
}

resource "azapi_resource" "mongodbDatabas" {
  type      = "Microsoft.DocumentDB/databaseAccounts/mongodbDatabases@2021-10-15"
  parent_id = azapi_resource.databaseAccount.id
  name      = "${var.resource_name}-mongodb"
  body = {
    properties = {
      options = {}
      resource = {
        id = "${var.resource_name}-mongodb"
      }
    }
  }
}

resource "azapi_resource" "mongodbUserDefinition" {
  type      = "Microsoft.DocumentDB/databaseAccounts/mongodbUserDefinitions@2022-11-15"
  parent_id = azapi_resource.databaseAccount.id
  name      = "${azapi_resource.mongodbDatabas.name}.myUserName"
  body = {
    properties = {
      databaseName = azapi_resource.mongodbDatabas.name
      mechanisms   = "SCRAM-SHA-256"
      password     = var.mongodb_user_password
      userName     = "myUserName"
    }
  }
}


```



## Arguments Reference

The following arguments are supported:

* `type` - (Required) The type of the resource. This should be set to `Microsoft.DocumentDB/databaseAccounts/mongodbUserDefinitions@api-version`. The available api-versions for this resource are: [`2021-10-15-preview`, `2021-11-15-preview`, `2022-02-15-preview`, `2022-05-15-preview`, `2022-08-15`, `2022-08-15-preview`, `2022-11-15`, `2022-11-15-preview`, `2023-03-01-preview`, `2023-03-15`, `2023-03-15-preview`, `2023-04-15`, `2023-09-15`, `2023-09-15-preview`, `2023-11-15`, `2023-11-15-preview`, `2024-02-15-preview`, `2024-05-15`, `2024-05-15-preview`, `2024-08-15`, `2024-09-01-preview`, `2024-11-15`, `2024-12-01-preview`, `2025-04-15`, `2025-05-01-preview`, `2025-10-15`, `2025-11-01-preview`].

* `parent_id` - (Required) The ID of the azure resource in which this resource is created. The allowed values are:  
  `/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.DocumentDB/databaseAccounts/{resourceName}`

* `name` - (Required) Specifies the name of the azure resource. Changing this forces a new resource to be created.

* `body` - (Required) Specifies the configuration of the resource. More information about the arguments in `body` can be found in the [Microsoft documentation](https://learn.microsoft.com/en-us/azure/templates/Microsoft.DocumentDB/databaseAccounts/mongodbUserDefinitions?pivots=deployment-language-terraform).

For other arguments, please refer to the [azapi_resource](https://registry.terraform.io/providers/Azure/azapi/latest/docs/resources/resource) documentation.

## Import

 ```shell
 # Azure resource can be imported using the resource id, e.g.
 terraform import azapi_resource.example /subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.DocumentDB/databaseAccounts/{resourceName}/mongodbUserDefinitions/{resourceName}
 
 # It also supports specifying API version by using the resource id with api-version as a query parameter, e.g.
 terraform import azapi_resource.example /subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.DocumentDB/databaseAccounts/{resourceName}/mongodbUserDefinitions/{resourceName}?api-version=2025-11-01-preview
 ```
