---
subcategory: "Microsoft.Web - App Service, Azure Functions"
page_title: "sites/slots/config"
description: |-
  Manages a App Service's Slot Virtual Network Association.
---

# Microsoft.Web/sites/slots/config - App Service's Slot Virtual Network Association

This article demonstrates how to use `azapi` provider to manage the App Service's Slot Virtual Network Association resource in Azure.

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

resource "azapi_resource" "serverfarm" {
  type      = "Microsoft.Web/serverfarms@2022-09-01"
  parent_id = azapi_resource.resourceGroup.id
  name      = var.resource_name
  location  = var.location
  body = {
    properties = {
      hyperV         = false
      perSiteScaling = false
      reserved       = false
      zoneRedundant  = false
    }
    sku = {
      name = "S1"
    }
  }
  schema_validation_enabled = false
  response_export_values    = ["*"]
}

resource "azapi_resource" "site" {
  type      = "Microsoft.Web/sites@2022-09-01"
  parent_id = azapi_resource.resourceGroup.id
  name      = var.resource_name
  location  = var.location
  body = {
    properties = {
      clientAffinityEnabled = false
      clientCertEnabled     = false
      clientCertMode        = "Required"
      enabled               = true
      httpsOnly             = false
      publicNetworkAccess   = "Enabled"
      serverFarmId          = azapi_resource.serverfarm.id
      siteConfig = {
        acrUseManagedIdentityCreds       = false
        alwaysOn                         = true
        autoHealEnabled                  = false
        ftpsState                        = "Disabled"
        http20Enabled                    = false
        loadBalancing                    = "LeastRequests"
        localMySqlEnabled                = false
        managedPipelineMode              = "Integrated"
        minTlsVersion                    = "1.2"
        publicNetworkAccess              = "Enabled"
        remoteDebuggingEnabled           = false
        scmIpSecurityRestrictionsUseMain = false
        scmMinTlsVersion                 = "1.2"
        use32BitWorkerProcess            = true
        vnetRouteAllEnabled              = false
        webSocketsEnabled                = false
        windowsFxVersion                 = ""
      }
      vnetRouteAllEnabled = false
    }
  }
  schema_validation_enabled = false
  response_export_values    = ["*"]
}

resource "azapi_resource" "slot" {
  type      = "Microsoft.Web/sites/slots@2022-09-01"
  parent_id = azapi_resource.site.id
  name      = var.resource_name
  location  = var.location
  body = {
    properties = {
      clientAffinityEnabled    = false
      clientCertEnabled        = false
      clientCertExclusionPaths = ""
      clientCertMode           = "Required"
      enabled                  = true
      httpsOnly                = false
      publicNetworkAccess      = "Enabled"
      serverFarmId             = azapi_resource.serverfarm.id
      siteConfig = {
        acrUseManagedIdentityCreds       = false
        alwaysOn                         = true
        autoHealEnabled                  = false
        ftpsState                        = "Disabled"
        http20Enabled                    = false
        loadBalancing                    = "LeastRequests"
        localMySqlEnabled                = false
        managedPipelineMode              = "Integrated"
        minTlsVersion                    = "1.2"
        publicNetworkAccess              = "Enabled"
        remoteDebuggingEnabled           = false
        scmIpSecurityRestrictionsUseMain = false
        scmMinTlsVersion                 = "1.2"
        use32BitWorkerProcess            = false
        vnetRouteAllEnabled              = false
        webSocketsEnabled                = false
        windowsFxVersion                 = ""
      }
      vnetRouteAllEnabled = false
    }
  }
  schema_validation_enabled = false
  response_export_values    = ["*"]
}

data "azapi_resource_id" "config" {
  type      = "Microsoft.Web/sites/slots/config@2022-09-01"
  parent_id = azapi_resource.slot.id
  name      = "azurestorageaccounts"
}

resource "azapi_resource_action" "config" {
  type        = "Microsoft.Web/sites/slots/config@2022-09-01"
  resource_id = data.azapi_resource_id.config.id
  method      = "PUT"
  body = {
    name = "azurestorageaccounts"
    properties = {
    }
  }
  response_export_values = ["*"]
}


```



## Arguments Reference

The following arguments are supported:

* `type` - (Required) The type of the resource. This should be set to `Microsoft.Web/sites/slots/config@api-version`. The available api-versions for this resource are: [`2015-08-01`, `2016-08-01`, `2018-02-01`, `2018-11-01`, `2019-08-01`, `2020-06-01`, `2020-09-01`, `2020-10-01`, `2020-12-01`, `2021-01-01`, `2021-01-15`, `2021-02-01`, `2021-03-01`, `2022-03-01`, `2022-09-01`, `2023-01-01`, `2023-12-01`, `2024-04-01`].

* `parent_id` - (Required) The ID of the azure resource in which this resource is created. The allowed values are:  
  `/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Web/sites/{resourceName}/slots/{resourceName}`

* `name` - (Required) Specifies the name of the azure resource. Changing this forces a new resource to be created.

* `body` - (Required) Specifies the configuration of the resource. More information about the arguments in `body` can be found in the [Microsoft documentation](https://learn.microsoft.com/en-us/azure/templates/Microsoft.Web/sites/slots/config?pivots=deployment-language-terraform).

## Import

 ```shell
 # Azure resource can be imported using the resource id, e.g.
 terraform import azapi_resource.example /subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Web/sites/{resourceName}/slots/{resourceName}/config/{resourceName}
 
 # It also supports specifying API version by using the resource id with api-version as a query parameter, e.g.
 terraform import azapi_resource.example /subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Web/sites/{resourceName}/slots/{resourceName}/config/{resourceName}?api-version=2024-04-01
 ```
