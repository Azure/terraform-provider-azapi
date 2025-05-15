---
subcategory: "Microsoft.KubernetesConfiguration - Azure Arc-enabled Kubernetes"
page_title: "fluxConfigurations"
description: |-
  Manages a Kubernetes Flux Configuration.
---

# Microsoft.KubernetesConfiguration/fluxConfigurations - Kubernetes Flux Configuration

This article demonstrates how to use `azapi` provider to manage the Kubernetes Flux Configuration resource in Azure.

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

resource "azapi_resource" "managedCluster" {
  type      = "Microsoft.ContainerService/managedClusters@2023-04-02-preview"
  parent_id = azapi_resource.resourceGroup.id
  name      = var.resource_name
  location  = var.location
  identity {
    type         = "SystemAssigned"
    identity_ids = []
  }
  body = {
    properties = {
      agentPoolProfiles = [
        {
          count  = 1
          mode   = "System"
          name   = "default"
          vmSize = "Standard_DS2_v2"
        },
      ]
      dnsPrefix = var.resource_name
    }
  }
  schema_validation_enabled = false
  response_export_values    = ["*"]
}

resource "azapi_resource" "extension" {
  type      = "Microsoft.KubernetesConfiguration/extensions@2022-11-01"
  parent_id = azapi_resource.managedCluster.id
  name      = var.resource_name
  body = {
    properties = {
      autoUpgradeMinorVersion = true
      extensionType           = "microsoft.flux"
    }
  }
  schema_validation_enabled = false
  response_export_values    = ["*"]
}

resource "azapi_resource" "fluxConfiguration" {
  type      = "Microsoft.KubernetesConfiguration/fluxConfigurations@2022-03-01"
  name      = var.resource_name
  parent_id = azapi_resource.managedCluster.id

  body = {
    properties = {
      scope      = "cluster"
      namespace  = "flux-system"
      sourceKind = "GitRepository"
      suspend    = false
      gitRepository = {
        url                   = "https://github.com/Azure/arc-k8s-demo"
        timeoutInSeconds      = 120
        syncIntervalInSeconds = 120
        repositoryRef = {
          branch = "branch"
        }
      }
      kustomizations = {
        shared = {
          path                   = "cluster-config/shared"
          timeoutInSeconds       = 600
          syncIntervalInSeconds  = 60
          retryIntervalInSeconds = 60
          prune                  = false
          force                  = false
        }
        applications = {
          path                   = "cluster-config/applications"
          dependsOn              = ["shared"]
          timeoutInSeconds       = 600
          syncIntervalInSeconds  = 60
          retryIntervalInSeconds = 60
          prune                  = false
          force                  = false
        }
      }
    }
  }

  depends_on = [
    azapi_resource.extension
  ]
}

```



## Arguments Reference

The following arguments are supported:

* `type` - (Required) The type of the resource. This should be set to `Microsoft.KubernetesConfiguration/fluxConfigurations@api-version`. The available api-versions for this resource are: [`2021-11-01-preview`, `2022-01-01-preview`, `2022-03-01`, `2022-07-01`, `2022-11-01`, `2023-05-01`, `2024-04-01-preview`, `2024-11-01`, `2025-04-01`].

* `parent_id` - (Required) The ID of the azure resource in which this resource is created. The allowed values are:  
  `{any azure resource id}`

* `name` - (Required) Specifies the name of the azure resource. Changing this forces a new resource to be created.

* `body` - (Required) Specifies the configuration of the resource. More information about the arguments in `body` can be found in the [Microsoft documentation](https://learn.microsoft.com/en-us/azure/templates/Microsoft.KubernetesConfiguration/fluxConfigurations?pivots=deployment-language-terraform).

For other arguments, please refer to the [azapi_resource](https://registry.terraform.io/providers/Azure/azapi/latest/docs/resources/resource) documentation.

## Import

 ```shell
 # Azure resource can be imported using the resource id, e.g.
 terraform import azapi_resource.example {any azure resource id}/providers/Microsoft.KubernetesConfiguration/fluxConfigurations/{resourceName}
 
 # It also supports specifying API version by using the resource id with api-version as a query parameter, e.g.
 terraform import azapi_resource.example {any azure resource id}/providers/Microsoft.KubernetesConfiguration/fluxConfigurations/{resourceName}?api-version=2025-04-01
 ```
