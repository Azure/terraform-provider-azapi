param location string = 'westeurope'
param resource_name string = 'acctest0001'

resource ExternalEndpoint 'Microsoft.Network/trafficManagerProfiles/ExternalEndpoints@2018-08-01' = {
  parent: trafficManagerProfile
  name: resource_name
  properties: {
    customHeaders: []
    endpointStatus: 'Enabled'
    subnets: []
    target: 'www.example.com'
    weight: 3
  }
}

resource trafficManagerProfile 'Microsoft.Network/trafficManagerProfiles@2018-08-01' = {
  location: 'global'
  name: resource_name
  properties: {
    dnsConfig: {
      relativeName: 'acctest-tmp-230630034107608613'
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

