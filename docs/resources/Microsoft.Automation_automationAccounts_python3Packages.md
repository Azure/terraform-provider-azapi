---
subcategory: "Microsoft.Automation - Automation"
page_title: "automationAccounts/python3Packages"
description: |-
  Manages a Automation Python3 Package.
---

# Microsoft.Automation/automationAccounts/python3Packages - Automation Python3 Package

This article demonstrates how to use `azapi` provider to manage the Automation Python3 Package resource in Azure.

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

resource "azapi_resource" "resourceGroup" {
  type     = "Microsoft.Resources/resourceGroups@2020-06-01"
  name     = var.resource_name
  location = var.location
}

resource "azapi_resource" "automationAccount" {
  type      = "Microsoft.Automation/automationAccounts@2023-11-01"
  parent_id = azapi_resource.resourceGroup.id
  name      = var.resource_name
  location  = var.location
  body = {
    properties = {
      disableLocalAuth = false
      encryption = {
        keySource = "Microsoft.Automation"
      }
      publicNetworkAccess = true
      sku = {
        name = "Basic"
      }
    }
  }
}

resource "azapi_resource" "python3Package" {
  type      = "Microsoft.Automation/automationAccounts/python3Packages@2023-11-01"
  parent_id = azapi_resource.automationAccount.id
  name      = var.resource_name
  body = {
    properties = {
      contentLink = {
        uri     = "https://files.pythonhosted.org/packages/py3/r/requests/requests-2.31.0-py3-none-any.whl"
        version = "2.31.0"
      }
    }
  }
  tags = {
    key = "foo"
  }
}


```



## Arguments Reference

The following arguments are supported:

* `type` - (Required) The type of the resource. This should be set to `Microsoft.Automation/automationAccounts/python3Packages@api-version`. The available api-versions for this resource are: [`2022-08-08`, `2023-05-15-preview`, `2023-11-01`, `2024-10-23`].

* `parent_id` - (Required) The ID of the azure resource in which this resource is created. The allowed values are:  
  `/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Automation/automationAccounts/{resourceName}`

* `name` - (Required) Specifies the name of the azure resource. Changing this forces a new resource to be created.

* `body` - (Required) Specifies the configuration of the resource. More information about the arguments in `body` can be found in the [Microsoft documentation](https://learn.microsoft.com/en-us/azure/templates/Microsoft.Automation/automationAccounts/python3Packages?pivots=deployment-language-terraform).

For other arguments, please refer to the [azapi_resource](https://registry.terraform.io/providers/Azure/azapi/latest/docs/resources/resource) documentation.

## Import

 ```shell
 # Azure resource can be imported using the resource id, e.g.
 terraform import azapi_resource.example /subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Automation/automationAccounts/{resourceName}/python3Packages/{resourceName}
 
 # It also supports specifying API version by using the resource id with api-version as a query parameter, e.g.
 terraform import azapi_resource.example /subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Automation/automationAccounts/{resourceName}/python3Packages/{resourceName}?api-version=2024-10-23
 ```
