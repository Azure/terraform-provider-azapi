param location string = 'westeurope'
param resource_name string = 'acctest0001'

resource AzureEndpoint 'Microsoft.Network/trafficManagerProfiles/AzureEndpoints@2018-08-01' = {
  parent: trafficManagerProfile
  name: resource_name
  properties: {
    customHeaders: []
    endpointStatus: 'Enabled'
    subnets: []
    targetResourceId: publicIPAddress.id
    weight: 3
  }
}

resource publicIPAddress 'Microsoft.Network/publicIPAddresses@2022-07-01' = {
  location: location
  name: resource_name
  properties: {
    ddosSettings: {
      protectionMode: 'VirtualNetworkInherited'
    }
    dnsSettings: {
      domainNameLabel: 'acctestpublicip-230630034107607730'
    }
    idleTimeoutInMinutes: 4
    publicIPAddressVersion: 'IPv4'
    publicIPAllocationMethod: 'Static'
  }
  sku: {
    name: 'Basic'
    tier: 'Regional'
  }
}

resource trafficManagerProfile 'Microsoft.Network/trafficManagerProfiles@2018-08-01' = {
  location: 'global'
  name: resource_name
  properties: {
    dnsConfig: {
      relativeName: 'acctest-tmp-230630034107607730'
      ttl: 30
    }
    monitorConfig: {
      expectedStatusCodeRanges: []
      intervalInSeconds: 30
      path: '/'
      port: 443
      protocol: 'HTTPS'
      timeoutInSeconds: 10
      toleratedNumberOfFailures: 3
    }
    trafficRoutingMethod: 'Weighted'
  }
}

