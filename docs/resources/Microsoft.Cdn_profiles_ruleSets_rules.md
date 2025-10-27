---
subcategory: "Microsoft.Cdn - Content Delivery Network"
page_title: "profiles/ruleSets/rules"
description: |-
  Manages a Front Door (standard/premium) Rule.
---

# Microsoft.Cdn/profiles/ruleSets/rules - Front Door (standard/premium) Rule

This article demonstrates how to use `azapi` provider to manage the Front Door (standard/premium) Rule resource in Azure.



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

variable "cdn_location" {
  type    = string
  default = "global"
}

resource "azapi_resource" "resourceGroup" {
  type     = "Microsoft.Resources/resourceGroups@2020-06-01"
  name     = var.resource_name
  location = var.location
}

resource "azapi_resource" "profile" {
  type      = "Microsoft.Cdn/profiles@2024-09-01"
  parent_id = azapi_resource.resourceGroup.id
  name      = "${var.resource_name}-profile"
  location  = var.cdn_location
  body = {
    properties = {
      originResponseTimeoutSeconds = 120
    }
    sku = {
      name = "Standard_AzureFrontDoor"
    }
  }
}

resource "azapi_resource" "originGroup" {
  type      = "Microsoft.Cdn/profiles/originGroups@2024-09-01"
  parent_id = azapi_resource.profile.id
  name      = "${var.resource_name}-origingroup"
  body = {
    properties = {
      loadBalancingSettings = {
        additionalLatencyInMilliseconds = 0
        sampleSize                      = 16
        successfulSamplesRequired       = 3
      }
      sessionAffinityState                                  = "Enabled"
      trafficRestorationTimeToHealedOrNewEndpointsInMinutes = 10
    }
  }
}

resource "azapi_resource" "ruleSet" {
  type      = "Microsoft.Cdn/profiles/ruleSets@2024-09-01"
  parent_id = azapi_resource.profile.id
  name      = "ruleSet${substr(var.resource_name, -4, -1)}"
}

resource "azapi_resource" "origin" {
  type      = "Microsoft.Cdn/profiles/originGroups/origins@2024-09-01"
  parent_id = azapi_resource.originGroup.id
  name      = "${var.resource_name}-origin"
  body = {
    properties = {
      enabledState                = "Enabled"
      enforceCertificateNameCheck = false
      hostName                    = "contoso.com"
      httpPort                    = 80
      httpsPort                   = 443
      originHostHeader            = "www.contoso.com"
      priority                    = 1
      weight                      = 1
    }
  }
}

resource "azapi_resource" "rule" {
  type      = "Microsoft.Cdn/profiles/ruleSets/rules@2024-09-01"
  parent_id = azapi_resource.ruleSet.id
  name      = "rule${substr(var.resource_name, -4, -1)}"
  body = {
    properties = {
      actions = [{
        name = "RouteConfigurationOverride"
        parameters = {
          cacheConfiguration = {
            cacheBehavior              = "OverrideIfOriginMissing"
            cacheDuration              = "23:59:59"
            isCompressionEnabled       = "Disabled"
            queryParameters            = "clientIp={client_ip}"
            queryStringCachingBehavior = "IgnoreSpecifiedQueryStrings"
          }
          originGroupOverride = {
            forwardingProtocol = "HttpsOnly"
            originGroup = {
              id = azapi_resource.originGroup.id
            }
          }
          typeName = "DeliveryRuleRouteConfigurationOverrideActionParameters"
        }
      }]
      conditions              = []
      matchProcessingBehavior = "Continue"
      order                   = 1
    }
  }
}


```



## Arguments Reference

The following arguments are supported:

* `type` - (Required) The type of the resource. This should be set to `Microsoft.Cdn/profiles/ruleSets/rules@api-version`. The available api-versions for this resource are: [`2020-09-01`, `2021-06-01`, `2022-05-01-preview`, `2022-11-01-preview`, `2023-05-01`, `2023-07-01-preview`, `2024-02-01`, `2024-05-01-preview`, `2024-06-01-preview`, `2024-09-01`, `2025-01-01-preview`, `2025-04-15`, `2025-05-01-preview`, `2025-06-01`].

* `parent_id` - (Required) The ID of the azure resource in which this resource is created. The allowed values are:  
  `/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Cdn/profiles/{resourceName}/ruleSets/{resourceName}`

* `name` - (Required) Specifies the name of the azure resource. Changing this forces a new resource to be created.

* `body` - (Required) Specifies the configuration of the resource. More information about the arguments in `body` can be found in the [Microsoft documentation](https://learn.microsoft.com/en-us/azure/templates/Microsoft.Cdn/profiles/ruleSets/rules?pivots=deployment-language-terraform).

For other arguments, please refer to the [azapi_resource](https://registry.terraform.io/providers/Azure/azapi/latest/docs/resources/resource) documentation.

## Import

 ```shell
 # Azure resource can be imported using the resource id, e.g.
 terraform import azapi_resource.example /subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Cdn/profiles/{resourceName}/ruleSets/{resourceName}/rules/{resourceName}
 
 # It also supports specifying API version by using the resource id with api-version as a query parameter, e.g.
 terraform import azapi_resource.example /subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Cdn/profiles/{resourceName}/ruleSets/{resourceName}/rules/{resourceName}?api-version=2025-06-01
 ```
