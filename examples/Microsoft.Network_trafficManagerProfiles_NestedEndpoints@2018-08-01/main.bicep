param location string = 'westeurope'
param resource_name string = 'acctest0001'

resource NestedEndpoint 'Microsoft.Network/trafficManagerProfiles/NestedEndpoints@2018-08-01' = {
  parent: trafficManagerProfile
  name: resource_name
  properties: {
    customHeaders: []
    endpointStatus: 'Enabled'
    minChildEndpoints: 5
    subnets: []
    targetResourceId: trafficManagerProfile2.id
    weight: 3
  }
}

resource trafficManagerProfile 'Microsoft.Network/trafficManagerProfiles@2018-08-01' = {
  location: 'global'
  name: resource_name
  properties: {
    dnsConfig: {
      relativeName: 'acctest-tmp-230630034107605443'
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

resource trafficManagerProfile2 'Microsoft.Network/trafficManagerProfiles@2018-08-01' = {
  location: 'global'
  name: resource_name
  properties: {
    dnsConfig: {
      relativeName: 'acctesttmpchild230630034107605443'
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
    trafficRoutingMethod: 'Priority'
  }
}

