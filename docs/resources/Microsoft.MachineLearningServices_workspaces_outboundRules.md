---
subcategory: "Microsoft.MachineLearningServices - Azure Machine Learning"
page_title: "workspaces/outboundRules"
description: |-
  Manages a Azure Machine Learning Workspace FQDN Network Outbound Rule.
---

# Microsoft.MachineLearningServices/workspaces/outboundRules - Azure Machine Learning Workspace FQDN Network Outbound Rule

This article demonstrates how to use `azapi` provider to manage the Azure Machine Learning Workspace FQDN Network Outbound Rule resource in Azure.

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

data "azapi_client_config" "current" {}

locals {
  base_name      = lower(var.resource_name)
  sa_base        = replace(local.base_name, "-", "")
  kv_base        = replace(local.base_name, "-", "")
  storage_name   = substr("sa${local.sa_base}", 0, 24)
  key_vault_name = substr("kv${local.kv_base}", 0, 24)
  ai_name        = "${var.resource_name}-ai"
  workspace_name = "${var.resource_name}-mlw"
  outbound_name  = "${var.resource_name}-outbound"
}

resource "azapi_resource" "resourceGroup" {
  type     = "Microsoft.Resources/resourceGroups@2020-06-01"
  name     = var.resource_name
  location = var.location
}

resource "azapi_resource" "component" {
  type      = "Microsoft.Insights/components@2020-02-02"
  parent_id = azapi_resource.resourceGroup.id
  name      = local.ai_name
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
  name      = local.key_vault_name
  location  = var.location
  body = {
    properties = {
      accessPolicies               = []
      createMode                   = "default"
      enablePurgeProtection        = true
      enableRbacAuthorization      = false
      enableSoftDelete             = true
      enabledForDeployment         = false
      enabledForDiskEncryption     = false
      enabledForTemplateDeployment = false
      publicNetworkAccess          = "Enabled"
      sku = {
        family = "A"
        name   = "standard"
      }
      tenantId = data.azapi_client_config.current.tenant_id
    }
  }
}

resource "azapi_resource" "storageAccount" {
  type      = "Microsoft.Storage/storageAccounts@2023-05-01"
  parent_id = azapi_resource.resourceGroup.id
  name      = local.storage_name
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
  name      = local.workspace_name
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
      managedNetwork = {
        isolationMode = "AllowOnlyApprovedOutbound"
      }
      publicNetworkAccess = "Enabled"
      storageAccount      = azapi_resource.storageAccount.id
      v1LegacyMode        = false
    }
    sku = {
      name = "Basic"
    }
  }
}

resource "azapi_resource" "outboundRule" {
  type      = "Microsoft.MachineLearningServices/workspaces/outboundRules@2024-04-01"
  parent_id = azapi_resource.workspace.id
  name      = local.outbound_name
  body = {
    properties = {
      category    = "UserDefined"
      destination = "www.microsoft.com"
      type        = "FQDN"
    }
  }
}

```



## Arguments Reference

The following arguments are supported:

* `type` - (Required) The type of the resource. This should be set to `Microsoft.MachineLearningServices/workspaces/outboundRules@api-version`. The available api-versions for this resource are: [`2023-04-01-preview`, `2023-06-01-preview`, `2023-08-01-preview`, `2023-10-01`, `2024-01-01-preview`, `2024-04-01`, `2024-04-01-preview`, `2024-07-01-preview`, `2024-10-01`, `2024-10-01-preview`, `2025-01-01-preview`, `2025-04-01`, `2025-04-01-preview`, `2025-06-01`, `2025-07-01-preview`].

* `parent_id` - (Required) The ID of the azure resource in which this resource is created. The allowed values are:  
  `/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.MachineLearningServices/workspaces/{resourceName}`

* `name` - (Required) Specifies the name of the azure resource. Changing this forces a new resource to be created.

* `body` - (Required) Specifies the configuration of the resource. More information about the arguments in `body` can be found in the [Microsoft documentation](https://learn.microsoft.com/en-us/azure/templates/Microsoft.MachineLearningServices/workspaces/outboundRules?pivots=deployment-language-terraform).

For other arguments, please refer to the [azapi_resource](https://registry.terraform.io/providers/Azure/azapi/latest/docs/resources/resource) documentation.

## Import

 ```shell
 # Azure resource can be imported using the resource id, e.g.
 terraform import azapi_resource.example /subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.MachineLearningServices/workspaces/{resourceName}/outboundRules/{resourceName}
 
 # It also supports specifying API version by using the resource id with api-version as a query parameter, e.g.
 terraform import azapi_resource.example /subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.MachineLearningServices/workspaces/{resourceName}/outboundRules/{resourceName}?api-version=2025-07-01-preview
 ```
