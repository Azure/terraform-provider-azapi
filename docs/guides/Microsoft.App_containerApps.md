---
subcategory: "Microsoft.App - Azure Container Apps"
page_title: "containerApps"
description: |-
  Manages a Container App.
---

# Microsoft.App/containerApps - Container App

This article demonstrates how to use `azapi` provider to manage the Container App resource in Azure.

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

resource "azapi_resource" "workspace" {
  type      = "Microsoft.OperationalInsights/workspaces@2022-10-01"
  parent_id = azapi_resource.resourceGroup.id
  name      = var.resource_name
  location  = var.location
  body = {
    properties = {
      features = {
        disableLocalAuth                            = false
        enableLogAccessUsingOnlyResourcePermissions = true
      }
      publicNetworkAccessForIngestion = "Enabled"
      publicNetworkAccessForQuery     = "Enabled"
      retentionInDays                 = 30
      sku = {
        name = "PerGB2018"
      }
      workspaceCapping = {
        dailyQuotaGb = -1
      }
    }
  }
  schema_validation_enabled = false
  response_export_values    = ["*"]
}

data "azapi_resource_action" "sharedKeys" {
  type                   = "Microsoft.OperationalInsights/workspaces@2020-08-01"
  resource_id            = azapi_resource.workspace.id
  action                 = "sharedKeys"
  response_export_values = ["*"]
}

resource "azapi_resource" "managedEnvironment" {
  type      = "Microsoft.App/managedEnvironments@2022-03-01"
  parent_id = azapi_resource.resourceGroup.id
  name      = var.resource_name
  location  = var.location
  body = {
    properties = {
      appLogsConfiguration = {
        destination = "log-analytics"
        logAnalyticsConfiguration = {
          customerId = azapi_resource.workspace.output.properties.customerId
          sharedKey  = data.azapi_resource_action.sharedKeys.output.primarySharedKey
        }
      }
      vnetConfiguration = {
      }
    }
  }
  schema_validation_enabled = false
  response_export_values    = ["*"]
}

resource "azapi_resource" "containerApp" {
  type      = "Microsoft.App/containerApps@2022-03-01"
  parent_id = azapi_resource.resourceGroup.id
  name      = var.resource_name
  location  = var.location
  body = {
    properties = {
      configuration = {
        activeRevisionsMode = "Single"
      }
      managedEnvironmentId = azapi_resource.managedEnvironment.id
      template = {
        containers = [
          {
            env = [
            ]
            image = "jackofallops/azure-containerapps-python-acctest:v0.0.1"
            name  = "acctest-cont-230630032906865620"
            probes = [
            ]
            resources = {
              cpu              = 0.25
              ephemeralStorage = "1Gi"
              memory           = "0.5Gi"
            }
            volumeMounts = [
            ]
          },
        ]
        scale = {
          maxReplicas = 10
        }
        volumes = [
        ]
      }
    }
  }
  schema_validation_enabled = false
  response_export_values    = ["*"]
}


```



## Arguments Reference

The following arguments are supported:

* `type` - (Required) The type of the resource. This should be set to `Microsoft.App/containerApps@api-version`. The available api-versions for this resource are: [`2022-01-01-preview`, `2022-03-01`, `2022-06-01-preview`, `2022-10-01`, `2022-11-01-preview`, `2023-04-01-preview`, `2023-05-01`, `2023-05-02-preview`, `2023-08-01-preview`, `2023-11-02-preview`, `2024-02-02-preview`, `2024-03-01`, `2024-08-02-preview`, `2024-10-02-preview`].

* `parent_id` - (Required) The ID of the azure resource in which this resource is created. The allowed values are:  
  `/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}`

* `name` - (Required) Specifies the name of the azure resource. Changing this forces a new resource to be created.

* `body` - (Required) Specifies the configuration of the resource. More information about the arguments in `body` can be found in the [Microsoft documentation](https://learn.microsoft.com/en-us/azure/templates/Microsoft.App/containerApps?pivots=deployment-language-terraform).

## Import

 ```shell
 # Azure resource can be imported using the resource id, e.g.
 terraform import azapi_resource.example /subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.App/containerApps/{resourceName}
 
 # It also supports specifying API version by using the resource id with api-version as a query parameter, e.g.
 terraform import azapi_resource.example /subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.App/containerApps/{resourceName}?api-version=2024-10-02-preview
 ```
