---
subcategory: "Microsoft.Network - Various Networking Services"
page_title: "applicationGateways"
description: |-
  Manages a Application Gateway.
---

# Microsoft.Network/applicationGateways - Application Gateway

This article demonstrates how to use `azapi` provider to manage the Application Gateway resource in Azure.

## Example Usage

### default

```hcl
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


```



## Arguments Reference

The following arguments are supported:

* `type` - (Required) The type of the resource. This should be set to `Microsoft.Network/applicationGateways@api-version`. The available api-versions for this resource are: [`2015-05-01-preview`, `2015-06-15`, `2016-03-30`, `2016-06-01`, `2016-09-01`, `2016-12-01`, `2017-03-01`, `2017-03-30`, `2017-06-01`, `2017-08-01`, `2017-09-01`, `2017-10-01`, `2017-11-01`, `2018-01-01`, `2018-02-01`, `2018-04-01`, `2018-06-01`, `2018-07-01`, `2018-08-01`, `2018-10-01`, `2018-11-01`, `2018-12-01`, `2019-02-01`, `2019-04-01`, `2019-06-01`, `2019-07-01`, `2019-08-01`, `2019-09-01`, `2019-11-01`, `2019-12-01`, `2020-03-01`, `2020-04-01`, `2020-05-01`, `2020-06-01`, `2020-07-01`, `2020-08-01`, `2020-11-01`, `2021-02-01`, `2021-03-01`, `2021-05-01`, `2021-08-01`, `2022-01-01`, `2022-05-01`, `2022-07-01`, `2022-09-01`, `2022-11-01`, `2023-02-01`, `2023-04-01`, `2023-05-01`, `2023-06-01`, `2023-09-01`, `2023-11-01`, `2024-01-01`, `2024-03-01`, `2024-05-01`, `2024-07-01`, `2024-10-01`].

* `parent_id` - (Required) The ID of the azure resource in which this resource is created. The allowed values are:  
  `/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}`

* `name` - (Required) Specifies the name of the azure resource. Changing this forces a new resource to be created.

* `body` - (Required) Specifies the configuration of the resource. More information about the arguments in `body` can be found in the [Microsoft documentation](https://learn.microsoft.com/en-us/azure/templates/Microsoft.Network/applicationGateways?pivots=deployment-language-terraform).

For other arguments, please refer to the [azapi_resource](https://registry.terraform.io/providers/Azure/azapi/latest/docs/resources/resource) documentation.

## Import

 ```shell
 # Azure resource can be imported using the resource id, e.g.
 terraform import azapi_resource.example /subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Network/applicationGateways/{resourceName}
 
 # It also supports specifying API version by using the resource id with api-version as a query parameter, e.g.
 terraform import azapi_resource.example /subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Network/applicationGateways/{resourceName}?api-version=2024-10-01
 ```
