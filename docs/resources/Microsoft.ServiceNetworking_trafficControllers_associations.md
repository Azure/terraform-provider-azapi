---
subcategory: "Microsoft.ServiceNetworking - Azure Service Networking"
page_title: "trafficControllers/associations"
description: |-
  Manages a association between an Application Gateway for Containers and a Subnet.
---

# Microsoft.ServiceNetworking/trafficControllers/associations - association between an Application Gateway for Containers and a Subnet

This article demonstrates how to use `azapi` provider to manage the association between an Application Gateway for Containers and a Subnet resource in Azure.

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

resource "azapi_resource" "trafficController" {
  type      = "Microsoft.ServiceNetworking/trafficControllers@2023-11-01"
  parent_id = azapi_resource.resourceGroup.id
  name      = "${var.resource_name}-tc"
  location  = var.location
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
      addressPrefix         = "10.0.1.0/24"
      defaultOutboundAccess = true
      delegations = [{
        name = "delegation"
        properties = {
          serviceName = "Microsoft.ServiceNetworking/trafficControllers"
        }
      }]
      privateEndpointNetworkPolicies    = "Disabled"
      privateLinkServiceNetworkPolicies = "Enabled"
      serviceEndpointPolicies           = []
      serviceEndpoints                  = []
    }
  }
}

resource "azapi_resource" "association" {
  type      = "Microsoft.ServiceNetworking/trafficControllers/associations@2023-11-01"
  parent_id = azapi_resource.trafficController.id
  name      = "${var.resource_name}-assoc"
  location  = var.location
  body = {
    properties = {
      associationType = "subnets"
      subnet = {
        id = azapi_resource.subnet.id
      }
    }
  }
}


```



## Arguments Reference

The following arguments are supported:

* `type` - (Required) The type of the resource. This should be set to `Microsoft.ServiceNetworking/trafficControllers/associations@api-version`. The available api-versions for this resource are: [`2022-10-01-preview`, `2023-05-01-preview`, `2023-11-01`, `2024-05-01-preview`, `2025-01-01`, `2025-03-01-preview`].

* `parent_id` - (Required) The ID of the azure resource in which this resource is created. The allowed values are:  
  `/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.ServiceNetworking/trafficControllers/{resourceName}`

* `name` - (Required) Specifies the name of the azure resource. Changing this forces a new resource to be created.

* `body` - (Required) Specifies the configuration of the resource. More information about the arguments in `body` can be found in the [Microsoft documentation](https://learn.microsoft.com/en-us/azure/templates/Microsoft.ServiceNetworking/trafficControllers/associations?pivots=deployment-language-terraform).

For other arguments, please refer to the [azapi_resource](https://registry.terraform.io/providers/Azure/azapi/latest/docs/resources/resource) documentation.

## Import

 ```shell
 # Azure resource can be imported using the resource id, e.g.
 terraform import azapi_resource.example /subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.ServiceNetworking/trafficControllers/{resourceName}/associations/{resourceName}
 
 # It also supports specifying API version by using the resource id with api-version as a query parameter, e.g.
 terraform import azapi_resource.example /subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.ServiceNetworking/trafficControllers/{resourceName}/associations/{resourceName}?api-version=2025-03-01-preview
 ```
