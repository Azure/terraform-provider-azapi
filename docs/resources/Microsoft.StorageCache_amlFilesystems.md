---
subcategory: "Microsoft.StorageCache - Azure HPC Cache"
page_title: "amlFilesystems"
description: |-
  Manages a Azure Managed Lustre File System.
---

# Microsoft.StorageCache/amlFilesystems - Azure Managed Lustre File System

This article demonstrates how to use `azapi` provider to manage the Azure Managed Lustre File System resource in Azure.

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

resource "azapi_resource" "virtualNetwork" {
  type      = "Microsoft.Network/virtualNetworks@2024-05-01"
  parent_id = azapi_resource.resourceGroup.id
  name      = "${var.resource_name}-vnet"
  location  = var.location
  body = {
    properties = {
      addressSpace = {
        addressPrefixes = ["10.0.0.0/16"]
      }
      dhcpOptions = {
        dnsServers = []
      }
      privateEndpointVNetPolicies = "Disabled"
      subnets                     = []
    }
  }
}

resource "azapi_resource" "subnet" {
  type      = "Microsoft.Network/virtualNetworks/subnets@2024-05-01"
  parent_id = azapi_resource.virtualNetwork.id
  name      = "${var.resource_name}-subnet"
  body = {
    properties = {
      addressPrefix                     = "10.0.2.0/24"
      defaultOutboundAccess             = true
      delegations                       = []
      privateEndpointNetworkPolicies    = "Disabled"
      privateLinkServiceNetworkPolicies = "Enabled"
      serviceEndpointPolicies           = []
      serviceEndpoints                  = []
    }
  }
}

resource "azapi_resource" "amlFilesystem" {
  type      = "Microsoft.StorageCache/amlFilesystems@2024-07-01"
  parent_id = azapi_resource.resourceGroup.id
  name      = "${var.resource_name}-amlfs"
  location  = var.location
  body = {
    properties = {
      filesystemSubnet = azapi_resource.subnet.id
      maintenanceWindow = {
        dayOfWeek    = "Friday"
        timeOfDayUTC = "22:00"
      }
      storageCapacityTiB = 8
    }
    sku = {
      name = "AMLFS-Durable-Premium-250"
    }
    zones = ["1"]
  }
}


```



## Arguments Reference

The following arguments are supported:

* `type` - (Required) The type of the resource. This should be set to `Microsoft.StorageCache/amlFilesystems@api-version`. The available api-versions for this resource are: [`2023-03-01-preview`, `2023-05-01`, `2023-11-01-preview`, `2024-03-01`, `2024-07-01`, `2025-07-01`].

* `parent_id` - (Required) The ID of the azure resource in which this resource is created. The allowed values are:  
  `/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}`

* `name` - (Required) Specifies the name of the azure resource. Changing this forces a new resource to be created.

* `body` - (Required) Specifies the configuration of the resource. More information about the arguments in `body` can be found in the [Microsoft documentation](https://learn.microsoft.com/en-us/azure/templates/Microsoft.StorageCache/amlFilesystems?pivots=deployment-language-terraform).

For other arguments, please refer to the [azapi_resource](https://registry.terraform.io/providers/Azure/azapi/latest/docs/resources/resource) documentation.

## Import

 ```shell
 # Azure resource can be imported using the resource id, e.g.
 terraform import azapi_resource.example /subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.StorageCache/amlFilesystems/{resourceName}
 
 # It also supports specifying API version by using the resource id with api-version as a query parameter, e.g.
 terraform import azapi_resource.example /subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.StorageCache/amlFilesystems/{resourceName}?api-version=2025-07-01
 ```
