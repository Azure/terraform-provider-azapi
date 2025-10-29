---
subcategory: "Microsoft.MachineLearningServices - Azure Machine Learning"
page_title: "workspaces/computes"
description: |-
  Manages a Machine Learning Compute.
---

# Microsoft.MachineLearningServices/workspaces/computes - Machine Learning Compute

This article demonstrates how to use `azapi` provider to manage the Machine Learning Compute resource in Azure.



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

data "azurerm_client_config" "current" {
}

resource "azapi_resource" "resourceGroup" {
  type     = "Microsoft.Resources/resourceGroups@2020-06-01"
  name     = var.resource_name
  location = var.location
  body = {
    tags = {
      stage = "test"
    }
  }
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
    sku = {
      name = "Standard_LRS"
    }
  }
  schema_validation_enabled = false
  response_export_values    = ["*"]
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
  schema_validation_enabled = false
  response_export_values    = ["*"]
}

resource "azapi_resource" "vault" {
  type      = "Microsoft.KeyVault/vaults@2021-10-01"
  parent_id = azapi_resource.resourceGroup.id
  name      = var.resource_name
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
      tenantId = data.azurerm_client_config.current.tenant_id
    }
  }
  schema_validation_enabled = false
  response_export_values    = ["*"]
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
      publicNetworkAccess = "Enabled"
      storageAccount      = azapi_resource.storageAccount.id
      v1LegacyMode        = false
    }
    sku = {
      name = "Basic"
      tier = "Basic"
    }
  }
  schema_validation_enabled = false
  ignore_casing             = true
  response_export_values    = ["*"]
}

resource "azapi_resource" "compute" {
  type      = "Microsoft.MachineLearningServices/workspaces/computes@2022-05-01"
  parent_id = azapi_resource.workspace.id
  name      = var.resource_name
  location  = var.location
  body = {
    properties = {
      computeLocation  = "westeurope"
      computeType      = "ComputeInstance"
      description      = ""
      disableLocalAuth = true
      properties = {
        vmSize = "STANDARD_D2_V2"
      }
    }
  }
  schema_validation_enabled = false
  response_export_values    = ["*"]
}


```



## Arguments Reference

The following arguments are supported:

* `type` - (Required) The type of the resource. This should be set to `Microsoft.MachineLearningServices/workspaces/computes@api-version`. The available api-versions for this resource are: [`2018-03-01-preview`, `2018-11-19`, `2019-05-01`, `2019-06-01`, `2019-11-01`, `2020-01-01`, `2020-02-18-preview`, `2020-03-01`, `2020-04-01`, `2020-05-01-preview`, `2020-05-15-preview`, `2020-06-01`, `2020-08-01`, `2020-09-01-preview`, `2021-01-01`, `2021-03-01-preview`, `2021-04-01`, `2021-07-01`, `2022-01-01-preview`, `2022-02-01-preview`, `2022-05-01`, `2022-06-01-preview`, `2022-10-01`, `2022-10-01-preview`, `2022-12-01-preview`, `2023-02-01-preview`, `2023-04-01`, `2023-04-01-preview`, `2023-06-01-preview`, `2023-08-01-preview`, `2023-10-01`, `2024-01-01-preview`, `2024-04-01`, `2024-04-01-preview`, `2024-07-01-preview`, `2024-10-01`, `2024-10-01-preview`, `2025-01-01-preview`, `2025-04-01`, `2025-04-01-preview`, `2025-06-01`, `2025-07-01-preview`, `2025-09-01`].

* `parent_id` - (Required) The ID of the azure resource in which this resource is created. The allowed values are:  
  `/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.MachineLearningServices/workspaces/{resourceName}`

* `name` - (Required) Specifies the name of the azure resource. Changing this forces a new resource to be created.

* `body` - (Required) Specifies the configuration of the resource. More information about the arguments in `body` can be found in the [Microsoft documentation](https://learn.microsoft.com/en-us/azure/templates/Microsoft.MachineLearningServices/workspaces/computes?pivots=deployment-language-terraform).

For other arguments, please refer to the [azapi_resource](https://registry.terraform.io/providers/Azure/azapi/latest/docs/resources/resource) documentation.

## Import

 ```shell
 # Azure resource can be imported using the resource id, e.g.
 terraform import azapi_resource.example /subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.MachineLearningServices/workspaces/{resourceName}/computes/{resourceName}
 
 # It also supports specifying API version by using the resource id with api-version as a query parameter, e.g.
 terraform import azapi_resource.example /subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.MachineLearningServices/workspaces/{resourceName}/computes/{resourceName}?api-version=2025-09-01
 ```
