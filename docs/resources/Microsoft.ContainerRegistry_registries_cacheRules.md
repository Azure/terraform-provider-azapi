---
subcategory: "Microsoft.ContainerRegistry - Container Registry"
page_title: "registries/cacheRules"
description: |-
  Manages a Azure Container Registry Cache Rule.
---

# Microsoft.ContainerRegistry/registries/cacheRules - Azure Container Registry Cache Rule

This article demonstrates how to use `azapi` provider to manage the Azure Container Registry Cache Rule resource in Azure.

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
      dataEndpointEnabled      = false
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
      name = "Basic"
    }
  }
}

resource "azapi_resource" "cacheRule" {
  type      = "Microsoft.ContainerRegistry/registries/cacheRules@2023-07-01"
  parent_id = azapi_resource.registry.id
  name      = "${var.resource_name}-cache-rule"
  body = {
    properties = {
      sourceRepository = "mcr.microsoft.com/hello-world"
      targetRepository = "target"
    }
  }
}

```



## Arguments Reference

The following arguments are supported:

* `type` - (Required) The type of the resource. This should be set to `Microsoft.ContainerRegistry/registries/cacheRules@api-version`. The available api-versions for this resource are: [`2023-01-01-preview`, `2023-06-01-preview`, `2023-07-01`, `2023-08-01-preview`, `2023-11-01-preview`, `2024-11-01-preview`, `2025-03-01-preview`, `2025-04-01`, `2025-05-01-preview`].

* `parent_id` - (Required) The ID of the azure resource in which this resource is created. The allowed values are:  
  `/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.ContainerRegistry/registries/{resourceName}`

* `name` - (Required) Specifies the name of the azure resource. Changing this forces a new resource to be created.

* `body` - (Required) Specifies the configuration of the resource. More information about the arguments in `body` can be found in the [Microsoft documentation](https://learn.microsoft.com/en-us/azure/templates/Microsoft.ContainerRegistry/registries/cacheRules?pivots=deployment-language-terraform).

For other arguments, please refer to the [azapi_resource](https://registry.terraform.io/providers/Azure/azapi/latest/docs/resources/resource) documentation.

## Import

 ```shell
 # Azure resource can be imported using the resource id, e.g.
 terraform import azapi_resource.example /subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.ContainerRegistry/registries/{resourceName}/cacheRules/{resourceName}
 
 # It also supports specifying API version by using the resource id with api-version as a query parameter, e.g.
 terraform import azapi_resource.example /subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.ContainerRegistry/registries/{resourceName}/cacheRules/{resourceName}?api-version=2025-05-01-preview
 ```
