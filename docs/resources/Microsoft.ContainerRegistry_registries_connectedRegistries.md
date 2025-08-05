---
subcategory: "Microsoft.ContainerRegistry - Container Registry"
page_title: "registries/connectedRegistries"
description: |-
  Manages a Container Connected Registry.
---

# Microsoft.ContainerRegistry/registries/connectedRegistries - Container Connected Registry

This article demonstrates how to use `azapi` provider to manage the Container Connected Registry resource in Azure.

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

resource "azapi_resource" "registry" {
  type      = "Microsoft.ContainerRegistry/registries@2023-11-01-preview"
  parent_id = azapi_resource.resourceGroup.id
  name      = "${var.resource_name}registry"
  location  = var.location
  body = {
    properties = {
      adminUserEnabled         = false
      anonymousPullEnabled     = false
      dataEndpointEnabled      = true
      networkRuleBypassOptions = "AzureServices"
      policies = {
        exportPolicy = {
          status = "enabled"
        }
        quarantinePolicy = {
          status = "disabled"
        }
        retentionPolicy = {}
        trustPolicy     = {}
      }
      publicNetworkAccess = "Enabled"
      zoneRedundancy      = "Disabled"
    }
    sku = {
      name = "Premium"
    }
  }
}

resource "azapi_resource" "scopeMap" {
  type      = "Microsoft.ContainerRegistry/registries/scopeMaps@2023-11-01-preview"
  parent_id = azapi_resource.registry.id
  name      = "${var.resource_name}scopemap"
  body = {
    properties = {
      actions     = ["repositories/hello-world/content/delete", "repositories/hello-world/content/read", "repositories/hello-world/content/write", "repositories/hello-world/metadata/read", "repositories/hello-world/metadata/write", "gateway/${var.resource_name}connectedregistry/config/read", "gateway/${var.resource_name}connectedregistry/config/write", "gateway/${var.resource_name}connectedregistry/message/read", "gateway/${var.resource_name}connectedregistry/message/write"]
      description = ""
    }
  }
}

resource "azapi_resource" "token" {
  type      = "Microsoft.ContainerRegistry/registries/tokens@2023-11-01-preview"
  parent_id = azapi_resource.registry.id
  name      = "${var.resource_name}token"
  body = {
    properties = {
      scopeMapId = azapi_resource.scopeMap.id
      status     = "enabled"
    }
  }
}

resource "azapi_resource" "connectedRegistry" {
  type      = "Microsoft.ContainerRegistry/registries/connectedRegistries@2023-11-01-preview"
  parent_id = azapi_resource.registry.id
  name      = "${var.resource_name}connectedregistry"
  body = {
    properties = {
      clientTokenIds = null
      logging = {
        auditLogStatus = "Disabled"
        logLevel       = "None"
      }
      mode = "ReadWrite"
      parent = {
        syncProperties = {
          messageTtl = "P1D"
          schedule   = "* * * * *"
          syncWindow = ""
          tokenId    = azapi_resource.token.id
        }
      }
    }
  }
}
```



## Arguments Reference

The following arguments are supported:

* `type` - (Required) The type of the resource. This should be set to `Microsoft.ContainerRegistry/registries/connectedRegistries@api-version`. The available api-versions for this resource are: [`2020-11-01-preview`, `2021-06-01-preview`, `2021-08-01-preview`, `2021-12-01-preview`, `2022-02-01-preview`, `2023-01-01-preview`, `2023-06-01-preview`, `2023-08-01-preview`, `2023-11-01-preview`, `2024-11-01-preview`, `2025-03-01-preview`, `2025-04-01`, `2025-05-01-preview`].

* `parent_id` - (Required) The ID of the azure resource in which this resource is created. The allowed values are:  
  `/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.ContainerRegistry/registries/{resourceName}`

* `name` - (Required) Specifies the name of the azure resource. Changing this forces a new resource to be created.

* `body` - (Required) Specifies the configuration of the resource. More information about the arguments in `body` can be found in the [Microsoft documentation](https://learn.microsoft.com/en-us/azure/templates/Microsoft.ContainerRegistry/registries/connectedRegistries?pivots=deployment-language-terraform).

For other arguments, please refer to the [azapi_resource](https://registry.terraform.io/providers/Azure/azapi/latest/docs/resources/resource) documentation.

## Import

 ```shell
 # Azure resource can be imported using the resource id, e.g.
 terraform import azapi_resource.example /subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.ContainerRegistry/registries/{resourceName}/connectedRegistries/{resourceName}
 
 # It also supports specifying API version by using the resource id with api-version as a query parameter, e.g.
 terraform import azapi_resource.example /subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.ContainerRegistry/registries/{resourceName}/connectedRegistries/{resourceName}?api-version=2025-05-01-preview
 ```
