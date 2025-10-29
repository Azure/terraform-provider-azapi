---
subcategory: "Microsoft.Compute - Virtual Machines, Virtual Machine Scale Sets"
page_title: "galleries/applications/versions"
description: |-
  Manages a Gallery Application Version.
---

# Microsoft.Compute/galleries/applications/versions - Gallery Application Version

This article demonstrates how to use `azapi` provider to manage the Gallery Application Version resource in Azure.



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

resource "azapi_resource" "gallery" {
  type      = "Microsoft.Compute/galleries@2022-03-03"
  parent_id = azapi_resource.resourceGroup.id
  name      = "${var.resource_name}sig"
  location  = var.location
  body = {
    properties = {
      description = ""
    }
  }
}

resource "azapi_resource" "storageAccount" {
  type      = "Microsoft.Storage/storageAccounts@2023-05-01"
  parent_id = azapi_resource.resourceGroup.id
  name      = "${var.resource_name}acc"
  location  = var.location
  body = {
    kind = "StorageV2"
    properties = {
      accessTier                   = "Hot"
      allowBlobPublicAccess        = true
      allowCrossTenantReplication  = false
      allowSharedKeyAccess         = true
      defaultToOAuthAuthentication = false
      dnsEndpointType              = "Standard"
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

resource "azapi_resource" "application" {
  type      = "Microsoft.Compute/galleries/applications@2022-03-03"
  parent_id = azapi_resource.gallery.id
  name      = "${var.resource_name}-app"
  location  = var.location
  body = {
    properties = {
      supportedOSType = "Linux"
    }
  }
}

resource "azapi_resource" "container" {
  type      = "Microsoft.Storage/storageAccounts/blobServices/containers@2023-05-01"
  parent_id = "${azapi_resource.storageAccount.id}/blobServices/default"
  name      = "mycontainer"
  body = {
    properties = {
      publicAccess = "Blob"
    }
  }
}

resource "azapi_resource" "version" {
  type      = "Microsoft.Compute/galleries/applications/versions@2022-03-03"
  parent_id = azapi_resource.application.id
  name      = "0.0.1"
  location  = var.location
  body = {
    properties = {
      publishingProfile = {
        enableHealthCheck = false
        excludeFromLatest = false
        manageActions = {
          install = "[install command]"
          remove  = "[remove command]"
          update  = ""
        }
        source = {
          defaultConfigurationLink = ""
          mediaLink                = "https://${azapi_resource.storageAccount.name}.blob.core.windows.net/mycontainer/myblob"
        }
        targetRegions = [{
          name                 = var.location
          regionalReplicaCount = 1
          storageAccountType   = "Standard_LRS"
        }]
      }
      safetyProfile = {
        allowDeletionOfReplicatedLocations = true
      }
    }
  }
  depends_on = [azapi_resource.container]
}


```



## Arguments Reference

The following arguments are supported:

* `type` - (Required) The type of the resource. This should be set to `Microsoft.Compute/galleries/applications/versions@api-version`. The available api-versions for this resource are: [`2019-03-01`, `2019-07-01`, `2019-12-01`, `2020-09-30`, `2021-07-01`, `2021-10-01`, `2022-01-03`, `2022-03-03`, `2022-08-03`, `2023-07-03`, `2024-03-03`].

* `parent_id` - (Required) The ID of the azure resource in which this resource is created. The allowed values are:  
  `/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Compute/galleries/{resourceName}/applications/{resourceName}`

* `name` - (Required) Specifies the name of the azure resource. Changing this forces a new resource to be created.

* `body` - (Required) Specifies the configuration of the resource. More information about the arguments in `body` can be found in the [Microsoft documentation](https://learn.microsoft.com/en-us/azure/templates/Microsoft.Compute/galleries/applications/versions?pivots=deployment-language-terraform).

For other arguments, please refer to the [azapi_resource](https://registry.terraform.io/providers/Azure/azapi/latest/docs/resources/resource) documentation.

## Import

 ```shell
 # Azure resource can be imported using the resource id, e.g.
 terraform import azapi_resource.example /subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Compute/galleries/{resourceName}/applications/{resourceName}/versions/{resourceName}
 
 # It also supports specifying API version by using the resource id with api-version as a query parameter, e.g.
 terraform import azapi_resource.example /subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Compute/galleries/{resourceName}/applications/{resourceName}/versions/{resourceName}?api-version=2024-03-03
 ```
