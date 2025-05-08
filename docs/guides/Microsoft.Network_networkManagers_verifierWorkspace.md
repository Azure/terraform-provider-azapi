---
subcategory: "Microsoft.Network - Various Networking Services"
page_title: "networkManagers/verifierWorkspace"
description: |-
  Manages a Network Managers Verifier Workspace.
---

# Microsoft.Network/networkManagers/verifierWorkspace - Network Managers Verifier Workspace

This article demonstrates how to use `azapi` provider to manage the Network Managers Verifier Workspace resource in Azure.

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

data "azapi_client_config" "current" {}

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
          "/subscriptions/${data.azapi_client_config.current.subscription_id}",
        ]
      }
    }
  }
  retry = {
    error_message_regex = ["CannotDeleteResource"]
  }
}

resource "azapi_resource" "verifierWorkspace" {
  type      = "Microsoft.Network/networkManagers/verifierWorkspaces@2024-01-01-preview"
  parent_id = azapi_resource.networkManager.id
  name      = var.resource_name
  location  = var.location
  body = {
    properties = {
      description = "A sample workspace"
    }
  }

  tags = {
    myTag = "testTag"
  }
}
```



## Arguments Reference

The following arguments are supported:

* `type` - (Required) The type of the resource. This should be set to `Microsoft.Network/networkManagers/verifierWorkspace@api-version`. The available api-versions for this resource are: [].

* `parent_id` - (Required) The ID of the azure resource in which this resource is created. The allowed values are:  
  

* `name` - (Required) Specifies the name of the azure resource. Changing this forces a new resource to be created.

* `body` - (Required) Specifies the configuration of the resource. More information about the arguments in `body` can be found in the [Microsoft documentation](https://learn.microsoft.com/en-us/azure/templates/Microsoft.Network/networkManagers/verifierWorkspace?pivots=deployment-language-terraform).

## Import

 ```shell
 # Azure resource can be imported using the resource id, e.g.
 terraform import azapi_resource.example 
 
 # It also supports specifying API version by using the resource id with api-version as a query parameter, e.g.
 terraform import azapi_resource.example ?api-version=API_VERSION
 ```
