---
subcategory: "Microsoft.Network - Various Networking Services"
page_title: "virtualNetworks/virtualNetworkPeerings"
description: |-
  Manages a virtual network peering which allows resources to access other.
---

# Microsoft.Network/virtualNetworks/virtualNetworkPeerings - virtual network peering which allows resources to access other

This article demonstrates how to use `azapi` provider to manage the virtual network peering which allows resources to access other resource in Azure.

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
      parameters = {
        prepareEncryption = {
          value = false
        }
        requireInfrastructureEncryption = {
          value = false
        }
      }
      publicNetworkAccess = "Enabled"
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
          "10.139.0.0/16"
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

* `type` - (Required) The type of the resource. This should be set to `Microsoft.Network/virtualNetworks/virtualNetworkPeerings@api-version`. The available api-versions for this resource are: [`2016-06-01`, `2016-09-01`, `2016-12-01`, `2017-03-01`, `2017-03-30`, `2017-06-01`, `2017-08-01`, `2017-09-01`, `2017-10-01`, `2017-11-01`, `2018-01-01`, `2018-02-01`, `2018-04-01`, `2018-06-01`, `2018-07-01`, `2018-08-01`, `2018-10-01`, `2018-11-01`, `2018-12-01`, `2019-02-01`, `2019-04-01`, `2019-06-01`, `2019-07-01`, `2019-08-01`, `2019-09-01`, `2019-11-01`, `2019-12-01`, `2020-03-01`, `2020-04-01`, `2020-05-01`, `2020-06-01`, `2020-07-01`, `2020-08-01`, `2020-11-01`, `2021-02-01`, `2021-03-01`, `2021-05-01`, `2021-08-01`, `2022-01-01`, `2022-05-01`, `2022-07-01`, `2022-09-01`, `2022-11-01`, `2023-02-01`, `2023-04-01`, `2023-05-01`, `2023-06-01`, `2023-09-01`, `2023-11-01`, `2024-01-01`, `2024-03-01`, `2024-05-01`, `2024-07-01`, `2024-10-01`].

* `parent_id` - (Required) The ID of the azure resource in which this resource is created. The allowed values are:  
  `/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Network/virtualNetworks/{resourceName}`

* `name` - (Required) Specifies the name of the azure resource. Changing this forces a new resource to be created.

* `body` - (Required) Specifies the configuration of the resource. More information about the arguments in `body` can be found in the [Microsoft documentation](https://learn.microsoft.com/en-us/azure/templates/Microsoft.Network/virtualNetworks/virtualNetworkPeerings?pivots=deployment-language-terraform).

For other arguments, please refer to the [azapi_resource](https://registry.terraform.io/providers/Azure/azapi/latest/docs/resources/resource) documentation.

## Import

 ```shell
 # Azure resource can be imported using the resource id, e.g.
 terraform import azapi_resource.example /subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Network/virtualNetworks/{resourceName}/virtualNetworkPeerings/{resourceName}
 
 # It also supports specifying API version by using the resource id with api-version as a query parameter, e.g.
 terraform import azapi_resource.example /subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Network/virtualNetworks/{resourceName}/virtualNetworkPeerings/{resourceName}?api-version=2024-10-01
 ```
