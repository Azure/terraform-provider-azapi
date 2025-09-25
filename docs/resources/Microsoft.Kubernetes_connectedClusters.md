---
subcategory: "Microsoft.Kubernetes - Azure Arc-enabled Kubernetes"
page_title: "connectedClusters"
description: |-
  Manages a Arc Kubernetes Cluster.
---

# Microsoft.Kubernetes/connectedClusters - Arc Kubernetes Cluster

This article demonstrates how to use `azapi` provider to manage the Arc Kubernetes Cluster resource in Azure.

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

resource "azapi_resource" "connectedCluster" {
  type      = "Microsoft.Kubernetes/connectedClusters@2024-01-01"
  parent_id = azapi_resource.resourceGroup.id
  name      = "${var.resource_name}-cc"
  location  = var.location
  identity {
    type         = "SystemAssigned"
    identity_ids = []
  }
  body = {
    kind = "ProvisionedCluster"
    properties = {
      agentPublicKeyCertificate = ""
      arcAgentProfile = {
        agentAutoUpgrade = "Enabled"
      }
    }
  }
}


```



## Arguments Reference

The following arguments are supported:

* `type` - (Required) The type of the resource. This should be set to `Microsoft.Kubernetes/connectedClusters@api-version`. The available api-versions for this resource are: [`2020-01-01-preview`, `2021-03-01`, `2021-04-01-preview`, `2021-10-01`, `2022-05-01-preview`, `2022-10-01-preview`, `2023-11-01-preview`, `2024-01-01`, `2024-02-01-preview`, `2024-06-01-preview`, `2024-07-01-preview`, `2024-07-15-preview`, `2024-12-01-preview`, `2025-08-01-preview`, `2025-12-01-preview`].

* `parent_id` - (Required) The ID of the azure resource in which this resource is created. The allowed values are:  
  `/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}`

* `name` - (Required) Specifies the name of the azure resource. Changing this forces a new resource to be created.

* `body` - (Required) Specifies the configuration of the resource. More information about the arguments in `body` can be found in the [Microsoft documentation](https://learn.microsoft.com/en-us/azure/templates/Microsoft.Kubernetes/connectedClusters?pivots=deployment-language-terraform).

For other arguments, please refer to the [azapi_resource](https://registry.terraform.io/providers/Azure/azapi/latest/docs/resources/resource) documentation.

## Import

 ```shell
 # Azure resource can be imported using the resource id, e.g.
 terraform import azapi_resource.example /subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Kubernetes/connectedClusters/{resourceName}
 
 # It also supports specifying API version by using the resource id with api-version as a query parameter, e.g.
 terraform import azapi_resource.example /subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Kubernetes/connectedClusters/{resourceName}?api-version=2025-12-01-preview
 ```
