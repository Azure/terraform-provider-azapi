terraform {
  required_providers {
    azapi = {
      source = "Azure/azapi"
    }
    time = {
      source = "Hashicorp/time"
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

data "azapi_client_config" "current" {}

resource "azapi_resource" "networkManager" {
  type      = "Microsoft.Network/networkManagers@2024-05-01"
  parent_id = azapi_resource.resourceGroup.id
  name      = var.resource_name
  location  = var.location
  body = {
    properties = {
      description                 = ""
      networkManagerScopeAccesses = []
      networkManagerScopes = {
        managementGroups = []
        subscriptions = [
          "/subscriptions/${data.azapi_client_config.current.subscription_id}"
        ]
      }
    }
  }
  tags = {
    "sampleTag" : "sampleTag"
  }
  schema_validation_enabled = false
  response_export_values    = ["*"]
}

resource "time_sleep" "wait_60_seconds" {
  depends_on = [azapi_resource.networkManager]

  destroy_duration = "60s"
}

resource "azapi_resource" "ipamPool" {
  type      = "Microsoft.Network/networkManagers/ipamPools@2024-05-01"
  parent_id = azapi_resource.networkManager.id
  name      = var.resource_name
  location  = var.location
  body = {
    properties = {
      addressPrefixes = [
        "10.0.0.0/24",
      ]
      description    = "Test description."
      parentPoolName = ""
      displayName    = "testDisplayName"
    }
  }
  tags = {
    "myTag" : "testTag"
  }
  schema_validation_enabled = false
  ignore_casing             = false
  ignore_missing_property   = false

  depends_on = [time_sleep.wait_60_seconds]
}

resource "azapi_resource" "vnet_withIPAM" {
  type      = "Microsoft.Network/virtualNetworks@2024-05-01"
  parent_id = azapi_resource.resourceGroup.id
  name      = var.resource_name
  location  = var.location
  body = {
    properties = {
      addressSpace = {
        ipamPoolPrefixAllocations = [
          {
            numberOfIpAddresses = "100",
            pool = {
              id = azapi_resource.ipamPool.id
            }
          }
        ]
      }
    }
  }
  lifecycle {
    ignore_changes = [body.properties.subnets]
  }
  schema_validation_enabled = false
  ignore_casing             = false
  ignore_missing_property   = false
}

resource "azapi_resource" "subnet_withIPAM" {
  type      = "Microsoft.Network/virtualNetworks/subnets@2024-05-01"
  parent_id = azapi_resource.vnet_withIPAM.id
  name      = var.resource_name
  body = {
    properties = {
      ipamPoolPrefixAllocations = [
        {
          numberOfIpAddresses = "100",
          pool = {
            id = azapi_resource.ipamPool.id
          }
        }
      ]
    }
  }
  schema_validation_enabled = false
  response_export_values    = ["*"]
}
