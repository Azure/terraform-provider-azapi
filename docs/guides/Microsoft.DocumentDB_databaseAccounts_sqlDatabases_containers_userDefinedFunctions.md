---
subcategory: "Microsoft.DocumentDB - Azure Cosmos DB"
page_title: "databaseAccounts/sqlDatabases/containers/userDefinedFunctions"
description: |-
  Manages a SQL User Defined Function.
---

# Microsoft.DocumentDB/databaseAccounts/sqlDatabases/containers/userDefinedFunctions - SQL User Defined Function

This article demonstrates how to use `azapi` provider to manage the SQL User Defined Function resource in Azure.

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

resource "azapi_resource" "databaseAccount" {
  type      = "Microsoft.DocumentDB/databaseAccounts@2021-10-15"
  parent_id = azapi_resource.resourceGroup.id
  name      = var.resource_name
  location  = var.location
  body = {
    kind = "GlobalDocumentDB"
    properties = {
      capabilities = [
      ]
      consistencyPolicy = {
        defaultConsistencyLevel = "Session"
        maxIntervalInSeconds    = 5
        maxStalenessPrefix      = 100
      }
      databaseAccountOfferType           = "Standard"
      defaultIdentity                    = "FirstPartyIdentity"
      disableKeyBasedMetadataWriteAccess = false
      disableLocalAuth                   = false
      enableAnalyticalStorage            = false
      enableAutomaticFailover            = false
      enableFreeTier                     = false
      enableMultipleWriteLocations       = false
      ipRules = [
      ]
      isVirtualNetworkFilterEnabled = false
      locations = [
        {
          failoverPriority = 0
          isZoneRedundant  = false
          locationName     = "West Europe"
        },
      ]
      networkAclBypass = "None"
      networkAclBypassResourceIds = [
      ]
      publicNetworkAccess = "Enabled"
      virtualNetworkRules = [
      ]
    }
  }
  schema_validation_enabled = false
  response_export_values    = ["*"]
}

resource "azapi_resource" "sqlDatabase" {
  type      = "Microsoft.DocumentDB/databaseAccounts/sqlDatabases@2021-10-15"
  parent_id = azapi_resource.databaseAccount.id
  name      = var.resource_name
  body = {
    properties = {
      options = {
      }
      resource = {
        id = var.resource_name
      }
    }
  }
  schema_validation_enabled = false
  response_export_values    = ["*"]
}

resource "azapi_resource" "container" {
  type      = "Microsoft.DocumentDB/databaseAccounts/sqlDatabases/containers@2023-04-15"
  parent_id = azapi_resource.sqlDatabase.id
  name      = var.resource_name
  body = {
    properties = {
      options = {
      }
      resource = {
        id = var.resource_name
        partitionKey = {
          kind = "Hash"
          paths = [
            "/definition/id",
          ]
        }
      }
    }
  }
  schema_validation_enabled = false
  response_export_values    = ["*"]
}

resource "azapi_resource" "userDefinedFunction" {
  type      = "Microsoft.DocumentDB/databaseAccounts/sqlDatabases/containers/userDefinedFunctions@2021-10-15"
  parent_id = azapi_resource.container.id
  name      = var.resource_name
  body = {
    properties = {
      options = {
      }
      resource = {
        body = "  \tfunction test() {\n\t\tvar context = getContext();\n\t\tvar response = context.getResponse();\n\t\tresponse.setBody('Hello, World');\n\t}\n"
        id   = var.resource_name
      }
    }
  }
  schema_validation_enabled = false
  response_export_values    = ["*"]
}


```



## Arguments Reference

The following arguments are supported:

* `type` - (Required) The type of the resource. This should be set to `Microsoft.DocumentDB/databaseAccounts/sqlDatabases/containers/userDefinedFunctions@api-version`. The available api-versions for this resource are: [`2019-08-01`, `2019-12-12`, `2020-03-01`, `2020-04-01`, `2020-06-01-preview`, `2020-09-01`, `2021-01-15`, `2021-03-01-preview`, `2021-03-15`, `2021-04-01-preview`, `2021-04-15`, `2021-05-15`, `2021-06-15`, `2021-07-01-preview`, `2021-10-15`, `2021-10-15-preview`, `2021-11-15-preview`, `2022-02-15-preview`, `2022-05-15`, `2022-05-15-preview`, `2022-08-15`, `2022-08-15-preview`, `2022-11-15`, `2022-11-15-preview`, `2023-03-01-preview`, `2023-03-15`, `2023-03-15-preview`, `2023-04-15`, `2023-09-15`, `2023-09-15-preview`, `2023-11-15`, `2023-11-15-preview`, `2024-02-15-preview`, `2024-05-15`, `2024-05-15-preview`, `2024-08-15`, `2024-09-01-preview`, `2024-11-15`, `2024-12-01-preview`].

* `parent_id` - (Required) The ID of the azure resource in which this resource is created. The allowed values are:  
  `/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.DocumentDB/databaseAccounts/{resourceName}/sqlDatabases/{resourceName}/containers/{resourceName}`

* `name` - (Required) Specifies the name of the azure resource. Changing this forces a new resource to be created.

* `body` - (Required) Specifies the configuration of the resource. More information about the arguments in `body` can be found in the [Microsoft documentation](https://learn.microsoft.com/en-us/azure/templates/Microsoft.DocumentDB/databaseAccounts/sqlDatabases/containers/userDefinedFunctions?pivots=deployment-language-terraform).

For other arguments, please refer to the [azapi_resource](https://registry.terraform.io/providers/Azure/azapi/latest/docs/resources/resource) documentation.

## Import

 ```shell
 # Azure resource can be imported using the resource id, e.g.
 terraform import azapi_resource.example /subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.DocumentDB/databaseAccounts/{resourceName}/sqlDatabases/{resourceName}/containers/{resourceName}/userDefinedFunctions/{resourceName}
 
 # It also supports specifying API version by using the resource id with api-version as a query parameter, e.g.
 terraform import azapi_resource.example /subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.DocumentDB/databaseAccounts/{resourceName}/sqlDatabases/{resourceName}/containers/{resourceName}/userDefinedFunctions/{resourceName}?api-version=2024-12-01-preview
 ```
