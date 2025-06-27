---
subcategory: "Microsoft.EventGrid - Event Grid"
page_title: "eventSubscriptions"
description: |-
  Manages a EventGrid Event Subscription.
---

# Microsoft.EventGrid/eventSubscriptions - EventGrid Event Subscription

This article demonstrates how to use `azapi` provider to manage the EventGrid Event Subscription resource in Azure.

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

resource "azapi_resource" "namespace" {
  type      = "Microsoft.EventHub/namespaces@2022-01-01-preview"
  parent_id = azapi_resource.resourceGroup.id
  name      = var.resource_name
  location  = var.location
  body = {
    properties = {
      disableLocalAuth     = false
      isAutoInflateEnabled = false
      publicNetworkAccess  = "Enabled"
      zoneRedundant        = false
    }
    sku = {
      capacity = 1
      name     = "Standard"
      tier     = "Standard"
    }
  }
  schema_validation_enabled = false
  response_export_values    = ["*"]
}

resource "azapi_resource" "storageAccount" {
  type      = "Microsoft.Storage/storageAccounts@2021-09-01"
  parent_id = azapi_resource.resourceGroup.id
  name      = var.resource_name
  location  = var.location
  body = {
    kind = "StorageV2"
    properties = {
      accessTier                   = "Hot"
      allowBlobPublicAccess        = true
      allowCrossTenantReplication  = true
      allowSharedKeyAccess         = true
      defaultToOAuthAuthentication = false
      encryption = {
        keySource = "Microsoft.Storage"
        services = {
          queue = {
            keyType = "Service"
          }
          table = {
            keyType = "Service"
          }
        }
      }
      isHnsEnabled      = false
      isNfsV3Enabled    = false
      isSftpEnabled     = false
      minimumTlsVersion = "TLS1_2"
      networkAcls = {
        defaultAction = "Allow"
      }
      publicNetworkAccess      = "Enabled"
      supportsHttpsTrafficOnly = true
    }
    sku = {
      name = "Standard_LRS"
    }
  }
  schema_validation_enabled = false
  response_export_values    = ["*"]
}

resource "azapi_resource" "eventhub" {
  type      = "Microsoft.EventHub/namespaces/eventhubs@2021-11-01"
  parent_id = azapi_resource.namespace.id
  name      = var.resource_name
  body = {
    properties = {
      messageRetentionInDays = 1
      partitionCount         = 1
      status                 = "Active"
    }
  }
  schema_validation_enabled = false
  response_export_values    = ["*"]
}

resource "azapi_resource" "eventSubscription" {
  type      = "Microsoft.EventGrid/eventSubscriptions@2021-12-01"
  parent_id = azapi_resource.storageAccount.id
  name      = var.resource_name
  body = {
    properties = {
      deadLetterDestination = null
      destination = {
        endpointType = "EventHub"
        properties = {
          deliveryAttributeMappings = null
          resourceId                = azapi_resource.eventhub.id
        }
      }
      eventDeliverySchema = "EventGridSchema"
      filter = {
        includedEventTypes = [
          "Microsoft.Storage.BlobCreated",
          "Microsoft.Storage.BlobRenamed",
        ]
      }
      labels = [
      ]
      retryPolicy = {
        eventTimeToLiveInMinutes = 144
        maxDeliveryAttempts      = 10
      }
    }
  }
  schema_validation_enabled = false
  response_export_values    = ["*"]
}


```



## Arguments Reference

The following arguments are supported:

* `type` - (Required) The type of the resource. This should be set to `Microsoft.EventGrid/eventSubscriptions@api-version`. The available api-versions for this resource are: [`2017-06-15-preview`, `2017-09-15-preview`, `2018-01-01`, `2018-05-01-preview`, `2018-09-15-preview`, `2019-01-01`, `2019-02-01-preview`, `2019-06-01`, `2020-01-01-preview`, `2020-04-01-preview`, `2020-06-01`, `2020-10-15-preview`, `2021-06-01-preview`, `2021-10-15-preview`, `2021-12-01`, `2022-06-15`, `2023-06-01-preview`, `2023-12-15-preview`, `2024-06-01-preview`, `2024-12-15-preview`, `2025-02-15`, `2025-04-01-preview`].

* `parent_id` - (Required) The ID of the azure resource in which this resource is created. The allowed values are:  
  `{any azure resource id}`

* `name` - (Required) Specifies the name of the azure resource. Changing this forces a new resource to be created.

* `body` - (Required) Specifies the configuration of the resource. More information about the arguments in `body` can be found in the [Microsoft documentation](https://learn.microsoft.com/en-us/azure/templates/Microsoft.EventGrid/eventSubscriptions?pivots=deployment-language-terraform).

For other arguments, please refer to the [azapi_resource](https://registry.terraform.io/providers/Azure/azapi/latest/docs/resources/resource) documentation.

## Import

 ```shell
 # Azure resource can be imported using the resource id, e.g.
 terraform import azapi_resource.example {any azure resource id}/providers/Microsoft.EventGrid/eventSubscriptions/{resourceName}
 
 # It also supports specifying API version by using the resource id with api-version as a query parameter, e.g.
 terraform import azapi_resource.example {any azure resource id}/providers/Microsoft.EventGrid/eventSubscriptions/{resourceName}?api-version=2025-04-01-preview
 ```
