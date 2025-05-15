---
subcategory: "Microsoft.Solutions - Azure Managed Applications"
page_title: "applicationDefinitions"
description: |-
  Manages a Managed Application Definition.
---

# Microsoft.Solutions/applicationDefinitions - Managed Application Definition

This article demonstrates how to use `azapi` provider to manage the Managed Application Definition resource in Azure.

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

data "azapi_resource_action" "roleDefinitions" {
  type                   = "Microsoft.Authorization@2018-01-01-preview"
  resource_id            = "/providers/Microsoft.Authorization"
  action                 = "roleDefinitions"
  method                 = "GET"
  response_export_values = ["*"]
}

resource "azapi_resource" "resourceGroup" {
  type     = "Microsoft.Resources/resourceGroups@2020-06-01"
  name     = var.resource_name
  location = var.location
}

resource "azapi_resource" "applicationDefinition" {
  type      = "Microsoft.Solutions/applicationDefinitions@2021-07-01"
  parent_id = azapi_resource.resourceGroup.id
  name      = var.resource_name
  location  = var.location
  body = {
    properties = {
      authorizations = [
        {
          principalId      = data.azurerm_client_config.current.object_id
          roleDefinitionId = data.azapi_resource_action.roleDefinitions.output.value[0].name
        },
      ]
      description    = "Test Managed App Definition"
      displayName    = "TestManagedAppDefinition"
      isEnabled      = true
      lockLevel      = "ReadOnly"
      packageFileUri = "https://github.com/Azure/azure-managedapp-samples/raw/master/Managed Application Sample Packages/201-managed-storage-account/managedstorage.zip"
    }
  }
  schema_validation_enabled = false
  response_export_values    = ["*"]
}


```



## Arguments Reference

The following arguments are supported:

* `type` - (Required) The type of the resource. This should be set to `Microsoft.Solutions/applicationDefinitions@api-version`. The available api-versions for this resource are: [`2017-09-01`, `2017-12-01`, `2018-02-01`, `2018-03-01`, `2018-06-01`, `2018-09-01-preview`, `2019-07-01`, `2020-08-21-preview`, `2021-02-01-preview`, `2021-07-01`, `2023-12-01-preview`].

* `parent_id` - (Required) The ID of the azure resource in which this resource is created. The allowed values are:  
  `/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}`

* `name` - (Required) Specifies the name of the azure resource. Changing this forces a new resource to be created.

* `body` - (Required) Specifies the configuration of the resource. More information about the arguments in `body` can be found in the [Microsoft documentation](https://learn.microsoft.com/en-us/azure/templates/Microsoft.Solutions/applicationDefinitions?pivots=deployment-language-terraform).

For other arguments, please refer to the [azapi_resource](https://registry.terraform.io/providers/Azure/azapi/latest/docs/resources/resource) documentation.

## Import

 ```shell
 # Azure resource can be imported using the resource id, e.g.
 terraform import azapi_resource.example /subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Solutions/applicationDefinitions/{resourceName}
 
 # It also supports specifying API version by using the resource id with api-version as a query parameter, e.g.
 terraform import azapi_resource.example /subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Solutions/applicationDefinitions/{resourceName}?api-version=2023-12-01-preview
 ```
