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
  type                      = "Microsoft.Resources/resourceGroups@2020-06-01"
  name                      = var.resource_name
  location                  = var.location
}

resource "azapi_resource" "virtualNetwork" {
  type      = "Microsoft.Network/virtualNetworks@2022-07-01"
  parent_id = azapi_resource.resourceGroup.id
  name      = var.resource_name
  location  = var.location
  body = jsonencode({
    properties = {
      addressSpace = {
        addressPrefixes = [
          "10.7.29.0/29",
        ]
      }
      dhcpOptions = {
        dnsServers = [
        ]
      }
      subnets = [
      ]
    }
  })
  schema_validation_enabled = false
  response_export_values    = ["*"]
  ignore_body_changes            = ["properties.subnets"]
}

resource "azapi_resource" "server" {
  type      = "Microsoft.DBforPostgreSQL/servers@2017-12-01"
  parent_id = azapi_resource.resourceGroup.id
  name      = var.resource_name
  location  = var.location
  body = jsonencode({
    properties = {
      administratorLogin         = "acctestun"
      administratorLoginPassword = "H@Sh1CoR3!"
      createMode                 = "Default"
      infrastructureEncryption   = "Disabled"
      minimalTlsVersion          = "TLS1_2"
      publicNetworkAccess        = "Enabled"
      sslEnforcement             = "Enabled"
      storageProfile = {
        backupRetentionDays = 7
        storageAutogrow     = "Enabled"
        storageMB           = 51200
      }
      version = "9.5"
    }
    sku = {
      capacity = 2
      family   = "Gen5"
      name     = "GP_Gen5_2"
      tier     = "GeneralPurpose"
    }
  })
  schema_validation_enabled = false
  response_export_values    = ["*"]
}

resource "azapi_resource" "subnet" {
  type      = "Microsoft.Network/virtualNetworks/subnets@2022-07-01"
  parent_id = azapi_resource.virtualNetwork.id
  name      = var.resource_name
  body = jsonencode({
    properties = {
      addressPrefix = "10.7.29.0/29"
      delegations = [
      ]
      privateEndpointNetworkPolicies    = "Enabled"
      privateLinkServiceNetworkPolicies = "Enabled"
      serviceEndpointPolicies = [
      ]
      serviceEndpoints = [
        {
          service = "Microsoft.Sql"
        },
      ]
    }
  })
  schema_validation_enabled = false
  response_export_values    = ["*"]
}

resource "azapi_resource" "virtualNetworkRule" {
  type      = "Microsoft.DBforPostgreSQL/servers/virtualNetworkRules@2017-12-01"
  parent_id = azapi_resource.server.id
  name      = var.resource_name
  body = jsonencode({
    properties = {
      ignoreMissingVnetServiceEndpoint = false
      virtualNetworkSubnetId           = azapi_resource.subnet.id
    }
  })
  schema_validation_enabled = false
  response_export_values    = ["*"]
}

