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

resource "azapi_resource" "managedCluster" {
  type      = "Microsoft.ContainerService/managedClusters@2025-02-01"
  parent_id = azapi_resource.resourceGroup.id
  name      = var.resource_name
  location  = var.location
  identity {
    type         = "SystemAssigned"
    identity_ids = []
  }
  body = {
    properties = {
      addonProfiles = {}
      agentPoolProfiles = [{
        count                  = 1
        enableAutoScaling      = false
        enableEncryptionAtHost = false
        enableFIPS             = false
        enableNodePublicIP     = false
        enableUltraSSD         = false
        kubeletDiskType        = ""
        mode                   = "System"
        name                   = "default"
        nodeLabels             = {}
        osDiskType             = "Managed"
        osType                 = "Linux"
        scaleDownMode          = "Delete"
        tags                   = {}
        type                   = "VirtualMachineScaleSets"
        upgradeSettings = {
          drainTimeoutInMinutes     = 0
          maxSurge                  = "10%"
          nodeSoakDurationInMinutes = 0
        }
        vmSize = "Standard_B2s"
      }]
      apiServerAccessProfile = {
        disableRunCommand              = false
        enablePrivateCluster           = false
        enablePrivateClusterPublicFQDN = false
      }
      autoUpgradeProfile = {
        nodeOSUpgradeChannel = "NodeImage"
        upgradeChannel       = "none"
      }
      azureMonitorProfile = {
        metrics = {
          enabled = false
        }
      }
      disableLocalAccounts = false
      dnsPrefix            = var.resource_name
      enableRBAC           = true
      kubernetesVersion    = ""
      metricsProfile = {
        costAnalysis = {
          enabled = false
        }
      }
      nodeResourceGroup = ""
      securityProfile   = {}
      servicePrincipalProfile = {
        clientId = "msi"
      }
      supportPlan = "KubernetesOfficial"
    }
    sku = {
      name = "Base"
      tier = "Free"
    }
  }
}

resource "azapi_resource" "fleet" {
  type      = "Microsoft.ContainerService/fleets@2024-04-01"
  parent_id = azapi_resource.resourceGroup.id
  name      = var.resource_name
  location  = var.location
  body = {
    properties = {}
  }
}

resource "azapi_resource" "member" {
  type      = "Microsoft.ContainerService/fleets/members@2024-04-01"
  parent_id = azapi_resource.fleet.id
  name      = var.resource_name
  body = {
    properties = {
      clusterResourceId = azapi_resource.managedCluster.id
      group             = "default"
    }
  }
}
