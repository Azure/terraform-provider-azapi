---
subcategory: "Microsoft.ContainerService - Azure Kubernetes Service (AKS)"
page_title: "managedClusters"
description: |-
  Manages a managed Kubernetes Cluster (also known as AKS / Azure Kubernetes Service).
---

# Microsoft.ContainerService/managedClusters - managed Kubernetes Cluster (also known as AKS / Azure Kubernetes Service)

This article demonstrates how to use `azapi` provider to manage the managed Kubernetes Cluster (also known as AKS / Azure Kubernetes Service) resource in Azure.



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
  type                      = "Microsoft.Resources/resourceGroups@2020-06-01"
  name                      = var.resource_name
  location                  = var.location
  schema_validation_enabled = false
  response_export_values    = ["*"]
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


```



## Arguments Reference

The following arguments are supported:

* `type` - (Required) The type of the resource. This should be set to `Microsoft.ContainerService/managedClusters@api-version`. The available api-versions for this resource are: [`2017-08-31`, `2018-03-31`, `2018-08-01-preview`, `2019-02-01`, `2019-04-01`, `2019-06-01`, `2019-08-01`, `2019-10-01`, `2019-11-01`, `2020-01-01`, `2020-02-01`, `2020-03-01`, `2020-04-01`, `2020-06-01`, `2020-07-01`, `2020-09-01`, `2020-11-01`, `2020-12-01`, `2021-02-01`, `2021-03-01`, `2021-05-01`, `2021-07-01`, `2021-08-01`, `2021-09-01`, `2021-10-01`, `2021-11-01-preview`, `2022-01-01`, `2022-01-02-preview`, `2022-02-01`, `2022-02-02-preview`, `2022-03-01`, `2022-03-02-preview`, `2022-04-01`, `2022-04-02-preview`, `2022-05-02-preview`, `2022-06-01`, `2022-06-02-preview`, `2022-07-01`, `2022-07-02-preview`, `2022-08-02-preview`, `2022-08-03-preview`, `2022-09-01`, `2022-09-02-preview`, `2022-10-02-preview`, `2022-11-01`, `2022-11-02-preview`, `2023-01-01`, `2023-01-02-preview`, `2023-02-01`, `2023-02-02-preview`, `2023-03-01`, `2023-03-02-preview`, `2023-04-01`, `2023-04-02-preview`, `2023-05-01`, `2023-05-02-preview`, `2023-06-01`, `2023-06-02-preview`, `2023-07-01`, `2023-07-02-preview`, `2023-08-01`, `2023-08-02-preview`, `2023-09-01`, `2023-09-02-preview`, `2023-10-01`, `2023-10-02-preview`, `2023-11-01`, `2023-11-02-preview`, `2024-01-01`, `2024-01-02-preview`, `2024-02-01`, `2024-02-02-preview`, `2024-03-02-preview`, `2024-04-02-preview`, `2024-05-01`, `2024-05-02-preview`, `2024-06-02-preview`, `2024-07-01`, `2024-07-02-preview`, `2024-08-01`, `2024-09-01`, `2024-09-02-preview`, `2024-10-01`, `2024-10-02-preview`, `2025-01-01`, `2025-01-02-preview`, `2025-02-01`, `2025-02-02-preview`, `2025-03-01`, `2025-03-02-preview`, `2025-04-01`, `2025-04-02-preview`, `2025-05-01`, `2025-05-02-preview`, `2025-06-02-preview`, `2025-07-01`, `2025-07-02-preview`, `2025-08-01`, `2025-08-02-preview`, `2025-09-01`, `2025-09-02-preview`, `2025-10-01`].

* `parent_id` - (Required) The ID of the azure resource in which this resource is created. The allowed values are:  
  `/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}`

* `name` - (Required) Specifies the name of the azure resource. Changing this forces a new resource to be created.

* `body` - (Required) Specifies the configuration of the resource. More information about the arguments in `body` can be found in the [Microsoft documentation](https://learn.microsoft.com/en-us/azure/templates/Microsoft.ContainerService/managedClusters?pivots=deployment-language-terraform).

For other arguments, please refer to the [azapi_resource](https://registry.terraform.io/providers/Azure/azapi/latest/docs/resources/resource) documentation.

## Import

 ```shell
 # Azure resource can be imported using the resource id, e.g.
 terraform import azapi_resource.example /subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.ContainerService/managedClusters/{resourceName}
 
 # It also supports specifying API version by using the resource id with api-version as a query parameter, e.g.
 terraform import azapi_resource.example /subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.ContainerService/managedClusters/{resourceName}?api-version=2025-10-01
 ```
