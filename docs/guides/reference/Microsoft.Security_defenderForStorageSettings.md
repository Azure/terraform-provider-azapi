---
subcategory: "Microsoft.Security - Security Center"
page_title: "defenderForStorageSettings"
description: |-
  Manages a Microsoft Defender for Storage.
---

# Microsoft.Security/defenderForStorageSettings - Microsoft Defender for Storage

This article demonstrates how to use `azapi` provider to manage the Microsoft Defender for Storage resource in Azure.

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
    }
    sku = {
      name = "Standard_LRS"
    }
  }
  schema_validation_enabled = false
  response_export_values    = ["*"]
}

resource "azapi_update_resource" "defenderForStorageSetting" {
  type      = "Microsoft.Security/defenderForStorageSettings@2022-12-01-preview"
  parent_id = azapi_resource.storageAccount.id
  name      = "current"
  body = {
    properties = {
      isEnabled = true
      malwareScanning = {
        onUpload = {
          capGBPerMonth = 5000
          isEnabled     = true
        }
      }
      sensitiveDataDiscovery = {
        isEnabled = true
      }
      overrideSubscriptionLevelSettings = true
    }
  }
}

```



## Arguments Reference

The following arguments are supported:

* `type` - (Required) The type of the resource. This should be set to `Microsoft.Security/defenderForStorageSettings@api-version`. The available api-versions for this resource are: [`2022-12-01-preview`, `2024-08-01-preview`, `2024-10-01-preview`, `2025-01-01`].

* `parent_id` - (Required) The ID of the azure resource in which this resource is created. The allowed values are:  
  `{any azure resource id}`

* `name` - (Required) Specifies the name of the azure resource. Changing this forces a new resource to be created.

* `body` - (Required) Specifies the configuration of the resource. More information about the arguments in `body` can be found in the [Microsoft documentation](https://learn.microsoft.com/en-us/azure/templates/Microsoft.Security/defenderForStorageSettings?pivots=deployment-language-terraform).

For other arguments, please refer to the [azapi_resource](https://registry.terraform.io/providers/Azure/azapi/latest/docs/resources/resource) documentation.

## Import

 ```shell
 # Azure resource can be imported using the resource id, e.g.
 terraform import azapi_resource.example {any azure resource id}/providers/Microsoft.Security/defenderForStorageSettings/{resourceName}
 
 # It also supports specifying API version by using the resource id with api-version as a query parameter, e.g.
 terraform import azapi_resource.example {any azure resource id}/providers/Microsoft.Security/defenderForStorageSettings/{resourceName}?api-version=2025-01-01
 ```
