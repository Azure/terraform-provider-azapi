---
subcategory: "Microsoft.CostManagement - Cost Management"
page_title: "exports"
description: |-
  Manages a Azure Cost Management Export.
---

# Microsoft.CostManagement/exports - Azure Cost Management Export

This article demonstrates how to use `azapi` provider to manage the Azure Cost Management Export resource in Azure.

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
  default = "westus"
}

resource "azapi_resource" "resourceGroup" {
  type     = "Microsoft.Resources/resourceGroups@2020-06-01"
  name     = var.resource_name
  location = var.location
}

resource "azapi_resource" "storageAccount" {
  type      = "Microsoft.Storage/storageAccounts@2023-05-01"
  parent_id = azapi_resource.resourceGroup.id
  name      = var.resource_name
  location  = var.location
  body = {
    kind = "StorageV2"
    properties = {
      accessTier                   = "Hot"
      allowBlobPublicAccess        = false
      allowCrossTenantReplication  = false
      allowSharedKeyAccess         = true
      defaultToOAuthAuthentication = false
      dnsEndpointType              = "Standard"
      encryption = {
        keySource = "Microsoft.Storage"
        services = {
          blob = {
            enabled = true
            keyType = "Account"
          }
          file = {
            enabled = true
            keyType = "Account"
          }
        }
      }
      isHnsEnabled       = false
      isLocalUserEnabled = true
      isNfsV3Enabled     = false
      isSftpEnabled      = false
      minimumTlsVersion  = "TLS1_2"
      networkAcls = {
        bypass              = "AzureServices"
        defaultAction       = "Allow"
        ipRules             = []
        resourceAccessRules = []
        virtualNetworkRules = []
      }
      publicNetworkAccess      = "Enabled"
      supportsHttpsTrafficOnly = true
    }
    sku = {
      name = "Standard_LRS"
    }
  }
}

resource "azapi_resource" "export" {
  type      = "Microsoft.CostManagement/exports@2023-08-01"
  parent_id = azapi_resource.resourceGroup.id
  name      = var.resource_name
  body = {
    properties = {
      definition = {
        timeframe = "TheLastMonth"
        type      = "Usage"
        dataSet = {
          granularity = "Daily"
        }
      }
      deliveryInfo = {
        destination = {
          container      = "exports"
          resourceId     = azapi_resource.storageAccount.id
          rootFolderPath = "ad-hoc"
        }
      }
      format = "Csv"
      schedule = {
        recurrence = "Monthly"
        recurrencePeriod = {
          from = timeadd(timestamp(), "24h")
          to   = timeadd(timestamp(), "744h")
        }
        status = "Active"
      }
    }
  }
}

```



## Arguments Reference

The following arguments are supported:

* `type` - (Required) The type of the resource. This should be set to `Microsoft.CostManagement/exports@api-version`. The available api-versions for this resource are: [`2019-01-01`, `2019-09-01`, `2019-10-01`, `2019-11-01`, `2020-06-01`, `2020-12-01-preview`, `2021-01-01`, `2021-10-01`, `2022-10-01`, `2023-03-01`, `2023-04-01-preview`, `2023-07-01-preview`, `2023-08-01`, `2023-09-01`, `2023-11-01`, `2024-08-01`, `2024-10-01-preview`, `2025-03-01`].

* `parent_id` - (Required) The ID of the azure resource in which this resource is created. The allowed values are:  
  `{any azure resource id}`

* `name` - (Required) Specifies the name of the azure resource. Changing this forces a new resource to be created.

* `body` - (Required) Specifies the configuration of the resource. More information about the arguments in `body` can be found in the [Microsoft documentation](https://learn.microsoft.com/en-us/azure/templates/Microsoft.CostManagement/exports?pivots=deployment-language-terraform).

For other arguments, please refer to the [azapi_resource](https://registry.terraform.io/providers/Azure/azapi/latest/docs/resources/resource) documentation.

## Import

 ```shell
 # Azure resource can be imported using the resource id, e.g.
 terraform import azapi_resource.example {any azure resource id}/providers/Microsoft.CostManagement/exports/{resourceName}
 
 # It also supports specifying API version by using the resource id with api-version as a query parameter, e.g.
 terraform import azapi_resource.example {any azure resource id}/providers/Microsoft.CostManagement/exports/{resourceName}?api-version=2025-03-01
 ```
