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

resource "azapi_resource" "publicIPAddress" {
  type      = "Microsoft.Network/publicIPAddresses@2022-07-01"
  parent_id = azapi_resource.resourceGroup.id
  name      = var.resource_name
  location  = var.location
  body = {
    properties = {
      ddosSettings = {
        protectionMode = "VirtualNetworkInherited"
      }
      idleTimeoutInMinutes     = 4
      publicIPAddressVersion   = "IPv4"
      publicIPAllocationMethod = "Static"
    }
    sku = {
      name = "Standard"
      tier = "Regional"
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
          "10.0.0.0/16",
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

resource "azapi_resource" "subnet" {
  type      = "Microsoft.Network/virtualNetworks/subnets@2022-07-01"
  parent_id = azapi_resource.virtualNetwork.id
  name      = "subnet-230630033653837171"
  body = {
    properties = {
      addressPrefix = "10.0.0.0/24"
      delegations = [
      ]
      privateEndpointNetworkPolicies    = "Enabled"
      privateLinkServiceNetworkPolicies = "Disabled"
      serviceEndpointPolicies = [
      ]
      serviceEndpoints = [
      ]
    }
  }
  schema_validation_enabled = false
  response_export_values    = ["*"]
}

data "azapi_resource_id" "applicationGateway" {
  type      = "Microsoft.Network/applicationGateways@2022-07-01"
  parent_id = azapi_resource.resourceGroup.id
  name      = var.resource_name
}

data "azapi_resource_id" "frontendIPConfiguration" {
  type      = "Microsoft.Network/applicationGateways/frontendIPConfigurations@2022-07-01"
  parent_id = data.azapi_resource_id.applicationGateway.id
  name      = "${azapi_resource.virtualNetwork.name}-feip"
}

data "azapi_resource_id" "frontendPort" {
  type      = "Microsoft.Network/applicationGateways/frontendPorts@2022-07-01"
  parent_id = data.azapi_resource_id.applicationGateway.id
  name      = "${azapi_resource.virtualNetwork.name}-feport"
}

data "azapi_resource_id" "backendAddressPool" {
  type      = "Microsoft.Network/applicationGateways/backendAddressPools@2022-07-01"
  parent_id = data.azapi_resource_id.applicationGateway.id
  name      = "${azapi_resource.virtualNetwork.name}-beap"
}

data "azapi_resource_id" "backendHttpSettingsCollection" {
  type      = "Microsoft.Network/applicationGateways/backendHttpSettingsCollection@2022-07-01"
  parent_id = data.azapi_resource_id.applicationGateway.id
  name      = "${azapi_resource.virtualNetwork.name}-be-htst"
}

data "azapi_resource_id" "httpListener" {
  type      = "Microsoft.Network/applicationGateways/httpListeners@2022-07-01"
  parent_id = data.azapi_resource_id.applicationGateway.id
  name      = "${azapi_resource.virtualNetwork.name}-httplstn"
}

resource "azapi_resource" "applicationGateway" {
  type      = "Microsoft.Network/applicationGateways@2022-07-01"
  parent_id = azapi_resource.resourceGroup.id
  name      = var.resource_name
  location  = var.location
  body = {
    properties = {
      authenticationCertificates = [
      ]
      backendAddressPools = [
        {
          name = data.azapi_resource_id.backendAddressPool.name
          properties = {
            backendAddresses = [
            ]
          }
        },
      ]
      backendHttpSettingsCollection = [
        {
          name = data.azapi_resource_id.backendHttpSettingsCollection.name
          properties = {
            authenticationCertificates = [
            ]
            cookieBasedAffinity            = "Disabled"
            path                           = ""
            pickHostNameFromBackendAddress = false
            port                           = 80
            protocol                       = "Http"
            requestTimeout                 = 1
            trustedRootCertificates = [
            ]
          }
        },
      ]
      customErrorConfigurations = [
      ]
      enableHttp2 = false
      frontendIPConfigurations = [
        {
          name = data.azapi_resource_id.frontendIPConfiguration.name
          properties = {
            privateIPAllocationMethod = "Dynamic"
            publicIPAddress = {
              id = azapi_resource.publicIPAddress.id
            }
          }
        },
      ]
      frontendPorts = [
        {
          name = data.azapi_resource_id.frontendPort.name
          properties = {
            port = 80
          }
        },
      ]
      gatewayIPConfigurations = [
        {
          name = "my-gateway-ip-configuration"
          properties = {
            subnet = {
              id = azapi_resource.subnet.id
            }
          }
        },
      ]
      httpListeners = [
        {
          name = data.azapi_resource_id.httpListener.name
          properties = {
            customErrorConfigurations = [
            ]
            frontendIPConfiguration = {
              id = data.azapi_resource_id.frontendIPConfiguration.id
            }
            frontendPort = {
              id = data.azapi_resource_id.frontendPort.id
            }
            protocol                    = "Http"
            requireServerNameIndication = false
          }
        },
      ]
      privateLinkConfigurations = [
      ]
      probes = [
      ]
      redirectConfigurations = [
      ]
      requestRoutingRules = [
        {
          name = "${azapi_resource.virtualNetwork.name}-rqrt"
          properties = {
            backendAddressPool = {
              id = data.azapi_resource_id.backendAddressPool.id
            }
            backendHttpSettings = {
              id = data.azapi_resource_id.backendHttpSettingsCollection.id
            }
            httpListener = {
              id = data.azapi_resource_id.httpListener.id
            }
            ruleType = "Basic"
            priority = 10
          }
        },
      ]
      rewriteRuleSets = [
      ]
      sku = {
        capacity = 2
        name     = "Standard_v2"
        tier     = "Standard_v2"
      }
      sslCertificates = [
      ]
      sslPolicy = {
      }
      sslProfiles = [
      ]
      trustedClientCertificates = [
      ]
      trustedRootCertificates = [
      ]
      urlPathMaps = [
      ]
    }
  }
  schema_validation_enabled = false
  response_export_values    = ["*"]
}

