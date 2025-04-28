---
subcategory: "Microsoft.MobileNetwork - Azure Private 5G Core"
page_title: "packetCoreControlPlanes/packetCoreDataPlanes"
description: |-
  Manages a Mobile Network Packet Core Data Plane.
---

# Microsoft.MobileNetwork/packetCoreControlPlanes/packetCoreDataPlanes - Mobile Network Packet Core Data Plane

This article demonstrates how to use `azapi` provider to manage the Mobile Network Packet Core Data Plane resource in Azure.

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
  default = "eastus"
}

resource "azapi_resource" "resourceGroup" {
  type     = "Microsoft.Resources/resourceGroups@2020-06-01"
  name     = var.resource_name
  location = var.location
}

resource "azapi_resource" "mobileNetwork" {
  type      = "Microsoft.MobileNetwork/mobileNetworks@2022-11-01"
  parent_id = azapi_resource.resourceGroup.id
  name      = var.resource_name
  location  = var.location
  body = {
    properties = {
      publicLandMobileNetworkIdentifier = {
        mcc = "001"
        mnc = "01"
      }
    }

  }
  schema_validation_enabled = false
  response_export_values    = ["*"]
}

resource "azapi_resource" "dataBoxEdgeDevice" {
  type      = "Microsoft.DataBoxEdge/dataBoxEdgeDevices@2022-03-01"
  parent_id = azapi_resource.resourceGroup.id
  name      = var.resource_name
  location  = var.location
  body = {
    sku = {
      name = "EdgeP_Base"
      tier = "Standard"
    }

  }
  schema_validation_enabled = false
  response_export_values    = ["*"]
}

resource "azapi_resource" "site" {
  type      = "Microsoft.MobileNetwork/mobileNetworks/sites@2022-11-01"
  parent_id = azapi_resource.mobileNetwork.id
  name      = var.resource_name
  location  = var.location
  body = {
    properties = {
    }

  }
  schema_validation_enabled = false
  response_export_values    = ["*"]
}

resource "azapi_resource" "packetCoreControlPlane" {
  type      = "Microsoft.MobileNetwork/packetCoreControlPlanes@2022-11-01"
  parent_id = azapi_resource.resourceGroup.id
  name      = var.resource_name
  location  = var.location
  body = {
    properties = {
      controlPlaneAccessInterface = {
      }
      localDiagnosticsAccess = {
        authenticationType = "AAD"
      }
      platform = {
        azureStackEdgeDevice = {
          id = azapi_resource.dataBoxEdgeDevice.id
        }
        type = "AKS-HCI"
      }
      sites = [
        {
          id = azapi_resource.site.id
        },
      ]
      sku   = "G0"
      ueMtu = 1440
    }

  }
  schema_validation_enabled = false
  response_export_values    = ["*"]
}

resource "azapi_resource" "packetCoreDataPlane" {
  type      = "Microsoft.MobileNetwork/packetCoreControlPlanes/packetCoreDataPlanes@2022-11-01"
  parent_id = azapi_resource.packetCoreControlPlane.id
  name      = var.resource_name
  location  = var.location
  body = {
    properties = {
      userPlaneAccessInterface = {
      }
    }

  }
  schema_validation_enabled = false
  response_export_values    = ["*"]
}


```



## Arguments Reference

The following arguments are supported:

* `type` - (Required) The type of the resource. This should be set to `Microsoft.MobileNetwork/packetCoreControlPlanes/packetCoreDataPlanes@api-version`. The available api-versions for this resource are: [`2022-03-01-preview`, `2022-04-01-preview`, `2022-11-01`, `2023-06-01`, `2023-09-01`, `2024-02-01`, `2024-04-01`].

* `parent_id` - (Required) The ID of the azure resource in which this resource is created. The allowed values are:  
  `/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.MobileNetwork/packetCoreControlPlanes/{resourceName}`

* `name` - (Required) Specifies the name of the azure resource. Changing this forces a new resource to be created.

* `body` - (Required) Specifies the configuration of the resource. More information about the arguments in `body` can be found in the [Microsoft documentation](https://learn.microsoft.com/en-us/azure/templates/Microsoft.MobileNetwork/packetCoreControlPlanes/packetCoreDataPlanes?pivots=deployment-language-terraform).

## Import

 ```shell
 # Azure resource can be imported using the resource id, e.g.
 terraform import azapi_resource.example /subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.MobileNetwork/packetCoreControlPlanes/{resourceName}/packetCoreDataPlanes/{resourceName}
 
 # It also supports specifying API version by using the resource id with api-version as a query parameter, e.g.
 terraform import azapi_resource.example /subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.MobileNetwork/packetCoreControlPlanes/{resourceName}/packetCoreDataPlanes/{resourceName}?api-version=2024-04-01
 ```
