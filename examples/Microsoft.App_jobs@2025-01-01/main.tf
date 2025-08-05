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
