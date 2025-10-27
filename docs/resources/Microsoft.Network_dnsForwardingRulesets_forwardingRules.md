---
subcategory: "Microsoft.Network - Various Networking Services"
page_title: "dnsForwardingRulesets/forwardingRules"
description: |-
  Manages a Private DNS Resolver Forwarding Rule.
---

# Microsoft.Network/dnsForwardingRulesets/forwardingRules - Private DNS Resolver Forwarding Rule

This article demonstrates how to use `azapi` provider to manage the Private DNS Resolver Forwarding Rule resource in Azure.



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
          "10.0.0.0/16",
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

resource "azapi_resource" "dnsResolver" {
  type      = "Microsoft.Network/dnsResolvers@2022-07-01"
  parent_id = azapi_resource.resourceGroup.id
  name      = var.resource_name
  location  = var.location
  body = {
    properties = {
      virtualNetwork = {
        id = azapi_resource.virtualNetwork.id
      }
    }

  }
  schema_validation_enabled = false
  response_export_values    = ["*"]
}

resource "azapi_resource" "subnet" {
  type      = "Microsoft.Network/virtualNetworks/subnets@2022-07-01"
  parent_id = azapi_resource.virtualNetwork.id
  name      = "outbounddns"
  body = {
    properties = {
      addressPrefix = "10.0.0.64/28"
      delegations = [
        {
          name = "Microsoft.Network.dnsResolvers"
          properties = {
            serviceName = "Microsoft.Network/dnsResolvers"
          }
        },
      ]
      privateEndpointNetworkPolicies    = "Enabled"
      privateLinkServiceNetworkPolicies = "Enabled"
      serviceEndpointPolicies = [
      ]
      serviceEndpoints = [
      ]
    }
  }
  schema_validation_enabled = false
  response_export_values    = ["*"]
}

resource "azapi_resource" "outboundEndpoint" {
  type      = "Microsoft.Network/dnsResolvers/outboundEndpoints@2022-07-01"
  parent_id = azapi_resource.dnsResolver.id
  name      = var.resource_name
  location  = var.location
  body = {
    properties = {
      subnet = {
        id = azapi_resource.subnet.id
      }
    }

  }
  schema_validation_enabled = false
  response_export_values    = ["*"]
}

resource "azapi_resource" "dnsForwardingRuleset" {
  type      = "Microsoft.Network/dnsForwardingRulesets@2022-07-01"
  parent_id = azapi_resource.resourceGroup.id
  name      = var.resource_name
  location  = var.location
  body = {
    properties = {
      dnsResolverOutboundEndpoints = [
        {
          id = azapi_resource.outboundEndpoint.id
        },
      ]
    }

  }
  schema_validation_enabled = false
  response_export_values    = ["*"]
}

resource "azapi_resource" "forwardingRule" {
  type      = "Microsoft.Network/dnsForwardingRulesets/forwardingRules@2022-07-01"
  parent_id = azapi_resource.dnsForwardingRuleset.id
  name      = var.resource_name
  body = {
    properties = {
      domainName          = "onprem.local."
      forwardingRuleState = "Enabled"
      metadata            = null
      targetDnsServers = [
        {
          ipAddress = "10.10.0.1"
          port      = 53
        },
      ]
    }
  }
  schema_validation_enabled = false
  response_export_values    = ["*"]
}


```



## Arguments Reference

The following arguments are supported:

* `type` - (Required) The type of the resource. This should be set to `Microsoft.Network/dnsForwardingRulesets/forwardingRules@api-version`. The available api-versions for this resource are: [`2020-04-01-preview`, `2022-07-01`, `2023-07-01-preview`, `2025-05-01`].

* `parent_id` - (Required) The ID of the azure resource in which this resource is created. The allowed values are:  
  `/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Network/dnsForwardingRulesets/{resourceName}`

* `name` - (Required) Specifies the name of the azure resource. Changing this forces a new resource to be created.

* `body` - (Required) Specifies the configuration of the resource. More information about the arguments in `body` can be found in the [Microsoft documentation](https://learn.microsoft.com/en-us/azure/templates/Microsoft.Network/dnsForwardingRulesets/forwardingRules?pivots=deployment-language-terraform).

For other arguments, please refer to the [azapi_resource](https://registry.terraform.io/providers/Azure/azapi/latest/docs/resources/resource) documentation.

## Import

 ```shell
 # Azure resource can be imported using the resource id, e.g.
 terraform import azapi_resource.example /subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Network/dnsForwardingRulesets/{resourceName}/forwardingRules/{resourceName}
 
 # It also supports specifying API version by using the resource id with api-version as a query parameter, e.g.
 terraform import azapi_resource.example /subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Network/dnsForwardingRulesets/{resourceName}/forwardingRules/{resourceName}?api-version=2025-05-01
 ```
