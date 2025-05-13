---
subcategory: "Microsoft.ContainerRegistry - Container Registry"
page_title: "registries/scopeMaps"
description: |-
  Manages a Azure Container Registry scope map.
---

# Microsoft.ContainerRegistry/registries/scopeMaps - Azure Container Registry scope map

This article demonstrates how to use `azapi` provider to manage the Azure Container Registry scope map resource in Azure.

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

resource "azapi_resource" "registry" {
  type      = "Microsoft.ContainerRegistry/registries@2021-08-01-preview"
  parent_id = azapi_resource.resourceGroup.id
  name      = var.resource_name
  location  = var.location
  body = {
    properties = {
      adminUserEnabled     = false
      anonymousPullEnabled = false
      dataEndpointEnabled  = false
      encryption = {
        status = "disabled"
      }
      networkRuleBypassOptions = "AzureServices"
      policies = {
        exportPolicy = {
          status = "enabled"
        }
        quarantinePolicy = {
          status = "disabled"
        }
        retentionPolicy = {
          status = "disabled"
        }
        trustPolicy = {
          status = "disabled"
        }
      }
      publicNetworkAccess = "Enabled"
      zoneRedundancy      = "Disabled"
    }
    sku = {
      name = "Premium"
      tier = "Premium"
    }
  }
  schema_validation_enabled = false
  response_export_values    = ["*"]
}

resource "azapi_resource" "scopeMap" {
  type      = "Microsoft.ContainerRegistry/registries/scopeMaps@2021-08-01-preview"
  parent_id = azapi_resource.registry.id
  name      = var.resource_name
  body = {
    properties = {
      actions = [
        "repositories/testrepo/content/read",
      ]
      description = ""
    }
  }
  schema_validation_enabled = false
  response_export_values    = ["*"]
}


```



## Arguments Reference

The following arguments are supported:

* `type` - (Required) The type of the resource. This should be set to `Microsoft.ContainerRegistry/registries/scopeMaps@api-version`. The available api-versions for this resource are: [`2019-05-01-preview`, `2020-11-01-preview`, `2021-06-01-preview`, `2021-08-01-preview`, `2021-12-01-preview`, `2022-02-01-preview`, `2022-12-01`, `2023-01-01-preview`, `2023-06-01-preview`, `2023-07-01`, `2023-08-01-preview`, `2023-11-01-preview`, `2024-11-01-preview`].

* `parent_id` - (Required) The ID of the azure resource in which this resource is created. The allowed values are:  
  `/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.ContainerRegistry/registries/{resourceName}`

* `name` - (Required) Specifies the name of the azure resource. Changing this forces a new resource to be created.

* `body` - (Required) Specifies the configuration of the resource. More information about the arguments in `body` can be found in the [Microsoft documentation](https://learn.microsoft.com/en-us/azure/templates/Microsoft.ContainerRegistry/registries/scopeMaps?pivots=deployment-language-terraform).

For other arguments, please refer to the [azapi_resource](https://registry.terraform.io/providers/Azure/azapi/latest/docs/resources/resource) documentation.

## Import

 ```shell
 # Azure resource can be imported using the resource id, e.g.
 terraform import azapi_resource.example /subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.ContainerRegistry/registries/{resourceName}/scopeMaps/{resourceName}
 
 # It also supports specifying API version by using the resource id with api-version as a query parameter, e.g.
 terraform import azapi_resource.example /subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.ContainerRegistry/registries/{resourceName}/scopeMaps/{resourceName}?api-version=2024-11-01-preview
 ```
