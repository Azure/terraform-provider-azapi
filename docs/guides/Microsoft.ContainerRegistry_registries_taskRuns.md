---
subcategory: "Microsoft.ContainerRegistry - Container Registry"
page_title: "registries/taskRuns"
description: |-
  Manages a Container Registry Task Runs.
---

# Microsoft.ContainerRegistry/registries/taskRuns - Container Registry Task Runs

This article demonstrates how to use `azapi` provider to manage the Container Registry Task Runs resource in Azure.

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
      name = "Standard"
      tier = "Standard"
    }
  }
  schema_validation_enabled = false
  response_export_values    = ["*"]
}

resource "azapi_resource" "taskRun" {
  type      = "Microsoft.ContainerRegistry/registries/taskRuns@2019-06-01-preview"
  parent_id = azapi_resource.registry.id
  name      = var.resource_name
  location  = var.location
  body = {
    properties = {
      runRequest = {
        type           = "DockerBuildRequest"
        sourceLocation = "https://github.com/Azure-Samples/aci-helloworld.git#master"
        dockerFilePath = "Dockerfile"
        platform = {
          os = "Linux"
        }
        imageNames = ["helloworld:{{.Run.ID}}", "helloworld:latest"]
      }
    }
  }
  ignore_missing_property = true
}

```



## Arguments Reference

The following arguments are supported:

* `type` - (Required) The type of the resource. This should be set to `Microsoft.ContainerRegistry/registries/taskRuns@api-version`. The available api-versions for this resource are: [`2019-06-01-preview`].

* `parent_id` - (Required) The ID of the azure resource in which this resource is created. The allowed values are:  
  `/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.ContainerRegistry/registries/{resourceName}`

* `name` - (Required) Specifies the name of the azure resource. Changing this forces a new resource to be created.

* `body` - (Required) Specifies the configuration of the resource. More information about the arguments in `body` can be found in the [Microsoft documentation](https://learn.microsoft.com/en-us/azure/templates/Microsoft.ContainerRegistry/registries/taskRuns?pivots=deployment-language-terraform).

For other arguments, please refer to the [azapi_resource](https://registry.terraform.io/providers/Azure/azapi/latest/docs/resources/resource) documentation.

## Import

 ```shell
 # Azure resource can be imported using the resource id, e.g.
 terraform import azapi_resource.example /subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.ContainerRegistry/registries/{resourceName}/taskRuns/{resourceName}
 
 # It also supports specifying API version by using the resource id with api-version as a query parameter, e.g.
 terraform import azapi_resource.example /subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.ContainerRegistry/registries/{resourceName}/taskRuns/{resourceName}?api-version=2019-06-01-preview
 ```
