---
subcategory: "Microsoft.DataFactory - Data Factory"
page_title: "factories/datasets"
description: |-
  Manages a Dataset inside an Azure Data Factory.
---

# Microsoft.DataFactory/factories/datasets - Dataset inside an Azure Data Factory

This article demonstrates how to use `azapi` provider to manage the Dataset inside an Azure Data Factory resource in Azure.



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

resource "azapi_resource" "factory" {
  type      = "Microsoft.DataFactory/factories@2018-06-01"
  parent_id = azapi_resource.resourceGroup.id
  name      = var.resource_name
  location  = var.location
  body = {
    properties = {
      publicNetworkAccess = "Enabled"
      repoConfiguration   = null
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

resource "azapi_resource" "linkedservice" {
  type      = "Microsoft.DataFactory/factories/linkedservices@2018-06-01"
  parent_id = azapi_resource.factory.id
  name      = var.resource_name
  body = {
    properties = {
      description = ""
      type        = "AzureBlobStorage"
      typeProperties = {
        serviceEndpoint = azapi_resource.storageAccount.output.properties.primaryEndpoints.blob
      }
    }
  }
  schema_validation_enabled = false
  response_export_values    = ["*"]
}

resource "azapi_resource" "dataset" {
  type      = "Microsoft.DataFactory/factories/datasets@2018-06-01"
  parent_id = azapi_resource.factory.id
  name      = var.resource_name
  body = {
    properties = {
      description = ""
      linkedServiceName = {
        referenceName = azapi_resource.linkedservice.name
        type          = "LinkedServiceReference"
      }
      type = "Json"
      typeProperties = {
        encodingName = "UTF-8"
        location = {
          container  = "container"
          fileName   = "bar.txt"
          folderPath = "foo/bar/"
          type       = "AzureBlobStorageLocation"
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

* `type` - (Required) The type of the resource. This should be set to `Microsoft.DataFactory/factories/datasets@api-version`. The available api-versions for this resource are: [`2017-09-01-preview`, `2018-06-01`].

* `parent_id` - (Required) The ID of the azure resource in which this resource is created. The allowed values are:  
  `/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.DataFactory/factories/{resourceName}`

* `name` - (Required) Specifies the name of the azure resource. Changing this forces a new resource to be created.

* `body` - (Required) Specifies the configuration of the resource. More information about the arguments in `body` can be found in the [Microsoft documentation](https://learn.microsoft.com/en-us/azure/templates/Microsoft.DataFactory/factories/datasets?pivots=deployment-language-terraform).

For other arguments, please refer to the [azapi_resource](https://registry.terraform.io/providers/Azure/azapi/latest/docs/resources/resource) documentation.

## Import

 ```shell
 # Azure resource can be imported using the resource id, e.g.
 terraform import azapi_resource.example /subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.DataFactory/factories/{resourceName}/datasets/{resourceName}
 
 # It also supports specifying API version by using the resource id with api-version as a query parameter, e.g.
 terraform import azapi_resource.example /subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.DataFactory/factories/{resourceName}/datasets/{resourceName}?api-version=2018-06-01
 ```
