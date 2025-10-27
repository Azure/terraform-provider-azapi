---
subcategory: "Microsoft.DBforPostgreSQL - Azure Database for PostgreSQL"
page_title: "serverGroupsv2"
description: |-
  Manages a Azure Cosmos DB for PostgreSQL Cluster.
---

# Microsoft.DBforPostgreSQL/serverGroupsv2 - Azure Cosmos DB for PostgreSQL Cluster

This article demonstrates how to use `azapi` provider to manage the Azure Cosmos DB for PostgreSQL Cluster resource in Azure.



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

variable "administrator_login_password" {
  type        = string
  description = "The administrator login password for the PostgreSQL server group"
  sensitive   = true
}

resource "azapi_resource" "resourceGroup" {
  type     = "Microsoft.Resources/resourceGroups@2020-06-01"
  name     = var.resource_name
  location = var.location
}

resource "azapi_resource" "serverGroupsv2" {
  type      = "Microsoft.DBforPostgreSQL/serverGroupsv2@2022-11-08"
  parent_id = azapi_resource.resourceGroup.id
  name      = var.resource_name
  location  = var.location
  body = {
    properties = {
      administratorLoginPassword      = var.administrator_login_password
      coordinatorEnablePublicIpAccess = true
      coordinatorServerEdition        = "GeneralPurpose"
      coordinatorStorageQuotaInMb     = 131072
      coordinatorVCores               = 2
      enableHa                        = false
      nodeCount                       = 0
      nodeEnablePublicIpAccess        = false
      nodeServerEdition               = "MemoryOptimized"
    }
  }
  schema_validation_enabled = false
  response_export_values    = ["*"]
}


```



## Arguments Reference

The following arguments are supported:

* `type` - (Required) The type of the resource. This should be set to `Microsoft.DBforPostgreSQL/serverGroupsv2@api-version`. The available api-versions for this resource are: [`2022-11-08`, `2023-03-02-preview`].

* `parent_id` - (Required) The ID of the azure resource in which this resource is created. The allowed values are:  
  `/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}`

* `name` - (Required) Specifies the name of the azure resource. Changing this forces a new resource to be created.

* `body` - (Required) Specifies the configuration of the resource. More information about the arguments in `body` can be found in the [Microsoft documentation](https://learn.microsoft.com/en-us/azure/templates/Microsoft.DBforPostgreSQL/serverGroupsv2?pivots=deployment-language-terraform).

For other arguments, please refer to the [azapi_resource](https://registry.terraform.io/providers/Azure/azapi/latest/docs/resources/resource) documentation.

## Import

 ```shell
 # Azure resource can be imported using the resource id, e.g.
 terraform import azapi_resource.example /subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.DBforPostgreSQL/serverGroupsv2/{resourceName}
 
 # It also supports specifying API version by using the resource id with api-version as a query parameter, e.g.
 terraform import azapi_resource.example /subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.DBforPostgreSQL/serverGroupsv2/{resourceName}?api-version=2023-03-02-preview
 ```
