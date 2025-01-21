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

resource "azapi_resource" "virtualNetwork" {
  type      = "Microsoft.Network/virtualNetworks@2022-07-01"
  parent_id = azapi_resource.resourceGroup.id
  name      = var.resource_name
  location  = var.location
  body = {
    properties = {
      addressSpace = {
        addressPrefixes = [
          "10.0.0.0/16",
        ]
      }
    }
  }
  schema_validation_enabled = false
  response_export_values    = ["*"]
  lifecycle {
    ignore_changes = [body.properties.subnets]
  }
}

resource "azapi_resource" "subnet" {
  type      = "Microsoft.Network/virtualNetworks/subnets@2022-07-01"
  parent_id = azapi_resource.virtualNetwork.id
  name      = var.resource_name
  body = {
    properties = {
      addressPrefix = "10.0.2.0/24"
    }
  }
  schema_validation_enabled = false
  response_export_values    = ["*"]
}

resource "azapi_resource" "networkConnection" {
  type      = "Microsoft.DevCenter/networkConnections@2023-04-01"
  name      = var.resource_name
  parent_id = azapi_resource.resourceGroup.id
  body = {
    location = var.location
    properties = {
      domainJoinType = "AzureADJoin"
      subnetId       = azapi_resource.subnet.id
    }
  }
  schema_validation_enabled = false
  response_export_values    = ["*"]
}

resource "azapi_resource" "devCenter" {
  type      = "Microsoft.DevCenter/devcenters@2023-04-01"
  parent_id = azapi_resource.resourceGroup.id
  name      = var.resource_name
  location  = var.location
  body = {
    identity = {
      type                   = "SystemAssigned"
      userAssignedIdentities = null
    }
  }
  schema_validation_enabled = false
  response_export_values    = ["*"]
}

resource "azapi_resource" "attachNetwork" {
  type      = "Microsoft.DevCenter/devcenters/attachednetworks@2023-04-01"
  name      = var.resource_name
  parent_id = azapi_resource.devCenter.id
  body = {
    properties = {
      networkConnectionId = azapi_resource.networkConnection.id
    }
  }
  schema_validation_enabled = false
  response_export_values    = ["*"]
}
