---
subcategory: "Microsoft.Network - Various Networking Services"
page_title: "networkManagers/routingConfigurations/ruleCollections/rules"
description: |-
  Manages a Network Manager Routing Rule.
---

# Microsoft.Network/networkManagers/routingConfigurations/ruleCollections/rules - Network Manager Routing Rule

This article demonstrates how to use `azapi` provider to manage the Network Manager Routing Rule resource in Azure.

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

data "azapi_client_config" "current" {
}

resource "azapi_resource" "resourceGroup" {
  type     = "Microsoft.Resources/resourceGroups@2020-06-01"
  name     = var.resource_name
  location = var.location
}

resource "azapi_resource" "networkManager" {
  type      = "Microsoft.Network/networkManagers@2024-05-01"
  parent_id = azapi_resource.resourceGroup.id
  name      = var.resource_name
  location  = var.location
  body = {
    properties = {
      description = ""
      networkManagerScopeAccesses = [
        "Routing",
      ]
      networkManagerScopes = {
        managementGroups = [
        ]
        subscriptions = [
          "/subscriptions/${data.azapi_client_config.current.subscription_id}",
        ]
      }
    }
  }
  schema_validation_enabled = false
  response_export_values    = ["*"]
}

resource "azapi_resource" "networkGroup" {
  type      = "Microsoft.Network/networkManagers/networkGroups@2024-05-01"
  parent_id = azapi_resource.networkManager.id
  name      = var.resource_name
  body = {
    properties = {
    }
  }
  schema_validation_enabled = false
  response_export_values    = ["*"]
}

resource "azapi_resource" "virtualNetwork" {
  type      = "Microsoft.Network/virtualNetworks@2024-05-01"
  parent_id = azapi_resource.resourceGroup.id
  name      = var.resource_name
  location  = var.location
  body = {
    properties = {
      addressSpace = {
        addressPrefixes = [
          "10.0.0.0/22",
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

resource "azapi_resource" "staticMember" {
  type      = "Microsoft.Network/networkManagers/networkGroups/staticMembers@2024-05-01"
  parent_id = azapi_resource.networkGroup.id
  name      = var.resource_name
  body = {
    properties = {
      resourceId = azapi_resource.virtualNetwork.id
    }
  }
  schema_validation_enabled = false
  response_export_values    = ["*"]
}

resource "azapi_resource" "routingConfiguration" {
  type      = "Microsoft.Network/networkManagers/routingConfigurations@2024-05-01"
  parent_id = azapi_resource.networkManager.id
  name      = var.resource_name
  body = {
    properties = {
      description = "example routing configuration"
    }
  }
  schema_validation_enabled = false
  response_export_values    = ["*"]
}

resource "azapi_resource" "ruleCollection" {
  type      = "Microsoft.Network/networkManagers/routingConfigurations/ruleCollections@2024-05-01"
  parent_id = azapi_resource.routingConfiguration.id
  name      = var.resource_name
  body = {
    properties = {
      description = "example rule collection"
      appliesTo = [
        {
          networkGroupId = azapi_resource.networkGroup.id
        },
      ]
    }
  }
  schema_validation_enabled = false
  response_export_values    = ["*"]
}

resource "azapi_resource" "rule" {
  type      = "Microsoft.Network/networkManagers/routingConfigurations/ruleCollections/rules@2024-05-01"
  parent_id = azapi_resource.ruleCollection.id
  name      = var.resource_name
  body = {
    kind = "Custom"
    properties = {
      description = "example rule"
      destination = {
        type               = "AddressPrefix" # Required, possible values are "AddressPrefix", "ServiceTag"
        destinationAddress = "10.0.0.0/16"   # Required
      }
      nextHop = {
        nextHopType    = "VirtualNetworkGateway" # Required, possible values are "Internet", "NoNextHop", "VirtualAppliance", "VirtualNetworkGateway", "VnetLocal"
        nextHopAddress = ""                      # Required if nextHopType is "VirtualAppliance"
      }
    }
  }
  schema_validation_enabled = false
  response_export_values    = ["*"]
}

```



## Arguments Reference

The following arguments are supported:

* `type` - (Required) The type of the resource. This should be set to `Microsoft.Network/networkManagers/routingConfigurations/ruleCollections/rules@api-version`. The available api-versions for this resource are: [`2023-03-01-preview`, `2024-03-01`, `2024-05-01`].

* `parent_id` - (Required) The ID of the azure resource in which this resource is created. The allowed values are:  
  `/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Network/networkManagers/{resourceName}/routingConfigurations/{resourceName}/ruleCollections/{resourceName}`

* `name` - (Required) Specifies the name of the azure resource. Changing this forces a new resource to be created.

* `body` - (Required) Specifies the configuration of the resource. More information about the arguments in `body` can be found in the [Microsoft documentation](https://learn.microsoft.com/en-us/azure/templates/Microsoft.Network/networkManagers/routingConfigurations/ruleCollections/rules?pivots=deployment-language-terraform).

For other arguments, please refer to the [azapi_resource](https://registry.terraform.io/providers/Azure/azapi/latest/docs/resources/resource) documentation.

## Import

 ```shell
 # Azure resource can be imported using the resource id, e.g.
 terraform import azapi_resource.example /subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Network/networkManagers/{resourceName}/routingConfigurations/{resourceName}/ruleCollections/{resourceName}/rules/{resourceName}
 
 # It also supports specifying API version by using the resource id with api-version as a query parameter, e.g.
 terraform import azapi_resource.example /subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Network/networkManagers/{resourceName}/routingConfigurations/{resourceName}/ruleCollections/{resourceName}/rules/{resourceName}?api-version=2024-05-01
 ```
