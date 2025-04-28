---
subcategory: "Microsoft.DocumentDB - Azure Cosmos DB"
page_title: "databaseAccounts/sqlRoleAssignments"
description: |-
  Manages a Cosmos DB SQL Role Assignment.
---

# Microsoft.DocumentDB/databaseAccounts/sqlRoleAssignments - Cosmos DB SQL Role Assignment

This article demonstrates how to use `azapi` provider to manage the Cosmos DB SQL Role Assignment resource in Azure.

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

resource "azapi_resource" "cluster" {
  type      = "Microsoft.Kusto/clusters@2023-05-02"
  parent_id = azapi_resource.resourceGroup.id
  name      = var.resource_name
  location  = var.location
  identity {
    type         = "SystemAssigned"
    identity_ids = []
  }
  body = {
    properties = {
      enableAutoStop                = true
      enableDiskEncryption          = false
      enableDoubleEncryption        = false
      enablePurge                   = false
      enableStreamingIngest         = false
      engineType                    = "V2"
      publicIPType                  = "IPv4"
      publicNetworkAccess           = "Enabled"
      restrictOutboundNetworkAccess = "Disabled"
      trustedExternalTenants = [
      ]
    }
    sku = {
      capacity = 1
      name     = "Dev(No SLA)_Standard_D11_v2"
      tier     = "Basic"
    }
  }
  schema_validation_enabled = false
  response_export_values    = ["*"]
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

data "azapi_resource" "sqlRoleDefinition" {
  type                   = "Microsoft.DocumentDB/databaseAccounts/sqlRoleDefinitions@2021-10-15"
  parent_id              = azapi_resource.databaseAccount.id
  name                   = "00000000-0000-0000-0000-000000000001"
  response_export_values = ["*"]
}

resource "azapi_resource" "database" {
  type      = "Microsoft.Kusto/clusters/databases@2023-05-02"
  parent_id = azapi_resource.cluster.id
  name      = var.resource_name
  location  = var.location
  body = {
    kind = "ReadWrite"
    properties = {
    }
  }
  schema_validation_enabled = false
  response_export_values    = ["*"]
}

resource "azapi_resource" "sqlRoleAssignment" {
  type      = "Microsoft.DocumentDB/databaseAccounts/sqlRoleAssignments@2021-10-15"
  parent_id = azapi_resource.databaseAccount.id
  name      = "ff419bf7-f8ca-ef51-00d2-3576700c341b"
  body = {
    properties = {
      principalId      = azapi_resource.cluster.output.identity.principalId
      roleDefinitionId = data.azapi_resource.sqlRoleDefinition.id
      scope            = azapi_resource.databaseAccount.id
    }
  }
  schema_validation_enabled = false
  response_export_values    = ["*"]
}


```



## Arguments Reference

The following arguments are supported:

* `type` - (Required) The type of the resource. This should be set to `Microsoft.DocumentDB/databaseAccounts/sqlRoleAssignments@api-version`. The available api-versions for this resource are: [`2020-06-01-preview`, `2021-03-01-preview`, `2021-04-01-preview`, `2021-04-15`, `2021-05-15`, `2021-06-15`, `2021-07-01-preview`, `2021-10-15`, `2021-10-15-preview`, `2021-11-15-preview`, `2022-02-15-preview`, `2022-05-15`, `2022-05-15-preview`, `2022-08-15`, `2022-08-15-preview`, `2022-11-15`, `2022-11-15-preview`, `2023-03-01-preview`, `2023-03-15`, `2023-03-15-preview`, `2023-04-15`, `2023-09-15`, `2023-09-15-preview`, `2023-11-15`, `2023-11-15-preview`, `2024-02-15-preview`, `2024-05-15`, `2024-05-15-preview`, `2024-08-15`, `2024-09-01-preview`, `2024-11-15`, `2024-12-01-preview`].

* `parent_id` - (Required) The ID of the azure resource in which this resource is created. The allowed values are:  
  `/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.DocumentDB/databaseAccounts/{resourceName}`

* `name` - (Required) Specifies the name of the azure resource. Changing this forces a new resource to be created.

* `body` - (Required) Specifies the configuration of the resource. More information about the arguments in `body` can be found in the [Microsoft documentation](https://learn.microsoft.com/en-us/azure/templates/Microsoft.DocumentDB/databaseAccounts/sqlRoleAssignments?pivots=deployment-language-terraform).

## Import

 ```shell
 # Azure resource can be imported using the resource id, e.g.
 terraform import azapi_resource.example /subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.DocumentDB/databaseAccounts/{resourceName}/sqlRoleAssignments/{resourceName}
 
 # It also supports specifying API version by using the resource id with api-version as a query parameter, e.g.
 terraform import azapi_resource.example /subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.DocumentDB/databaseAccounts/{resourceName}/sqlRoleAssignments/{resourceName}?api-version=2024-12-01-preview
 ```
