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
