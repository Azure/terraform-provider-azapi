---
subcategory: "Microsoft.MachineLearningServices - Azure Machine Learning"
page_title: "workspaces/codes/versions"
description: |-
  Manages a Machine Learning Services Workspace Codes Versions.
---

# Microsoft.MachineLearningServices/workspaces/codes/versions - Machine Learning Services Workspace Codes Versions

This article demonstrates how to use `azapi` provider to manage the Machine Learning Services Workspace Codes Versions resource in Azure.

## Example Usage

### default

```hcl
terraform {
  required_providers {
    azapi = {
      source = "Azure/azapi"
    }
    azurerm = {
      source = "hashicorp/azurerm"
    }
  }
}

provider "azurerm" {
  features {
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

data "azapi_client_config" "current" {
}

resource "azapi_resource" "resourceGroup" {
  type                      = "Microsoft.Resources/resourceGroups@2020-06-01"
  name                      = var.resource_name
  location                  = var.location
  body                      = {}
  schema_validation_enabled = false
  response_export_values    = ["*"]
}

resource "azapi_resource" "storageAccount" {
  type      = "Microsoft.Storage/storageAccounts@2021-09-01"
  parent_id = azapi_resource.resourceGroup.id
  name      = var.resource_name
  location  = var.location
  body = {
    kind = "StorageV2"
    properties = {
      accessTier                   = "Hot"
      allowBlobPublicAccess        = true
      allowCrossTenantReplication  = true
      allowSharedKeyAccess         = true
      defaultToOAuthAuthentication = false
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
      isHnsEnabled      = false
      isNfsV3Enabled    = false
      isSftpEnabled     = false
      minimumTlsVersion = "TLS1_2"
      networkAcls = {
        defaultAction = "Allow"
      }
      publicNetworkAccess      = "Enabled"
      supportsHttpsTrafficOnly = true
    }
    sku = {
      name = "Standard_LRS"
    }
  }
}

data "azapi_resource" "blobService" {
  type      = "Microsoft.Storage/storageAccounts/blobServices@2022-09-01"
  parent_id = azapi_resource.storageAccount.id
  name      = "default"
}

resource "azapi_resource" "container" {
  type      = "Microsoft.Storage/storageAccounts/blobServices/containers@2022-09-01"
  name      = var.resource_name
  parent_id = data.azapi_resource.blobService.id
  body = {
    properties = {
      metadata = {
        key = "value"
      }
    }
  }
}

resource "azurerm_storage_blob" "example" {
  name                   = "codes.py"
  storage_account_name   = azapi_resource.storageAccount.name
  storage_container_name = azapi_resource.container.name
  type                   = "Block"
  source                 = "/Users/luheng/go/playground/issues/icm/azapi_ml_code/output.md"
}

resource "azapi_resource" "component" {
  type      = "Microsoft.Insights/components@2020-02-02"
  parent_id = azapi_resource.resourceGroup.id
  name      = var.resource_name
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
  type      = "Microsoft.KeyVault/vaults@2021-10-01"
  parent_id = azapi_resource.resourceGroup.id
  name      = var.resource_name
  location  = var.location
  body = {
    properties = {
      accessPolicies = [
        {
          objectId = "45a2d1ea-488a-44b0-bb2e-3cd8e485ebef"
          permissions = {
            certificates = [
              "all",
            ]
            keys = [
              "all",
            ]
            secrets = [
              "all",
            ]
            storage = []
          }
          tenantId = data.azapi_client_config.current.tenant_id
        }
      ]
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
  lifecycle {
    ignore_changes = [body.properties.accessPolicies]
  }
}

resource "azapi_resource" "workspace" {
  type      = "Microsoft.MachineLearningServices/workspaces@2022-05-01"
  parent_id = azapi_resource.resourceGroup.id
  name      = var.resource_name
  location  = var.location
  identity {
    type         = "SystemAssigned"
    identity_ids = []
  }
  body = {
    properties = {
      applicationInsights = azapi_resource.component.id
      keyVault            = azapi_resource.vault.id
      publicNetworkAccess = "Disabled"
      storageAccount      = azapi_resource.storageAccount.id
      v1LegacyMode        = false
    }
    sku = {
      name = "Basic"
      tier = "Basic"
    }
  }
  ignore_casing             = true
  schema_validation_enabled = false
  response_export_values    = ["*"]
}

data "azapi_resource_id" "code" {
  type      = "Microsoft.MachineLearningServices/workspaces/codes@2024-10-01"
  parent_id = azapi_resource.workspace.id
  name      = "mycode"
}

resource "azapi_resource" "codeVersion" {
  type      = "Microsoft.MachineLearningServices/workspaces/codes/versions@2024-10-01"
  parent_id = data.azapi_resource_id.code.id
  name      = "1"
  body = {
    properties = {
      codeUri     = "${azapi_resource.storageAccount.output.properties.primaryEndpoints.blob}${azapi_resource.container.name}"
      description = "this is a test code version"
      isArchived  = false
      tags = {
        env = "prod"
      }
      properties = {

      }
    }
  }
  depends_on = [azurerm_storage_blob.example]
}

```



## Arguments Reference

The following arguments are supported:

* `type` - (Required) The type of the resource. This should be set to `Microsoft.MachineLearningServices/workspaces/codes/versions@api-version`. The available api-versions for this resource are: [`2021-03-01-preview`, `2022-02-01-preview`, `2022-05-01`, `2022-06-01-preview`, `2022-10-01`, `2022-10-01-preview`, `2022-12-01-preview`, `2023-02-01-preview`, `2023-04-01`, `2023-04-01-preview`, `2023-06-01-preview`, `2023-08-01-preview`, `2023-10-01`, `2024-01-01-preview`, `2024-04-01`, `2024-04-01-preview`, `2024-07-01-preview`, `2024-10-01`, `2024-10-01-preview`, `2025-01-01-preview`, `2025-04-01`, `2025-04-01-preview`, `2025-06-01`].

* `parent_id` - (Required) The ID of the azure resource in which this resource is created. The allowed values are:  
  `/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.MachineLearningServices/workspaces/{resourceName}/codes/{resourceName}`

* `name` - (Required) Specifies the name of the azure resource. Changing this forces a new resource to be created.

* `body` - (Required) Specifies the configuration of the resource. More information about the arguments in `body` can be found in the [Microsoft documentation](https://learn.microsoft.com/en-us/azure/templates/Microsoft.MachineLearningServices/workspaces/codes/versions?pivots=deployment-language-terraform).

For other arguments, please refer to the [azapi_resource](https://registry.terraform.io/providers/Azure/azapi/latest/docs/resources/resource) documentation.

## Import

 ```shell
 # Azure resource can be imported using the resource id, e.g.
 terraform import azapi_resource.example /subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.MachineLearningServices/workspaces/{resourceName}/codes/{resourceName}/versions/{resourceName}
 
 # It also supports specifying API version by using the resource id with api-version as a query parameter, e.g.
 terraform import azapi_resource.example /subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.MachineLearningServices/workspaces/{resourceName}/codes/{resourceName}/versions/{resourceName}?api-version=2025-06-01
 ```
