param location string = 'westeurope'
param resource_name string = 'acctest0001'

resource localNetworkGateway 'Microsoft.Network/localNetworkGateways@2022-07-01' = {
  location: location
  name: resource_name
  properties: {
    gatewayIpAddress: '168.62.225.23'
    localNetworkAddressSpace: {
      addressPrefixes: [
        '10.1.1.0/24'
      ]
    }
  }
}

