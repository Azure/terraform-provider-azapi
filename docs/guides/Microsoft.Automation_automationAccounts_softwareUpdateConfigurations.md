---
subcategory: "Microsoft.Automation - Automation"
page_title: "automationAccounts/softwareUpdateConfigurations"
description: |-
  Manages a Automation Software Update Configuration.
---

# Microsoft.Automation/automationAccounts/softwareUpdateConfigurations - Automation Software Update Configuration

This article demonstrates how to use `azapi` provider to manage the Automation Software Update Configuration resource in Azure.

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

resource "azapi_resource" "automationAccount" {
  type      = "Microsoft.Automation/automationAccounts@2021-06-22"
  parent_id = azapi_resource.resourceGroup.id
  name      = var.resource_name
  location  = var.location
  body = {
    properties = {
      encryption = {
        keySource = "Microsoft.Automation"
      }
      publicNetworkAccess = true
      sku = {
        name = "Basic"
      }
    }
  }
  schema_validation_enabled = false
  response_export_values    = ["*"]
}

resource "azapi_resource" "softwareUpdateConfiguration" {
  type      = "Microsoft.Automation/automationAccounts/softwareUpdateConfigurations@2019-06-01"
  parent_id = azapi_resource.automationAccount.id
  name      = var.resource_name
  body = {
    properties = {
      scheduleInfo = {
        description             = ""
        expiryTimeOffsetMinutes = 0
        frequency               = "OneTime"
        interval                = 0
        isEnabled               = true
        nextRunOffsetMinutes    = 0
        startTimeOffsetMinutes  = 0
        timeZone                = "Etc/UTC"
      }
      updateConfiguration = {
        duration = "PT2H"
        linux = {
          excludedPackageNameMasks = [
          ]
          includedPackageClassifications = "Security"
          includedPackageNameMasks = [
          ]
          rebootSetting = "IfRequired"
        }
        operatingSystem = "Linux"
        targets = {
          azureQueries = [
            {
              locations = [
                "westeurope",
              ]
              scope = [
                azapi_resource.resourceGroup.id,
              ]
            },
          ]
        }
      }
    }
  }
  schema_validation_enabled = false
  response_export_values    = ["*"]
}


```



## Arguments Reference

The following arguments are supported:

* `type` - (Required) The type of the resource. This should be set to `Microsoft.Automation/automationAccounts/softwareUpdateConfigurations@api-version`. The available api-versions for this resource are: [`2017-05-15-preview`, `2019-06-01`, `2023-05-15-preview`, `2024-10-23`].

* `parent_id` - (Required) The ID of the azure resource in which this resource is created. The allowed values are:  
  `/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Automation/automationAccounts/{resourceName}`

* `name` - (Required) Specifies the name of the azure resource. Changing this forces a new resource to be created.

* `body` - (Required) Specifies the configuration of the resource. More information about the arguments in `body` can be found in the [Microsoft documentation](https://learn.microsoft.com/en-us/azure/templates/Microsoft.Automation/automationAccounts/softwareUpdateConfigurations?pivots=deployment-language-terraform).

## Import

 ```shell
 # Azure resource can be imported using the resource id, e.g.
 terraform import azapi_resource.example /subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Automation/automationAccounts/{resourceName}/softwareUpdateConfigurations/{resourceName}
 
 # It also supports specifying API version by using the resource id with api-version as a query parameter, e.g.
 terraform import azapi_resource.example /subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Automation/automationAccounts/{resourceName}/softwareUpdateConfigurations/{resourceName}?api-version=2024-10-23
 ```
