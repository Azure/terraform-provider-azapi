---
subcategory: "Microsoft.DocumentDB - Azure Cosmos DB"
page_title: "mongoClusters/users"
description: |-
  Manages a Azure Cosmos DB for MongoDB (vCore) User.
---

# Microsoft.DocumentDB/mongoClusters/users - Azure Cosmos DB for MongoDB (vCore) User

This article demonstrates how to use `azapi` provider to manage the Azure Cosmos DB for MongoDB (vCore) User resource in Azure.



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
  default = "westus3"
}

resource "azapi_resource" "resourceGroup" {
  type     = "Microsoft.Resources/resourceGroups@2020-06-01"
  name     = var.resource_name
  location = var.location
}

resource "azapi_resource" "userAssignedIdentity" {
  type                      = "Microsoft.ManagedIdentity/userAssignedIdentities@2023-01-31"
  parent_id                 = azapi_resource.resourceGroup.id
  name                      = var.resource_name
  location                  = var.location
  schema_validation_enabled = false
  response_export_values    = ["*"]
}

data "azapi_client_config" "current" {}

resource "azapi_resource" "mongoClusterSSDv2" {
  type      = "Microsoft.DocumentDB/mongoClusters@2025-09-01"
  parent_id = azapi_resource.resourceGroup.id
  name      = "${var.resource_name}-ssdv2"
  location  = var.location
  body = {
    properties = {
      authConfig = {
        allowedModes = ["MicrosoftEntraID"]
      }
      compute = {
        tier = "M30"
      }
      highAvailability = {
        targetMode = "Disabled"
      }
      serverVersion       = "6.0"
      publicNetworkAccess = "Disabled"
      sharding = {
        shardCount = 1
      }
      storage = {
        sizeGb = 64
        type   = "PremiumSSDv2"
      }
    }
  }
  schema_validation_enabled = false
  response_export_values    = ["*"]
}

resource "azapi_resource" "mongoUser_EntraServicePrincipal" {
  type      = "Microsoft.DocumentDB/mongoClusters/users@2025-09-01"
  name      = azapi_resource.userAssignedIdentity.output.properties.principalId
  parent_id = azapi_resource.mongoClusterSSDv2.id

  body = {
    properties = {
      roles = [
        {
          role = "root"
          db   = "admin"
        }
      ]
      identityProvider = {
        type = "MicrosoftEntraID"
        properties = {
          principalType = "ServicePrincipal"
        }
      }
    }
  }
  schema_validation_enabled = false
  response_export_values    = ["*"]
}

```



## Arguments Reference

The following arguments are supported:

* `type` - (Required) The type of the resource. This should be set to `Microsoft.DocumentDB/mongoClusters/users@api-version`. The available api-versions for this resource are: [`2025-04-01-preview`, `2025-07-01-preview`].

* `parent_id` - (Required) The ID of the azure resource in which this resource is created. The allowed values are:  
  `/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.DocumentDB/mongoClusters/{resourceName}`

* `name` - (Required) Specifies the name of the azure resource. Changing this forces a new resource to be created.

* `body` - (Required) Specifies the configuration of the resource. More information about the arguments in `body` can be found in the [Microsoft documentation](https://learn.microsoft.com/en-us/azure/templates/Microsoft.DocumentDB/mongoClusters/users?pivots=deployment-language-terraform).

For other arguments, please refer to the [azapi_resource](https://registry.terraform.io/providers/Azure/azapi/latest/docs/resources/resource) documentation.

## Import

 ```shell
 # Azure resource can be imported using the resource id, e.g.
 terraform import azapi_resource.example /subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.DocumentDB/mongoClusters/{resourceName}/users/{resourceName}
 
 # It also supports specifying API version by using the resource id with api-version as a query parameter, e.g.
 terraform import azapi_resource.example /subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.DocumentDB/mongoClusters/{resourceName}/users/{resourceName}?api-version=2025-07-01-preview
 ```
