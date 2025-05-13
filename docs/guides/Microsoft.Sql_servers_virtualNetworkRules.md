---
subcategory: "Microsoft.Sql - Azure SQL Database, Azure SQL Managed Instance, Azure Synapse Analytics"
page_title: "servers/virtualNetworkRules"
description: |-
  Manages a Azure SQL Virtual Network Rule.
---

# Microsoft.Sql/servers/virtualNetworkRules - Azure SQL Virtual Network Rule

This article demonstrates how to use `azapi` provider to manage the Azure SQL Virtual Network Rule resource in Azure.

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

resource "azapi_resource" "server" {
  type      = "Microsoft.Sql/servers@2021-02-01-preview"
  parent_id = azapi_resource.resourceGroup.id
  name      = var.resource_name
  location  = var.location
  body = {
    properties = {
      administratorLogin            = "missadmin"
      administratorLoginPassword    = "P@55W0rD!!oyefp"
      minimalTlsVersion             = "1.2"
      publicNetworkAccess           = "Enabled"
      restrictOutboundNetworkAccess = "Disabled"
      version                       = "12.0"
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
          "10.7.28.0/23",
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

resource "azapi_resource" "subnet" {
  type      = "Microsoft.Network/virtualNetworks/subnets@2022-07-01"
  parent_id = azapi_resource.virtualNetwork.id
  name      = var.resource_name
  body = {
    properties = {
      addressPrefix = "10.7.28.0/25"
      delegations = [
      ]
      privateEndpointNetworkPolicies    = "Enabled"
      privateLinkServiceNetworkPolicies = "Enabled"
      serviceEndpointPolicies = [
      ]
      serviceEndpoints = [
        {
          service = "Microsoft.Sql"
        },
      ]
    }
  }
  schema_validation_enabled = false
  response_export_values    = ["*"]
}

resource "azapi_resource" "virtualNetworkRule" {
  type      = "Microsoft.Sql/servers/virtualNetworkRules@2020-11-01-preview"
  parent_id = azapi_resource.server.id
  name      = var.resource_name
  body = {
    properties = {
      ignoreMissingVnetServiceEndpoint = false
      virtualNetworkSubnetId           = azapi_resource.subnet.id
    }
  }
  schema_validation_enabled = false
  response_export_values    = ["*"]
}


```



## Arguments Reference

The following arguments are supported:

* `type` - (Required) The type of the resource. This should be set to `Microsoft.Sql/servers/virtualNetworkRules@api-version`. The available api-versions for this resource are: [`2015-05-01-preview`, `2020-02-02-preview`, `2020-08-01-preview`, `2020-11-01-preview`, `2021-02-01-preview`, `2021-05-01-preview`, `2021-08-01-preview`, `2021-11-01`, `2021-11-01-preview`, `2022-02-01-preview`, `2022-05-01-preview`, `2022-08-01-preview`, `2022-11-01-preview`, `2023-02-01-preview`, `2023-05-01-preview`, `2023-08-01-preview`, `2024-05-01-preview`].

* `parent_id` - (Required) The ID of the azure resource in which this resource is created. The allowed values are:  
  `/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Sql/servers/{resourceName}`

* `name` - (Required) Specifies the name of the azure resource. Changing this forces a new resource to be created.

* `body` - (Required) Specifies the configuration of the resource. More information about the arguments in `body` can be found in the [Microsoft documentation](https://learn.microsoft.com/en-us/azure/templates/Microsoft.Sql/servers/virtualNetworkRules?pivots=deployment-language-terraform).

For other arguments, please refer to the [azapi_resource](https://registry.terraform.io/providers/Azure/azapi/latest/docs/resources/resource) documentation.

## Import

 ```shell
 # Azure resource can be imported using the resource id, e.g.
 terraform import azapi_resource.example /subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Sql/servers/{resourceName}/virtualNetworkRules/{resourceName}
 
 # It also supports specifying API version by using the resource id with api-version as a query parameter, e.g.
 terraform import azapi_resource.example /subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Sql/servers/{resourceName}/virtualNetworkRules/{resourceName}?api-version=2024-05-01-preview
 ```
