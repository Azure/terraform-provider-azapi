---
subcategory: "Microsoft.ServiceLinker - Service Linker"
page_title: "linkers"
description: |-
  Manages a service connector.
---

# Microsoft.ServiceLinker/linkers - service connector

This article demonstrates how to use `azapi` provider to manage the service connector resource in Azure.

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

resource "azapi_resource" "Spring" {
  type      = "Microsoft.AppPlatform/Spring@2023-05-01-preview"
  parent_id = azapi_resource.resourceGroup.id
  name      = var.resource_name
  location  = var.location
  body = {
    properties = {
      zoneRedundant = false
    }
    sku = {
      name = "S0"
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
        defaultConsistencyLevel = "BoundedStaleness"
        maxIntervalInSeconds    = 10
        maxStalenessPrefix      = 200
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
        throughput = 400
      }
      resource = {
        id = var.resource_name
      }
    }
  }
  schema_validation_enabled = false
  response_export_values    = ["*"]
}

resource "azapi_resource" "app" {
  type      = "Microsoft.AppPlatform/Spring/apps@2023-05-01-preview"
  parent_id = azapi_resource.Spring.id
  name      = var.resource_name
  location  = var.location
  identity {
    type         = "SystemAssigned"
    identity_ids = []
  }
  body = {
    properties = {
      customPersistentDisks = [
      ]
      enableEndToEndTLS = false
      public            = false
    }
  }
  schema_validation_enabled = false
  response_export_values    = ["*"]
}

resource "azapi_resource" "deployment" {
  type      = "Microsoft.AppPlatform/Spring/apps/deployments@2023-05-01-preview"
  parent_id = azapi_resource.app.id
  name      = "deploy-q4uff"
  body = {
    properties = {
      deploymentSettings = {
        environmentVariables = {
        }
        resourceRequests = {
          cpu    = "1"
          memory = "1Gi"
        }
      }
      source = {
        jvmOptions     = ""
        relativePath   = "<default>"
        runtimeVersion = "Java_8"
        type           = "Jar"
      }
    }
    sku = {
      capacity = 1
      name     = "S0"
      tier     = "Standard"
    }
  }
  schema_validation_enabled = false
  response_export_values    = ["*"]
}

resource "azapi_resource" "linker" {
  type      = "Microsoft.ServiceLinker/linkers@2022-05-01"
  parent_id = azapi_resource.deployment.id
  name      = var.resource_name
  body = {
    properties = {
      authInfo = {
        authType = "systemAssignedIdentity"
      }
      clientType = "none"
      targetService = {
        id                 = azapi_resource.sqlDatabase.id
        resourceProperties = null
        type               = "AzureResource"
      }
    }
  }
  schema_validation_enabled = false
  response_export_values    = ["*"]
}


```



## Arguments Reference

The following arguments are supported:

* `type` - (Required) The type of the resource. This should be set to `Microsoft.ServiceLinker/linkers@api-version`. The available api-versions for this resource are: [`2021-11-01-preview`, `2022-01-01-preview`, `2022-05-01`, `2022-11-01-preview`, `2023-04-01-preview`, `2024-04-01`, `2024-07-01-preview`].

* `parent_id` - (Required) The ID of the azure resource in which this resource is created. The allowed values are:  
  `{any azure resource id}`

* `name` - (Required) Specifies the name of the azure resource. Changing this forces a new resource to be created.

* `body` - (Required) Specifies the configuration of the resource. More information about the arguments in `body` can be found in the [Microsoft documentation](https://learn.microsoft.com/en-us/azure/templates/Microsoft.ServiceLinker/linkers?pivots=deployment-language-terraform).

## Import

 ```shell
 # Azure resource can be imported using the resource id, e.g.
 terraform import azapi_resource.example {any azure resource id}/providers/Microsoft.ServiceLinker/linkers/{resourceName}
 
 # It also supports specifying API version by using the resource id with api-version as a query parameter, e.g.
 terraform import azapi_resource.example {any azure resource id}/providers/Microsoft.ServiceLinker/linkers/{resourceName}?api-version=2024-07-01-preview
 ```
