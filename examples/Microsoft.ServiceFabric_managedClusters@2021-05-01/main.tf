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

resource "azapi_resource" "managedCluster" {
  type      = "Microsoft.ServiceFabric/managedClusters@2021-05-01"
  parent_id = azapi_resource.resourceGroup.id
  name      = var.resource_name
  location  = var.location
  body = {
    properties = {
      addonFeatures = [
        "DnsService",
      ]
      adminPassword             = "NotV3ryS3cur3P@$$w0rd"
      adminUserName             = "testUser"
      clientConnectionPort      = 12345
      clusterUpgradeCadence     = "Wave0"
      dnsName                   = var.resource_name
      httpGatewayConnectionPort = 23456
      loadBalancingRules = [
        {
          backendPort      = 8000
          frontendPort     = 443
          probeProtocol    = "http"
          probeRequestPath = "/"
          protocol         = "tcp"
        },
      ]
      networkSecurityRules = [
        {
          access = "allow"
          destinationAddressPrefixes = [
            "0.0.0.0/0",
          ]
          destinationPortRanges = [
            "443",
          ]
          direction = "inbound"
          name      = "rule443-allow-fe"
          priority  = 1000
          protocol  = "tcp"
          sourceAddressPrefixes = [
            "0.0.0.0/0",
          ]
          sourcePortRanges = [
            "1-65535",
          ]
        },
      ]
    }
    sku = {
      name = "Standard"
    }
    tags = {
      Test = "value"
    }
  }
  schema_validation_enabled = false
  response_export_values    = ["*"]
}

