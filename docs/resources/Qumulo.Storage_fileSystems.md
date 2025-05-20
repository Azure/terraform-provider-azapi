---
subcategory: "Qumulo.Storage - Qumulo Storage"
page_title: "fileSystems"
description: |-
  Manages a Qumulo File System.
---

# Qumulo.Storage/fileSystems - Qumulo File System

This article demonstrates how to use `azapi` provider to manage the Qumulo File System resource in Azure.

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

variable "qumulo_password" {
  type      = string
  default   = ")^X#ZX#JRyIY}t9"
  sensitive = true
}

resource "azapi_resource" "resourceGroup" {
  type     = "Microsoft.Resources/resourceGroups@2020-06-01"
  name     = var.resource_name
  location = var.location
}

resource "azapi_resource" "vnet" {
  type      = "Microsoft.Network/virtualNetworks@2024-05-01"
  name      = var.resource_name
  location  = var.location
  parent_id = azapi_resource.resourceGroup.id

  body = {
    properties = {
      addressSpace = {
        addressPrefixes = ["10.0.0.0/16"]
      }
      privateEndpointVNetPolicies = "Disabled"
      subnets                     = []
    }
  }

  schema_validation_enabled = false
  response_export_values    = ["*"]
  lifecycle {
    # This is to avoid vnet change to overwrite the subnets
    ignore_changes = [body.properties.subnets]
  }
}

resource "azapi_resource" "subnet" {
  type      = "Microsoft.Network/virtualNetworks/subnets@2024-05-01"
  name      = var.resource_name
  location  = var.location
  parent_id = azapi_resource.vnet.id

  body = {
    properties = {
      addressPrefix         = "10.0.1.0/24"
      defaultOutboundAccess = true
      delegations = [{
        name = "delegation"
        properties = {
          actions     = ["Microsoft.Network/virtualNetworks/subnets/join/action"]
          serviceName = "Qumulo.Storage/fileSystems"
        }
      }]
      privateEndpointNetworkPolicies    = "Disabled"
      privateLinkServiceNetworkPolicies = "Enabled"
    }
  }

  schema_validation_enabled = false
  response_export_values    = ["*"]
}

resource "azapi_resource" "qumuloFileSystem" {
  type      = "Qumulo.Storage/fileSystems@2024-06-19"
  name      = var.resource_name
  location  = var.location
  parent_id = azapi_resource.resourceGroup.id

  body = {
    properties = {
      adminPassword     = var.qumulo_password
      availabilityZone  = "1"
      delegatedSubnetId = azapi_resource.subnet.id
      marketplaceDetails = {
        offerId     = "qumulo-saas-mpp"
        planId      = "azure-native-qumulo-v3"
        publisherId = "qumulo1584033880660"
      }
      storageSku = "Cold_LRS"
      userDetails = {
        email = "test@test.com"
      }
    }
  }

  tags = {
    environment = "terraform-acctests"
    some_key    = "some-value"
  }

  schema_validation_enabled = false
  response_export_values    = ["*"]
}

```



## Arguments Reference

The following arguments are supported:

* `type` - (Required) The type of the resource. This should be set to `Qumulo.Storage/fileSystems@api-version`. The available api-versions for this resource are: [`2022-06-27-preview`, `2022-10-12`, `2022-10-12-preview`, `2024-06-19`].

* `parent_id` - (Required) The ID of the azure resource in which this resource is created. The allowed values are:  
  `/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}`

* `name` - (Required) Specifies the name of the azure resource. Changing this forces a new resource to be created.

* `body` - (Required) Specifies the configuration of the resource. More information about the arguments in `body` can be found in the [Microsoft documentation](https://learn.microsoft.com/en-us/azure/templates/Qumulo.Storage/fileSystems?pivots=deployment-language-terraform).

For other arguments, please refer to the [azapi_resource](https://registry.terraform.io/providers/Azure/azapi/latest/docs/resources/resource) documentation.

## Import

 ```shell
 # Azure resource can be imported using the resource id, e.g.
 terraform import azapi_resource.example /subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Qumulo.Storage/fileSystems/{resourceName}
 
 # It also supports specifying API version by using the resource id with api-version as a query parameter, e.g.
 terraform import azapi_resource.example /subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Qumulo.Storage/fileSystems/{resourceName}?api-version=2024-06-19
 ```
