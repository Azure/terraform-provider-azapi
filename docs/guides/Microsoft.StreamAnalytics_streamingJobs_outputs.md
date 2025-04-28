---
subcategory: "Microsoft.StreamAnalytics - Azure Stream Analytics"
page_title: "streamingJobs/outputs"
description: |-
  Manages a Stream Analytics Output Table.
---

# Microsoft.StreamAnalytics/streamingJobs/outputs - Stream Analytics Output Table

This article demonstrates how to use `azapi` provider to manage the Stream Analytics Output Table resource in Azure.

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

resource "azapi_resource" "streamingJob" {
  type      = "Microsoft.StreamAnalytics/streamingJobs@2020-03-01"
  parent_id = azapi_resource.resourceGroup.id
  name      = var.resource_name
  location  = var.location
  body = {
    properties = {
      cluster = {
      }
      compatibilityLevel                 = "1.0"
      contentStoragePolicy               = "SystemAccount"
      dataLocale                         = "en-GB"
      eventsLateArrivalMaxDelayInSeconds = 60
      eventsOutOfOrderMaxDelayInSeconds  = 50
      eventsOutOfOrderPolicy             = "Adjust"
      jobType                            = "Cloud"
      outputErrorPolicy                  = "Drop"
      sku = {
        name = "Standard"
      }
      transformation = {
        name = "main"
        properties = {
          query          = "    SELECT *\n    INTO [YourOutputAlias]\n    FROM [YourInputAlias]\n"
          streamingUnits = 3
        }
      }
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

data "azapi_resource_action" "listKeys" {
  type                   = "Microsoft.Storage/storageAccounts@2021-09-01"
  resource_id            = azapi_resource.storageAccount.id
  action                 = "listKeys"
  response_export_values = ["*"]
}

resource "azapi_resource" "output" {
  type      = "Microsoft.StreamAnalytics/streamingJobs/outputs@2021-10-01-preview"
  parent_id = azapi_resource.streamingJob.id
  name      = var.resource_name
  body = {
    properties = {
      datasource = {
        properties = {
          accountKey   = data.azapi_resource_action.listKeys.output.keys[0].value
          accountName  = azapi_resource.storageAccount.name
          batchSize    = 100
          partitionKey = "foo"
          rowKey       = "bar"
          table        = "foobar"
        }
        type = "Microsoft.Storage/Table"
      }
      serialization = null
    }
  }
  schema_validation_enabled = false
  response_export_values    = ["*"]
}


```



## Arguments Reference

The following arguments are supported:

* `type` - (Required) The type of the resource. This should be set to `Microsoft.StreamAnalytics/streamingJobs/outputs@api-version`. The available api-versions for this resource are: [`2016-03-01`, `2017-04-01-preview`, `2020-03-01`, `2021-10-01-preview`].

* `parent_id` - (Required) The ID of the azure resource in which this resource is created. The allowed values are:  
  `/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.StreamAnalytics/streamingJobs/{resourceName}`

* `name` - (Required) Specifies the name of the azure resource. Changing this forces a new resource to be created.

* `body` - (Required) Specifies the configuration of the resource. More information about the arguments in `body` can be found in the [Microsoft documentation](https://learn.microsoft.com/en-us/azure/templates/Microsoft.StreamAnalytics/streamingJobs/outputs?pivots=deployment-language-terraform).

## Import

 ```shell
 # Azure resource can be imported using the resource id, e.g.
 terraform import azapi_resource.example /subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.StreamAnalytics/streamingJobs/{resourceName}/outputs/{resourceName}
 
 # It also supports specifying API version by using the resource id with api-version as a query parameter, e.g.
 terraform import azapi_resource.example /subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.StreamAnalytics/streamingJobs/{resourceName}/outputs/{resourceName}?api-version=2021-10-01-preview
 ```
