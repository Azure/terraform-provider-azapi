param location string = 'westeurope'
param resource_name string = 'acctest0001'

resource applicationGateway 'Microsoft.Network/applicationGateways@2022-07-01' = {
  location: location
  name: resource_name
  properties: {
    authenticationCertificates: []
    backendAddressPools: [
      {
        name: data.azapi_resource_id.backendAddressPool.name
        properties: {
          backendAddresses: []
        }
      }
    ]
    backendHttpSettingsCollection: [
      {
        name: data.azapi_resource_id.backendHttpSettingsCollection.name
        properties: {
          authenticationCertificates: []
          cookieBasedAffinity: 'Disabled'
          path: ''
          pickHostNameFromBackendAddress: false
          port: 80
          protocol: 'Http'
          requestTimeout: 1
          trustedRootCertificates: []
        }
      }
    ]
    customErrorConfigurations: []
    enableHttp2: false
    frontendIPConfigurations: [
      {
        name: data.azapi_resource_id.frontendIPConfiguration.name
        properties: {
          privateIPAllocationMethod: 'Dynamic'
          publicIPAddress: {
            id: publicIPAddress.id
          }
        }
      }
    ]
    frontendPorts: [
      {
        name: data.azapi_resource_id.frontendPort.name
        properties: {
          port: 80
        }
      }
    ]
    gatewayIPConfigurations: [
      {
        name: 'my-gateway-ip-configuration'
        properties: {
          subnet: {
            id: subnet.id
          }
        }
      }
    ]
    httpListeners: [
      {
        name: data.azapi_resource_id.httpListener.name
        properties: {
          customErrorConfigurations: []
          frontendIPConfiguration: {
            id: data.azapi_resource_id.frontendIPConfiguration.id
          }
          frontendPort: {
            id: data.azapi_resource_id.frontendPort.id
          }
          protocol: 'Http'
          requireServerNameIndication: false
        }
      }
    ]
    privateLinkConfigurations: []
    probes: []
    redirectConfigurations: []
    requestRoutingRules: [
      {
        name: '${virtualNetwork.name}-rqrt'
        properties: {
          backendAddressPool: {
            id: data.azapi_resource_id.backendAddressPool.id
          }
          backendHttpSettings: {
            id: data.azapi_resource_id.backendHttpSettingsCollection.id
          }
          httpListener: {
            id: data.azapi_resource_id.httpListener.id
          }
          priority: 10
          ruleType: 'Basic'
        }
      }
    ]
    rewriteRuleSets: []
    sku: {
      capacity: 2
      name: 'Standard_v2'
      tier: 'Standard_v2'
    }
    sslCertificates: []
    sslPolicy: {}
    sslProfiles: []
    trustedClientCertificates: []
    trustedRootCertificates: []
    urlPathMaps: []
  }
}

resource publicIPAddress 'Microsoft.Network/publicIPAddresses@2022-07-01' = {
  location: location
  name: resource_name
  properties: {
    ddosSettings: {
      protectionMode: 'VirtualNetworkInherited'
    }
    idleTimeoutInMinutes: 4
    publicIPAddressVersion: 'IPv4'
    publicIPAllocationMethod: 'Static'
  }
  sku: {
    name: 'Standard'
    tier: 'Regional'
  }
}

resource subnet 'Microsoft.Network/virtualNetworks/subnets@2022-07-01' = {
  parent: virtualNetwork
  name: 'subnet-230630033653837171'
  properties: {
    addressPrefix: '10.0.0.0/24'
    delegations: []
    privateEndpointNetworkPolicies: 'Enabled'
    privateLinkServiceNetworkPolicies: 'Disabled'
    serviceEndpointPolicies: []
    serviceEndpoints: []
  }
}

resource virtualNetwork 'Microsoft.Network/virtualNetworks@2022-07-01' = {
  location: location
  name: resource_name
  properties: {
    addressSpace: {
      addressPrefixes: [
        '10.0.0.0/16'
      ]
    }
    dhcpOptions: {
      dnsServers: []
    }
    subnets: []
  }
}

