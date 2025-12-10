---
subcategory: "Microsoft.App - Azure Container Apps"
page_title: "jobs"
description: |-
  Manages a Container App Job.
---

# Microsoft.App/jobs - Container App Job

This article demonstrates how to use `azapi` provider to manage the Container App Job resource in Azure.



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

resource "azapi_resource" "workspace" {
  type      = "Microsoft.OperationalInsights/workspaces@2023-09-01"
  parent_id = azapi_resource.resourceGroup.id
  name      = "${var.resource_name}-law"
  location  = var.location
  body = {
    properties = {
      sku = {
        name = "PerGB2018"
      }
      retentionInDays = 30
    }
  }
}

data "azapi_resource_action" "workspace_keys" {
  type                   = "Microsoft.OperationalInsights/workspaces@2023-09-01"
  resource_id            = azapi_resource.workspace.id
  action                 = "listKeys"
  method                 = "POST"
  response_export_values = ["*"]
}

resource "azapi_resource" "managedEnvironment" {
  type      = "Microsoft.App/managedEnvironments@2025-01-01"
  parent_id = azapi_resource.resourceGroup.id
  name      = "${var.resource_name}-env"
  location  = var.location
  body = {
    properties = {
      appLogsConfiguration = {
        destination = "log-analytics"
        logAnalyticsConfiguration = {
          customerId = azapi_resource.workspace.output.properties.customerId
          sharedKey  = data.azapi_resource_action.workspace_keys.output.primarySharedKey
        }
      }
    }
  }
}

resource "azapi_resource" "job" {
  type      = "Microsoft.App/jobs@2025-01-01"
  parent_id = azapi_resource.resourceGroup.id
  name      = "${var.resource_name}-cajob"
  location  = var.location
  body = {
    properties = {
      configuration = {
        manualTriggerConfig = {
          parallelism            = 4
          replicaCompletionCount = 1
        }
        replicaRetryLimit = 10
        replicaTimeout    = 10
        triggerType       = "Manual"
      }
      environmentId = azapi_resource.managedEnvironment.id
      template = {
        containers = [{
          env    = []
          image  = "jackofallops/azure-containerapps-python-acctest:v0.0.1"
          name   = "testcontainerappsjob0"
          probes = []
          resources = {
            cpu    = 0.5
            memory = "1Gi"
          }
          volumeMounts = []
        }]
        initContainers = []
        volumes        = []
      }
    }
  }
}

```



## Arguments Reference

The following arguments are supported:

* `type` - (Required) The type of the resource. This should be set to `Microsoft.App/jobs@api-version`. The available api-versions for this resource are: [`2022-11-01-preview`, `2023-04-01-preview`, `2023-05-01`, `2023-05-02-preview`, `2023-08-01-preview`, `2023-11-02-preview`, `2024-02-02-preview`, `2024-03-01`, `2024-08-02-preview`, `2024-10-02-preview`, `2025-01-01`, `2025-02-02-preview`, `2025-07-01`, `2025-10-02-preview`].

* `parent_id` - (Required) The ID of the azure resource in which this resource is created. The allowed values are:  
  `/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}`

* `name` - (Required) Specifies the name of the azure resource. Changing this forces a new resource to be created.

* `body` - (Required) Specifies the configuration of the resource. More information about the arguments in `body` can be found in the [Microsoft documentation](https://learn.microsoft.com/en-us/azure/templates/Microsoft.App/jobs?pivots=deployment-language-terraform).

For other arguments, please refer to the [azapi_resource](https://registry.terraform.io/providers/Azure/azapi/latest/docs/resources/resource) documentation.

## Import

 ```shell
 # Azure resource can be imported using the resource id, e.g.
 terraform import azapi_resource.example /subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.App/jobs/{resourceName}
 
 # It also supports specifying API version by using the resource id with api-version as a query parameter, e.g.
 terraform import azapi_resource.example /subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.App/jobs/{resourceName}?api-version=2025-10-02-preview
 ```
