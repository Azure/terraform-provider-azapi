---
subcategory: "Microsoft.Media - Media Services"
page_title: "mediaServices/assets/assetFilters"
description: |-
  Manages a Media Services Asset Filters.
---

# Microsoft.Media/mediaServices/assets/assetFilters - Media Services Asset Filters

This article demonstrates how to use `azapi` provider to manage the Media Services Asset Filters resource in Azure.

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
      name = "Standard_GRS"
    }
  }
  schema_validation_enabled = false
  response_export_values    = ["*"]
}

resource "azapi_resource" "mediaService" {
  type      = "Microsoft.Media/mediaServices@2021-11-01"
  parent_id = azapi_resource.resourceGroup.id
  name      = var.resource_name
  location  = var.location
  body = {
    properties = {
      publicNetworkAccess = "Enabled"
      storageAccounts = [
        {
          id   = azapi_resource.storageAccount.id
          type = "Primary"
        },
      ]
    }
  }
  schema_validation_enabled = false
  response_export_values    = ["*"]
}

resource "azapi_resource" "asset" {
  type      = "Microsoft.Media/mediaServices/assets@2022-08-01"
  parent_id = azapi_resource.mediaService.id
  name      = var.resource_name
  body = {
    properties = {
      description = ""
    }
  }
  schema_validation_enabled = false
  response_export_values    = ["*"]
}

resource "azapi_resource" "assetFilter" {
  type      = "Microsoft.Media/mediaServices/assets/assetFilters@2022-08-01"
  parent_id = azapi_resource.asset.id
  name      = var.resource_name
  body = {
    properties = {
      firstQuality = {
        bitrate = 0
      }
    }
  }
  schema_validation_enabled = false
  response_export_values    = ["*"]
}


```



## Arguments Reference

The following arguments are supported:

* `type` - (Required) The type of the resource. This should be set to `Microsoft.Media/mediaServices/assets/assetFilters@api-version`. The available api-versions for this resource are: [].

* `parent_id` - (Required) The ID of the azure resource in which this resource is created. The allowed values are:  
  

* `name` - (Required) Specifies the name of the azure resource. Changing this forces a new resource to be created.

* `body` - (Required) Specifies the configuration of the resource. More information about the arguments in `body` can be found in the [Microsoft documentation](https://learn.microsoft.com/en-us/azure/templates/Microsoft.Media/mediaServices/assets/assetFilters?pivots=deployment-language-terraform).

For other arguments, please refer to the [azapi_resource](https://registry.terraform.io/providers/Azure/azapi/latest/docs/resources/resource) documentation.

## Import

 ```shell
 # Azure resource can be imported using the resource id, e.g.
 terraform import azapi_resource.example 
 
 # It also supports specifying API version by using the resource id with api-version as a query parameter, e.g.
 terraform import azapi_resource.example ?api-version=API_VERSION
 ```
