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

data "azapi_client_config" "current" {}

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

resource "azapi_resource" "component" {
  type      = "Microsoft.Insights/components@2020-02-02"
  parent_id = azapi_resource.resourceGroup.id
  name      = "ai-${var.resource_name}"
  location  = var.location
  body = {
    kind = "web"
    properties = {
      Application_Type                = "web"
      DisableIpMasking                = false
      DisableLocalAuth                = false
      ForceCustomerStorageForProfiler = false
      RetentionInDays                 = 90
      SamplingPercentage              = 100
      publicNetworkAccessForIngestion = "Enabled"
      publicNetworkAccessForQuery     = "Enabled"
    }
  }
}

resource "azapi_resource" "vault" {
  type      = "Microsoft.KeyVault/vaults@2023-02-01"
  parent_id = azapi_resource.resourceGroup.id
  name      = "kv${var.resource_name}"
  location  = var.location
  body = {
    properties = {
      accessPolicies               = []
      createMode                   = "default"
      enableRbacAuthorization      = false
      enabledForDeployment         = false
      enabledForDiskEncryption     = false
      enabledForTemplateDeployment = false
      publicNetworkAccess          = "Enabled"
      sku = {
        family = "A"
        name   = "standard"
      }
      softDeleteRetentionInDays = 7
      tenantId                  = data.azapi_client_config.current.tenant_id
    }
  }
}

resource "azapi_resource" "managedCluster" {
  type      = "Microsoft.ContainerService/managedClusters@2025-02-01"
  parent_id = azapi_resource.resourceGroup.id
  name      = "aks-${var.resource_name}"
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
      dnsPrefix            = "aks-${var.resource_name}"
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

resource "azapi_resource" "storageAccount" {
  type      = "Microsoft.Storage/storageAccounts@2023-05-01"
  parent_id = azapi_resource.resourceGroup.id
  name      = "st${var.resource_name}"
  location  = var.location
  body = {
    kind = "StorageV2"
    properties = {
      accessTier                   = "Hot"
      allowBlobPublicAccess        = true
      allowCrossTenantReplication  = false
      allowSharedKeyAccess         = true
      defaultToOAuthAuthentication = false
      dnsEndpointType              = "Standard"
      encryption = {
        keySource = "Microsoft.Storage"
        services = {
          queue = {
            keyType = "Service"
          }
          table = {
            keyType = "Service"
          }
        }
      }
      isHnsEnabled       = false
      isLocalUserEnabled = true
      isNfsV3Enabled     = false
      isSftpEnabled      = false
      minimumTlsVersion  = "TLS1_2"
      networkAcls = {
        bypass              = "AzureServices"
        defaultAction       = "Allow"
        ipRules             = []
        resourceAccessRules = []
        virtualNetworkRules = []
      }
      publicNetworkAccess      = "Enabled"
      supportsHttpsTrafficOnly = true
    }
    sku = {
      name = "Standard_LRS"
    }
  }
}

resource "azapi_resource" "workspace" {
  type      = "Microsoft.MachineLearningServices/workspaces@2024-04-01"
  parent_id = azapi_resource.resourceGroup.id
  name      = "mlw-${var.resource_name}"
  location  = var.location
  identity {
    type         = "SystemAssigned"
    identity_ids = []
  }
  body = {
    kind = "Default"
    properties = {
      applicationInsights = azapi_resource.component.id
      keyVault            = azapi_resource.vault.id
      publicNetworkAccess = "Enabled"
      storageAccount      = azapi_resource.storageAccount.id
      v1LegacyMode        = false
    }
    sku = {
      name = "Basic"
      tier = "Basic"
    }
  }
}

resource "azapi_resource" "trustedAccessRoleBinding" {
  type      = "Microsoft.ContainerService/managedClusters/trustedAccessRoleBindings@2025-02-01"
  parent_id = azapi_resource.managedCluster.id
  name      = "tarb-${var.resource_name}"
  body = {
    properties = {
      roles            = ["Microsoft.MachineLearningServices/workspaces/mlworkload"]
      sourceResourceId = azapi_resource.workspace.id
    }
  }
}
