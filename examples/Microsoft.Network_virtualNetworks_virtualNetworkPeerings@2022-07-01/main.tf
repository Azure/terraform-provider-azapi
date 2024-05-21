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

data "azapi_resource_id" "workspace_resource_group" {
  type      = "Microsoft.Resources/resourceGroups@2020-06-01"
  parent_id = azapi_resource.resourceGroup.parent_id
  name      = "databricks-rg-${var.resource_name}"
}

resource "azapi_resource" "workspace" {
  type      = "Microsoft.Databricks/workspaces@2023-02-01"
  parent_id = azapi_resource.resourceGroup.id
  name      = var.resource_name
  location  = var.location
  body = {
    properties = {
      managedResourceGroupId = data.azapi_resource_id.workspace_resource_group.id
      parameters = {
        prepareEncryption = {
          value = false
        }
        requireInfrastructureEncryption = {
          value = false
        }
      }
      publicNetworkAccess = "Enabled"
    }
    sku = {
      name = "standard"
    }
  }
  schema_validation_enabled = false
  response_export_values    = ["*"]
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
          "10.0.1.0/24",
        ]
      }
      dhcpOptions = {
        dnsServers = [
        ]
      }
      subnets = [
      ]
    }
  }
  schema_validation_enabled = false
  response_export_values    = ["*"]
  lifecycle {
    ignore_changes = [body.properties.subnets]
  }
}

resource "azapi_resource" "virtualNetworkPeering" {
  type      = "Microsoft.Databricks/workspaces/virtualNetworkPeerings@2023-02-01"
  parent_id = azapi_resource.workspace.id
  name      = var.resource_name
  body = {
    properties = {
      allowForwardedTraffic     = false
      allowGatewayTransit       = false
      allowVirtualNetworkAccess = true
      databricksAddressSpace = {
        addressPrefixes = [
          "10.139.0.0/16"
        ]
      }
      remoteAddressSpace = {
        addressPrefixes = [
          "10.0.1.0/24",
        ]
      }
      remoteVirtualNetwork = {
        id = azapi_resource.virtualNetwork.id
      }
      useRemoteGateways = false
    }
  }
  schema_validation_enabled = false
  response_export_values    = ["*"]
}


