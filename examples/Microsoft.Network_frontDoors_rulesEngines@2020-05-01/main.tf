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
