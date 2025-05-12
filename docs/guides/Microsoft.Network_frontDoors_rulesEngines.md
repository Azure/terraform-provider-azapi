---
subcategory: "Microsoft.Network - Various Networking Services"
page_title: "frontDoors/rulesEngines"
description: |-
  Manages a Azure Front Door (classic) Rules Engine configuration and rules.
---

# Microsoft.Network/frontDoors/rulesEngines - Azure Front Door (classic) Rules Engine configuration and rules

This article demonstrates how to use `azapi` provider to manage the Azure Front Door (classic) Rules Engine configuration and rules resource in Azure.

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

resource "azapi_resource" "resourceGroup" {
  type     = "Microsoft.Resources/resourceGroups@2020-06-01"
  name     = var.resource_name
  location = var.location
}

locals {
  backend_name        = "backend-bing"
  endpoint_name       = "frontend-endpoint"
  health_probe_name   = "health-probe"
  load_balancing_name = "load-balancing-setting"
}

resource "azurerm_frontdoor" "test" {
  name                = "acctest-FD-test"
  resource_group_name = azapi_resource.resourceGroup.name

  backend_pool_settings {
    enforce_backend_pools_certificate_name_check = false
  }

  routing_rule {
    name               = "routing-rule"
    accepted_protocols = ["Http", "Https"]
    patterns_to_match  = ["/*"]
    frontend_endpoints = [local.endpoint_name]
    forwarding_configuration {
      forwarding_protocol = "MatchRequest"
      backend_pool_name   = local.backend_name
    }
  }

  backend_pool_load_balancing {
    name = local.load_balancing_name
  }

  backend_pool_health_probe {
    name = local.health_probe_name
  }

  backend_pool {
    name = local.backend_name
    backend {
      host_header = "www.bing.com"
      address     = "www.bing.com"
      http_port   = 80
      https_port  = 443
    }

    load_balancing_name = local.load_balancing_name
    health_probe_name   = local.health_probe_name
  }

  frontend_endpoint {
    name      = local.endpoint_name
    host_name = "acctest-FD-test.azurefd.net"
  }
}

resource "azapi_resource" "rulesEngine" {
  type      = "Microsoft.Network/frontDoors/rulesEngines@2020-05-01"
  parent_id = azurerm_frontdoor.test.id
  name      = var.resource_name
  body = {
    properties = {
      rules = [
        {
          name     = var.resource_name
          priority = 0
          action = {
            routeConfigurationOverride = {
              redirectType     = "Found"
              redirectProtocol = "HttpsOnly"
              customHost       = "customhost.org"
              "@odata.type"    = "#Microsoft.Azure.FrontDoor.Models.FrontdoorRedirectConfiguration"
            }
          }
          matchProcessingBehavior = "Continue"
        }
      ]
    }
  }
}

```



## Arguments Reference

The following arguments are supported:

* `type` - (Required) The type of the resource. This should be set to `Microsoft.Network/frontDoors/rulesEngines@api-version`. The available api-versions for this resource are: [`2020-01-01`, `2020-04-01`, `2020-05-01`, `2021-06-01`].

* `parent_id` - (Required) The ID of the azure resource in which this resource is created. The allowed values are:  
  `/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Network/frontDoors/{resourceName}`

* `name` - (Required) Specifies the name of the azure resource. Changing this forces a new resource to be created.

* `body` - (Required) Specifies the configuration of the resource. More information about the arguments in `body` can be found in the [Microsoft documentation](https://learn.microsoft.com/en-us/azure/templates/Microsoft.Network/frontDoors/rulesEngines?pivots=deployment-language-terraform).

For other arguments, please refer to the [azapi_resource](https://registry.terraform.io/providers/Azure/azapi/latest/docs/resources/resource) documentation.

## Import

 ```shell
 # Azure resource can be imported using the resource id, e.g.
 terraform import azapi_resource.example /subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Network/frontDoors/{resourceName}/rulesEngines/{resourceName}
 
 # It also supports specifying API version by using the resource id with api-version as a query parameter, e.g.
 terraform import azapi_resource.example /subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Network/frontDoors/{resourceName}/rulesEngines/{resourceName}?api-version=2021-06-01
 ```
