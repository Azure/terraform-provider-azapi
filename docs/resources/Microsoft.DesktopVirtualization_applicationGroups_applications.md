---
subcategory: "Microsoft.DesktopVirtualization - Azure Virtual Desktop"
page_title: "applicationGroups/applications"
description: |-
  Manages a Virtual Desktop Application.
---

# Microsoft.DesktopVirtualization/applicationGroups/applications - Virtual Desktop Application

This article demonstrates how to use `azapi` provider to manage the Virtual Desktop Application resource in Azure.

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
  default = "westeurope"
}

resource "azapi_resource" "resourceGroup" {
  type     = "Microsoft.Resources/resourceGroups@2020-06-01"
  name     = var.resource_name
  location = var.location
}

resource "azapi_resource" "hostPool" {
  type      = "Microsoft.DesktopVirtualization/hostPools@2023-09-05"
  parent_id = azapi_resource.resourceGroup.id
  name      = var.resource_name
  location  = var.location
  body = {
    properties = {
      hostPoolType          = "Pooled"
      loadBalancerType      = "BreadthFirst"
      maxSessionLimit       = 999999
      preferredAppGroupType = "Desktop"
      publicNetworkAccess   = "Enabled"
      startVMOnConnect      = false
      validationEnvironment = false
    }
  }
  schema_validation_enabled = false
  response_export_values    = ["*"]
}

resource "azapi_resource" "applicationGroup" {
  type      = "Microsoft.DesktopVirtualization/applicationGroups@2023-09-05"
  parent_id = azapi_resource.resourceGroup.id
  name      = var.resource_name
  location  = var.location
  body = {
    properties = {
      applicationGroupType = "RemoteApp"
      hostPoolArmPath      = azapi_resource.hostPool.id
    }
  }
  schema_validation_enabled = false
  ignore_casing             = true
  response_export_values    = ["*"]
}

resource "azapi_resource" "application" {
  type      = "Microsoft.DesktopVirtualization/applicationGroups/applications@2023-09-05"
  parent_id = azapi_resource.applicationGroup.id
  name      = var.resource_name
  location  = var.location
  body = {
    properties = {
      commandLineSetting = "DoNotAllow"
      filePath           = "C:\\Program Files\\Google\\Chrome\\Application\\chrome.exe"
      showInPortal       = false
    }
  }
  schema_validation_enabled = false
  response_export_values    = ["*"]
}

```



## Arguments Reference

The following arguments are supported:

* `type` - (Required) The type of the resource. This should be set to `Microsoft.DesktopVirtualization/applicationGroups/applications@api-version`. The available api-versions for this resource are: [`2019-01-23-preview`, `2019-09-24-preview`, `2019-12-10-preview`, `2020-09-21-preview`, `2020-10-19-preview`, `2020-11-02-preview`, `2020-11-10-preview`, `2021-01-14-preview`, `2021-02-01-preview`, `2021-03-09-preview`, `2021-04-01-preview`, `2021-07-12`, `2021-09-03-preview`, `2022-02-10-preview`, `2022-04-01-preview`, `2022-09-09`, `2022-10-14-preview`, `2023-09-05`, `2023-10-04-preview`, `2023-11-01-preview`, `2024-01-16-preview`, `2024-03-06-preview`, `2024-04-03`, `2024-04-08-preview`, `2024-08-08-preview`, `2024-11-01-preview`, `2025-03-01-preview`].

* `parent_id` - (Required) The ID of the azure resource in which this resource is created. The allowed values are:  
  `/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.DesktopVirtualization/applicationGroups/{resourceName}`

* `name` - (Required) Specifies the name of the azure resource. Changing this forces a new resource to be created.

* `body` - (Required) Specifies the configuration of the resource. More information about the arguments in `body` can be found in the [Microsoft documentation](https://learn.microsoft.com/en-us/azure/templates/Microsoft.DesktopVirtualization/applicationGroups/applications?pivots=deployment-language-terraform).

For other arguments, please refer to the [azapi_resource](https://registry.terraform.io/providers/Azure/azapi/latest/docs/resources/resource) documentation.

## Import

 ```shell
 # Azure resource can be imported using the resource id, e.g.
 terraform import azapi_resource.example /subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.DesktopVirtualization/applicationGroups/{resourceName}/applications/{resourceName}
 
 # It also supports specifying API version by using the resource id with api-version as a query parameter, e.g.
 terraform import azapi_resource.example /subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.DesktopVirtualization/applicationGroups/{resourceName}/applications/{resourceName}?api-version=2025-03-01-preview
 ```
