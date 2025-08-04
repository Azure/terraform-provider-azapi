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

