---
subcategory: "Microsoft.HDInsight - HDInsight"
page_title: "clusters"
description: |-
  Manages a HDInsight Cluster.
---

# Microsoft.HDInsight/clusters - HDInsight Cluster

This article demonstrates how to use `azapi` provider to manage the HDInsight Cluster resource in Azure.

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

variable "vm_username" {
  type        = string
  description = "The username for the HDInsight cluster virtual machines"
}

variable "vm_password" {
  type        = string
  description = "The password for the HDInsight cluster virtual machines"
  sensitive   = true
}

variable "rest_credential_password" {
  type        = string
  description = "The REST API credential password for the HDInsight cluster gateway"
  sensitive   = true
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

data "azapi_resource" "blobService" {
  type      = "Microsoft.Storage/storageAccounts/blobServices@2022-09-01"
  parent_id = azapi_resource.storageAccount.id
  name      = "default"
}

resource "azapi_resource" "container" {
  type      = "Microsoft.Storage/storageAccounts/blobServices/containers@2022-09-01"
  name      = var.resource_name
  parent_id = data.azapi_resource.blobService.id
  body = {
    properties = {
      metadata = {
        key = "value"
      }
    }
  }
  response_export_values = ["*"]
}

resource "azapi_resource" "cluster" {
  type      = "Microsoft.HDInsight/clusters@2018-06-01-preview"
  parent_id = azapi_resource.resourceGroup.id
  name      = var.resource_name
  location  = var.location
  body = {
    properties = {
      clusterDefinition = {
        componentVersion = {
          Spark = "2.4"
        }
        configurations = {
          gateway = {
            "restAuthCredential.isEnabled" = true
            "restAuthCredential.password"  = var.rest_credential_password
            "restAuthCredential.username"  = "acctestusrgw"
          }
        }
        kind = "Spark"
      }
      clusterVersion = "4.0.3000.1"
      computeProfile = {
        roles = [
          {
            hardwareProfile = {
              vmSize = "standard_a4_v2"
            }
            name = "headnode"
            osProfile = {
              linuxOperatingSystemProfile = {
                password = var.vm_password
                username = var.vm_username
              }
            }
            targetInstanceCount = 2
          },
          {
            hardwareProfile = {
              vmSize = "standard_a4_v2"
            }
            name = "workernode"
            osProfile = {
              linuxOperatingSystemProfile = {
                password = var.vm_password
                username = var.vm_username
              }
            }
            targetInstanceCount = 3
          },
          {
            hardwareProfile = {
              vmSize = "standard_a2_v2"
            }
            name = "zookeepernode"
            osProfile = {
              linuxOperatingSystemProfile = {
                password = var.vm_password
                username = var.vm_username
              }
            }
            targetInstanceCount = 3
          },
        ]
      }
      encryptionInTransitProperties = {
        isEncryptionInTransitEnabled = false
      }
      minSupportedTlsVersion = "1.2"
      osType                 = "Linux"
      storageProfile = {
        storageaccounts = [
          {
            container  = azapi_resource.container.name
            isDefault  = true
            key        = data.azapi_resource_action.listKeys.output.keys[0].value
            name       = "${azapi_resource.storageAccount.name}.blob.core.windows.net"
            resourceId = azapi_resource.storageAccount.id
          },
        ]
      }
      tier = "standard"
    }
  }
  schema_validation_enabled = false
  response_export_values    = ["*"]
  timeouts {
    create = "180m"
    update = "180m"
    delete = "60m"
  }
}


```



## Arguments Reference

The following arguments are supported:

* `type` - (Required) The type of the resource. This should be set to `Microsoft.HDInsight/clusters@api-version`. The available api-versions for this resource are: [`2015-03-01-preview`, `2018-06-01-preview`, `2021-06-01`, `2023-04-15-preview`, `2023-08-15-preview`, `2024-08-01-preview`, `2025-01-15-preview`].

* `parent_id` - (Required) The ID of the azure resource in which this resource is created. The allowed values are:  
  `/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}`

* `name` - (Required) Specifies the name of the azure resource. Changing this forces a new resource to be created.

* `body` - (Required) Specifies the configuration of the resource. More information about the arguments in `body` can be found in the [Microsoft documentation](https://learn.microsoft.com/en-us/azure/templates/Microsoft.HDInsight/clusters?pivots=deployment-language-terraform).

For other arguments, please refer to the [azapi_resource](https://registry.terraform.io/providers/Azure/azapi/latest/docs/resources/resource) documentation.

## Import

 ```shell
 # Azure resource can be imported using the resource id, e.g.
 terraform import azapi_resource.example /subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.HDInsight/clusters/{resourceName}
 
 # It also supports specifying API version by using the resource id with api-version as a query parameter, e.g.
 terraform import azapi_resource.example /subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.HDInsight/clusters/{resourceName}?api-version=2025-01-15-preview
 ```
