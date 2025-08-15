---
subcategory: "Microsoft.ContainerService - Azure Kubernetes Service (AKS)"
page_title: "managedClusters/trustedAccessRoleBindings"
description: |-
  Manages a Kubernetes Cluster Trusted Access Role Binding.
---

# Microsoft.ContainerService/managedClusters/trustedAccessRoleBindings - Kubernetes Cluster Trusted Access Role Binding

This article demonstrates how to use `azapi` provider to manage the Kubernetes Cluster Trusted Access Role Binding resource in Azure.

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

```



## Arguments Reference

The following arguments are supported:

* `type` - (Required) The type of the resource. This should be set to `Microsoft.ContainerService/managedClusters/trustedAccessRoleBindings@api-version`. The available api-versions for this resource are: [`2022-04-02-preview`, `2022-05-02-preview`, `2022-06-02-preview`, `2022-07-02-preview`, `2022-08-02-preview`, `2022-08-03-preview`, `2022-09-02-preview`, `2022-10-02-preview`, `2022-11-02-preview`, `2023-01-02-preview`, `2023-02-02-preview`, `2023-03-02-preview`, `2023-04-02-preview`, `2023-05-02-preview`, `2023-06-02-preview`, `2023-07-02-preview`, `2023-08-02-preview`, `2023-09-01`, `2023-09-02-preview`, `2023-10-01`, `2023-10-02-preview`, `2023-11-01`, `2023-11-02-preview`, `2024-01-01`, `2024-01-02-preview`, `2024-02-01`, `2024-02-02-preview`, `2024-03-02-preview`, `2024-04-02-preview`, `2024-05-01`, `2024-05-02-preview`, `2024-06-02-preview`, `2024-07-01`, `2024-07-02-preview`, `2024-08-01`, `2024-09-01`, `2024-09-02-preview`, `2024-10-01`, `2024-10-02-preview`, `2025-01-01`, `2025-01-02-preview`, `2025-02-01`, `2025-02-02-preview`, `2025-03-01`, `2025-03-02-preview`, `2025-04-01`, `2025-04-02-preview`, `2025-05-01`, `2025-05-02-preview`, `2025-06-02-preview`].

* `parent_id` - (Required) The ID of the azure resource in which this resource is created. The allowed values are:  
  `/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.ContainerService/managedClusters/{resourceName}`

* `name` - (Required) Specifies the name of the azure resource. Changing this forces a new resource to be created.

* `body` - (Required) Specifies the configuration of the resource. More information about the arguments in `body` can be found in the [Microsoft documentation](https://learn.microsoft.com/en-us/azure/templates/Microsoft.ContainerService/managedClusters/trustedAccessRoleBindings?pivots=deployment-language-terraform).

For other arguments, please refer to the [azapi_resource](https://registry.terraform.io/providers/Azure/azapi/latest/docs/resources/resource) documentation.

## Import

 ```shell
 # Azure resource can be imported using the resource id, e.g.
 terraform import azapi_resource.example /subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.ContainerService/managedClusters/{resourceName}/trustedAccessRoleBindings/{resourceName}
 
 # It also supports specifying API version by using the resource id with api-version as a query parameter, e.g.
 terraform import azapi_resource.example /subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.ContainerService/managedClusters/{resourceName}/trustedAccessRoleBindings/{resourceName}?api-version=2025-06-02-preview
 ```
