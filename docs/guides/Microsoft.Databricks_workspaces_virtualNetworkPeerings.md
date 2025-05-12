---
subcategory: "Microsoft.Databricks - Azure Databricks"
page_title: "workspaces/virtualNetworkPeerings"
description: |-
  Manages a Databricks Virtual Network Peering.
---

# Microsoft.Databricks/workspaces/virtualNetworkPeerings - Databricks Virtual Network Peering

This article demonstrates how to use `azapi` provider to manage the Databricks Virtual Network Peering resource in Azure.

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

data "azapi_resource_id" "workspace_resource_group" {
  type      = "Microsoft.Resources/resourceGroups@2020-06-01"
  parent_id = azapi_resource.resourceGroup.parent_id
  name      = "databricks-rg-${var.resource_name}"
}

resource "azapi_resource" "workspace" {
  type      = "Microsoft.Databricks/workspaces@2023-02-01"
  parent_id = azapi_resource.resourceGroup.id
  name      = var.resource_name
  location  = var.location
  body = {
    properties = {
      managedResourceGroupId = data.azapi_resource_id.workspace_resource_group.id
      publicNetworkAccess    = "Enabled"
    }
    sku = {
      name = "standard"
    }
  }
  schema_validation_enabled = false
  response_export_values    = ["*"]
}

resource "azapi_resource" "virtualNetwork" {
  type      = "Microsoft.Network/virtualNetworks@2022-07-01"
  parent_id = azapi_resource.resourceGroup.id
  name      = var.resource_name
  location  = var.location
  body = {
    properties = {
      addressSpace = {
        addressPrefixes = [
          "10.0.1.0/24",
        ]
      }
      dhcpOptions = {
        dnsServers = [
        ]
      }
      subnets = [
      ]
    }
  }
  schema_validation_enabled = false
  response_export_values    = ["*"]
  lifecycle {
    ignore_changes = [body.properties.subnets]
  }
}

resource "azapi_resource" "virtualNetworkPeering" {
  type      = "Microsoft.Databricks/workspaces/virtualNetworkPeerings@2023-02-01"
  parent_id = azapi_resource.workspace.id
  name      = var.resource_name
  body = {
    properties = {
      allowForwardedTraffic     = false
      allowGatewayTransit       = false
      allowVirtualNetworkAccess = true
      databricksAddressSpace = {
        addressPrefixes = [
          "10.139.0.0/16",
        ]
      }
      remoteAddressSpace = {
        addressPrefixes = [
          "10.0.1.0/24",
        ]
      }
      remoteVirtualNetwork = {
        id = azapi_resource.virtualNetwork.id
      }
      useRemoteGateways = false
    }
  }
  schema_validation_enabled = false
  response_export_values    = ["*"]
}


```



## Arguments Reference

The following arguments are supported:

* `type` - (Required) The type of the resource. This should be set to `Microsoft.Databricks/workspaces/virtualNetworkPeerings@api-version`. The available api-versions for this resource are: [`2018-04-01`, `2021-04-01-preview`, `2022-04-01-preview`, `2023-02-01`, `2023-09-15-preview`, `2024-05-01`, `2024-09-01-preview`].

* `parent_id` - (Required) The ID of the azure resource in which this resource is created. The allowed values are:  
  `/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Databricks/workspaces/{resourceName}`

* `name` - (Required) Specifies the name of the azure resource. Changing this forces a new resource to be created.

* `body` - (Required) Specifies the configuration of the resource. More information about the arguments in `body` can be found in the [Microsoft documentation](https://learn.microsoft.com/en-us/azure/templates/Microsoft.Databricks/workspaces/virtualNetworkPeerings?pivots=deployment-language-terraform).

For other arguments, please refer to the [azapi_resource](https://registry.terraform.io/providers/Azure/azapi/latest/docs/resources/resource) documentation.

## Import

 ```shell
 # Azure resource can be imported using the resource id, e.g.
 terraform import azapi_resource.example /subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Databricks/workspaces/{resourceName}/virtualNetworkPeerings/{resourceName}
 
 # It also supports specifying API version by using the resource id with api-version as a query parameter, e.g.
 terraform import azapi_resource.example /subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Databricks/workspaces/{resourceName}/virtualNetworkPeerings/{resourceName}?api-version=2024-09-01-preview
 ```
