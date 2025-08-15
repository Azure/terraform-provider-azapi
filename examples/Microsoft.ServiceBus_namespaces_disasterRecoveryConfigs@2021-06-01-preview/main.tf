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

variable "secondary_location" {
  type    = string
  default = "centralus"
}

resource "azapi_resource" "resourceGroup" {
  type     = "Microsoft.Resources/resourceGroups@2020-06-01"
  name     = var.resource_name
  location = var.location
}

resource "azapi_resource" "resourceGroup_1" {
  type     = "Microsoft.Resources/resourceGroups@2020-06-01"
  name     = "${var.resource_name}rg2"
  location = var.secondary_location
}

resource "azapi_resource" "namespace" {
  type      = "Microsoft.ServiceBus/namespaces@2022-10-01-preview"
  parent_id = azapi_resource.resourceGroup.id
  name      = "${var.resource_name}ns1"
  location  = var.location
  body = {
    properties = {
      disableLocalAuth           = false
      minimumTlsVersion          = "1.2"
      premiumMessagingPartitions = 1
      publicNetworkAccess        = "Enabled"
    }
    sku = {
      capacity = 1
      name     = "Premium"
      tier     = "Premium"
    }
  }
}

resource "azapi_resource" "namespace_1" {
  type      = "Microsoft.ServiceBus/namespaces@2022-10-01-preview"
  parent_id = azapi_resource.resourceGroup_1.id
  name      = "${var.resource_name}ns2"
  location  = var.secondary_location
  body = {
    properties = {
      disableLocalAuth           = false
      minimumTlsVersion          = "1.2"
      premiumMessagingPartitions = 1
      publicNetworkAccess        = "Enabled"
    }
    sku = {
      capacity = 1
      name     = "Premium"
      tier     = "Premium"
    }
  }
}

resource "azapi_resource" "disasterRecoveryConfig" {
  type      = "Microsoft.ServiceBus/namespaces/disasterRecoveryConfigs@2021-06-01-preview"
  parent_id = azapi_resource.namespace.id
  name      = "${var.resource_name}alias"
  body = {
    properties = {
      partnerNamespace = azapi_resource.namespace_1.id
    }
  }
}

