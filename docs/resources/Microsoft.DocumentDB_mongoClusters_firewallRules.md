---
subcategory: "Microsoft.DocumentDB - Azure Cosmos DB"
page_title: "mongoClusters/firewallRules"
description: |-
  Manages a Azure Cosmos DB for MongoDB (vCore) Firewall Rule.
---

# Microsoft.DocumentDB/mongoClusters/firewallRules - Azure Cosmos DB for MongoDB (vCore) Firewall Rule

This article demonstrates how to use `azapi` provider to manage the Azure Cosmos DB for MongoDB (vCore) Firewall Rule resource in Azure.



## Example Usage

### default

```hcl
terraform {
  required_providers {
    azapi = {
      source = "Azure/azapi"
    }
    time = {
      source = "hashicorp/time"
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

data "azapi_client_config" "current" {}

resource "azapi_resource" "mongoCluster" {
  type      = "Microsoft.DocumentDB/mongoClusters@2025-09-01"
  parent_id = azapi_resource.resourceGroup.id
  name      = var.resource_name
  location  = var.location
  body = {
    properties = {
      authConfig = {
        allowedModes = ["MicrosoftEntraID"]
      }
      compute = {
        tier = "M40"
      }
      highAvailability = {
        targetMode = "Disabled"
      }
      previewFeatures = [
        "ShardRebalancer"
      ]
      publicNetworkAccess = "Enabled"
      serverVersion       = "5.0"
      sharding = {
        shardCount = 1
      }
      storage = {
        sizeGb = 32
      }
    }
  }
  tags = {
    Environment = "Test"
  }
  schema_validation_enabled = false
  response_export_values    = ["*"]
}

resource "azapi_resource" "firewallRule" {
  type      = "Microsoft.DocumentDB/mongoClusters/firewallRules@2025-09-01"
  parent_id = azapi_resource.mongoCluster.id
  name      = var.resource_name
  body = {
    properties = {
      endIpAddress   = "0.0.0.0"
      startIpAddress = "0.0.0.0"
    }
  }
  schema_validation_enabled = false
  response_export_values    = ["*"]
}

```



## Arguments Reference

The following arguments are supported:

* `type` - (Required) The type of the resource. This should be set to `Microsoft.DocumentDB/mongoClusters/firewallRules@api-version`. The available api-versions for this resource are: [`2023-03-01-preview`, `2023-03-15-preview`, `2023-09-15-preview`, `2023-11-15-preview`, `2024-02-15-preview`, `2024-03-01-preview`, `2024-06-01-preview`, `2024-07-01`, `2024-10-01-preview`, `2025-04-01-preview`, `2025-07-01-preview`].

* `parent_id` - (Required) The ID of the azure resource in which this resource is created. The allowed values are:  
  `/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.DocumentDB/mongoClusters/{resourceName}`

* `name` - (Required) Specifies the name of the azure resource. Changing this forces a new resource to be created.

* `body` - (Required) Specifies the configuration of the resource. More information about the arguments in `body` can be found in the [Microsoft documentation](https://learn.microsoft.com/en-us/azure/templates/Microsoft.DocumentDB/mongoClusters/firewallRules?pivots=deployment-language-terraform).

For other arguments, please refer to the [azapi_resource](https://registry.terraform.io/providers/Azure/azapi/latest/docs/resources/resource) documentation.

## Import

 ```shell
 # Azure resource can be imported using the resource id, e.g.
 terraform import azapi_resource.example /subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.DocumentDB/mongoClusters/{resourceName}/firewallRules/{resourceName}
 
 # It also supports specifying API version by using the resource id with api-version as a query parameter, e.g.
 terraform import azapi_resource.example /subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.DocumentDB/mongoClusters/{resourceName}/firewallRules/{resourceName}?api-version=2025-07-01-preview
 ```
