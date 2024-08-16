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

resource "azapi_resource" "networkSecurityGroup" {
  type      = "Microsoft.Network/networkSecurityGroups@2022-07-01"
  parent_id = azapi_resource.resourceGroup.id
  name      = "mi-security-group1-230630034008554952"
  location  = var.location
  body = {
    properties = {
      securityRules = [
      ]
    }
  }
  schema_validation_enabled = false
  response_export_values    = ["*"]
  lifecycle {
    ignore_changes = [body.properties.securityRules]
  }
}

resource "azapi_resource" "securityRule" {
  type      = "Microsoft.Network/networkSecurityGroups/securityRules@2022-09-01"
  parent_id = azapi_resource.networkSecurityGroup.id
  name      = "allow_management_inbound"
  body = {
    properties = {
      access                   = "Allow"
      destinationAddressPrefix = "*"
      destinationPortRange     = ""
      destinationPortRanges = [
        "9000",
        "1438",
        "1440",
        "9003",
        "1452",
      ]
      direction           = "Inbound"
      priority            = 106
      protocol            = "Tcp"
      sourceAddressPrefix = "*"
      sourcePortRange     = "*"
    }
  }
  schema_validation_enabled = false
  response_export_values    = ["*"]
}

