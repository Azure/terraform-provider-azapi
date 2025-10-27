---
subcategory: "Microsoft.Network - Various Networking Services"
page_title: "networkManagerConnections"
description: |-
  Manages a Network Manager Management Group Connection.
---

# Microsoft.Network/networkManagerConnections - Network Manager Management Group Connection

This article demonstrates how to use `azapi` provider to manage the Network Manager Management Group Connection resource in Azure.



## Example Usage

### default

```hcl
terraform {
  required_providers {
    azapi = {
      source = "Azure/azapi"
    }
    azurerm = {
      source = "hashicorp/azurerm"
    }
  }
}

provider "azurerm" {
  features {
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

data "azurerm_client_config" "current" {
}

data "azapi_resource" "subscription" {
  type                   = "Microsoft.Resources/subscriptions@2021-01-01"
  resource_id            = "/subscriptions/${data.azurerm_client_config.current.subscription_id}"
  response_export_values = ["*"]
}

resource "azapi_resource" "resourceGroup" {
  type     = "Microsoft.Resources/resourceGroups@2020-06-01"
  name     = var.resource_name
  location = var.location
}

resource "azapi_resource" "networkManager" {
  type      = "Microsoft.Network/networkManagers@2022-09-01"
  parent_id = azapi_resource.resourceGroup.id
  name      = var.resource_name
  location  = var.location
  body = {
    properties = {
      description = ""
      networkManagerScopeAccesses = [
        "SecurityAdmin",
      ]
      networkManagerScopes = {
        managementGroups = [
        ]
        subscriptions = [
          data.azapi_resource.subscription.id,
        ]
      }
    }
  }
  schema_validation_enabled = false
  response_export_values    = ["*"]
}

resource "azapi_resource" "networkManagerConnection" {
  type      = "Microsoft.Network/networkManagerConnections@2022-09-01"
  parent_id = data.azapi_resource.subscription.id
  name      = var.resource_name
  body = {
    properties = {
      networkManagerId = azapi_resource.networkManager.id
    }
  }
  schema_validation_enabled = false
  response_export_values    = ["*"]
}


```



## Arguments Reference

The following arguments are supported:

* `type` - (Required) The type of the resource. This should be set to `Microsoft.Network/networkManagerConnections@api-version`. The available api-versions for this resource are: [`2021-05-01-preview`, `2022-01-01`, `2022-02-01-preview`, `2022-04-01-preview`, `2022-05-01`, `2022-07-01`, `2022-09-01`, `2022-11-01`, `2023-02-01`, `2023-04-01`, `2023-05-01`, `2023-06-01`, `2023-09-01`, `2023-11-01`, `2024-01-01`, `2024-03-01`, `2024-05-01`, `2024-07-01`, `2024-10-01`].

* `parent_id` - (Required) The ID of the azure resource in which this resource is created. The allowed values are:  
  `/subscriptions/{subscriptionId}`

* `name` - (Required) Specifies the name of the azure resource. Changing this forces a new resource to be created.

* `body` - (Required) Specifies the configuration of the resource. More information about the arguments in `body` can be found in the [Microsoft documentation](https://learn.microsoft.com/en-us/azure/templates/Microsoft.Network/networkManagerConnections?pivots=deployment-language-terraform).

For other arguments, please refer to the [azapi_resource](https://registry.terraform.io/providers/Azure/azapi/latest/docs/resources/resource) documentation.

## Import

 ```shell
 # Azure resource can be imported using the resource id, e.g.
 terraform import azapi_resource.example /subscriptions/{subscriptionId}/providers/Microsoft.Network/networkManagerConnections/{resourceName}
 
 # It also supports specifying API version by using the resource id with api-version as a query parameter, e.g.
 terraform import azapi_resource.example /subscriptions/{subscriptionId}/providers/Microsoft.Network/networkManagerConnections/{resourceName}?api-version=2024-10-01
 ```
