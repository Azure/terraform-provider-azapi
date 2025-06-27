---
subcategory: "Microsoft.Network - Various Networking Services"
page_title: "privateEndpoints"
description: |-
  Manages a Private Endpoint.
---

# Microsoft.Network/privateEndpoints - Private Endpoint

This article demonstrates how to use `azapi` provider to manage the Private Endpoint resource in Azure.

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

resource "azapi_resource" "virtualNetwork" {
  type      = "Microsoft.Network/virtualNetworks@2022-07-01"
  parent_id = azapi_resource.resourceGroup.id
  name      = var.resource_name
  location  = var.location
  body = {
    properties = {
      addressSpace = {
        addressPrefixes = [
          "10.5.0.0/16",
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

resource "azapi_resource" "publicIPAddress" {
  type      = "Microsoft.Network/publicIPAddresses@2022-07-01"
  parent_id = azapi_resource.resourceGroup.id
  name      = var.resource_name
  location  = var.location
  body = {
    properties = {
      ddosSettings = {
        protectionMode = "VirtualNetworkInherited"
      }
      idleTimeoutInMinutes     = 4
      publicIPAddressVersion   = "IPv4"
      publicIPAllocationMethod = "Static"
    }
    sku = {
      name = "Standard"
      tier = "Regional"
    }
  }
  schema_validation_enabled = false
  response_export_values    = ["*"]
}

resource "azapi_resource" "loadBalancer" {
  type      = "Microsoft.Network/loadBalancers@2022-07-01"
  parent_id = azapi_resource.resourceGroup.id
  name      = var.resource_name
  location  = var.location
  body = {
    properties = {
      frontendIPConfigurations = [
        {
          name = var.resource_name
          properties = {
            publicIPAddress = {
              id = azapi_resource.publicIPAddress.id
            }
          }
        },
      ]
    }
    sku = {
      name = "Standard"
      tier = "Regional"
    }
  }
  schema_validation_enabled = false
  response_export_values    = ["*"]
}

resource "azapi_resource" "subnet" {
  type      = "Microsoft.Network/virtualNetworks/subnets@2022-07-01"
  parent_id = azapi_resource.virtualNetwork.id
  name      = var.resource_name
  body = {
    properties = {
      addressPrefix = "10.5.4.0/24"
      delegations = [
      ]
      privateEndpointNetworkPolicies    = "Enabled"
      privateLinkServiceNetworkPolicies = "Disabled"
      serviceEndpointPolicies = [
      ]
      serviceEndpoints = [
      ]
    }
  }
  schema_validation_enabled = false
  response_export_values    = ["*"]
}

resource "azapi_resource" "privateLinkService" {
  type      = "Microsoft.Network/privateLinkServices@2022-07-01"
  parent_id = azapi_resource.resourceGroup.id
  name      = var.resource_name
  location  = var.location
  body = {
    properties = {
      autoApproval = {
        subscriptions = [
        ]
      }
      enableProxyProtocol = false
      fqdns = [
      ]
      ipConfigurations = [
        {
          name = "primaryIpConfiguration-230630033653892379"
          properties = {
            primary                   = true
            privateIPAddress          = ""
            privateIPAddressVersion   = "IPv4"
            privateIPAllocationMethod = "Dynamic"
            subnet = {
              id = azapi_resource.subnet.id
            }
          }
        },
      ]
      loadBalancerFrontendIpConfigurations = [
        {
          id = azapi_resource.loadBalancer.output.properties.frontendIPConfigurations[0].id
        },
      ]
      visibility = {
        subscriptions = [
        ]
      }
    }
  }
  schema_validation_enabled = false
  response_export_values    = ["*"]
}

resource "azapi_resource" "privateEndpoint" {
  type      = "Microsoft.Network/privateEndpoints@2022-07-01"
  parent_id = azapi_resource.resourceGroup.id
  name      = var.resource_name
  location  = var.location
  body = {
    properties = {
      privateLinkServiceConnections = [
        {
          name = azapi_resource.privateLinkService.name
          properties = {
            privateLinkServiceId = azapi_resource.privateLinkService.id
          }
        }
      ]
      subnet = {
        id = azapi_resource.subnet.id
      }
    }
  }
  schema_validation_enabled = false
  response_export_values    = ["*"]
}

```



## Arguments Reference

The following arguments are supported:

* `type` - (Required) The type of the resource. This should be set to `Microsoft.Network/privateEndpoints@api-version`. The available api-versions for this resource are: [`2019-04-01`, `2019-06-01`, `2019-07-01`, `2019-08-01`, `2019-09-01`, `2019-11-01`, `2019-12-01`, `2020-03-01`, `2020-04-01`, `2020-05-01`, `2020-06-01`, `2020-07-01`, `2020-08-01`, `2020-11-01`, `2021-02-01`, `2021-03-01`, `2021-05-01`, `2021-08-01`, `2022-01-01`, `2022-05-01`, `2022-07-01`, `2022-09-01`, `2022-11-01`, `2023-02-01`, `2023-04-01`, `2023-05-01`, `2023-06-01`, `2023-09-01`, `2023-11-01`, `2024-01-01`, `2024-03-01`, `2024-05-01`, `2024-07-01`].

* `parent_id` - (Required) The ID of the azure resource in which this resource is created. The allowed values are:  
  `/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}`

* `name` - (Required) Specifies the name of the azure resource. Changing this forces a new resource to be created.

* `body` - (Required) Specifies the configuration of the resource. More information about the arguments in `body` can be found in the [Microsoft documentation](https://learn.microsoft.com/en-us/azure/templates/Microsoft.Network/privateEndpoints?pivots=deployment-language-terraform).

For other arguments, please refer to the [azapi_resource](https://registry.terraform.io/providers/Azure/azapi/latest/docs/resources/resource) documentation.

## Import

 ```shell
 # Azure resource can be imported using the resource id, e.g.
 terraform import azapi_resource.example /subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Network/privateEndpoints/{resourceName}
 
 # It also supports specifying API version by using the resource id with api-version as a query parameter, e.g.
 terraform import azapi_resource.example /subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Network/privateEndpoints/{resourceName}?api-version=2024-07-01
 ```
