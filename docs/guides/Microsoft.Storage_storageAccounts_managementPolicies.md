---
subcategory: "Microsoft.Storage - Storage"
page_title: "storageAccounts/managementPolicies"
description: |-
  Manages a Azure Storage Account Management Policy.
---

# Microsoft.Storage/storageAccounts/managementPolicies - Azure Storage Account Management Policy

This article demonstrates how to use `azapi` provider to manage the Azure Storage Account Management Policy resource in Azure.

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
    kind = "BlobStorage"
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

resource "azapi_resource" "managementPolicy" {
  type      = "Microsoft.Storage/storageAccounts/managementPolicies@2021-09-01"
  parent_id = azapi_resource.storageAccount.id
  name      = "default"
  body = {
    properties = {
      policy = {
        rules = [
          {
            definition = {
              actions = {
                baseBlob = {
                  delete = {
                    daysAfterModificationGreaterThan = 100
                  }
                  tierToArchive = {
                    daysAfterModificationGreaterThan = 50
                  }
                  tierToCool = {
                    daysAfterModificationGreaterThan = 10
                  }
                }
                snapshot = {
                  delete = {
                    daysAfterCreationGreaterThan = 30
                  }
                }
              }
              filters = {
                blobTypes = [
                  "blockBlob",
                ]
                prefixMatch = [
                  "container1/prefix1",
                ]
              }
            }
            enabled = true
            name    = "rule-1"
            type    = "Lifecycle"
          },
        ]
      }
    }
  }
  schema_validation_enabled = false
  response_export_values    = ["*"]
}


```



## Arguments Reference

The following arguments are supported:

* `type` - (Required) The type of the resource. This should be set to `Microsoft.Storage/storageAccounts/managementPolicies@api-version`. The available api-versions for this resource are: [`2018-03-01-preview`, `2018-11-01`, `2019-04-01`, `2019-06-01`, `2020-08-01-preview`, `2021-01-01`, `2021-02-01`, `2021-04-01`, `2021-06-01`, `2021-08-01`, `2021-09-01`, `2022-05-01`, `2022-09-01`, `2023-01-01`, `2023-04-01`, `2023-05-01`, `2024-01-01`].

* `parent_id` - (Required) The ID of the azure resource in which this resource is created. The allowed values are:  
  `/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Storage/storageAccounts/{resourceName}`

* `name` - (Required) Specifies the name of the azure resource. Changing this forces a new resource to be created.

* `body` - (Required) Specifies the configuration of the resource. More information about the arguments in `body` can be found in the [Microsoft documentation](https://learn.microsoft.com/en-us/azure/templates/Microsoft.Storage/storageAccounts/managementPolicies?pivots=deployment-language-terraform).

## Import

 ```shell
 # Azure resource can be imported using the resource id, e.g.
 terraform import azapi_resource.example /subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Storage/storageAccounts/{resourceName}/managementPolicies/{resourceName}
 
 # It also supports specifying API version by using the resource id with api-version as a query parameter, e.g.
 terraform import azapi_resource.example /subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Storage/storageAccounts/{resourceName}/managementPolicies/{resourceName}?api-version=2024-01-01
 ```
