param location string = 'westeurope'
param resource_name string = 'acctest0001'

resource dnsForwardingRuleset 'Microsoft.Network/dnsForwardingRulesets@2022-07-01' = {
  location: location
  name: resource_name
  properties: {
    dnsResolverOutboundEndpoints: [
      {
        id: outboundEndpoint.id
      }
    ]
  }
}

resource dnsResolver 'Microsoft.Network/dnsResolvers@2022-07-01' = {
  location: location
  name: resource_name
  properties: {
    virtualNetwork: {
      id: virtualNetwork.id
    }
  }
}

resource outboundEndpoint 'Microsoft.Network/dnsResolvers/outboundEndpoints@2022-07-01' = {
  parent: dnsResolver
  location: location
  name: resource_name
  properties: {
    subnet: {
      id: subnet.id
    }
  }
}

resource subnet 'Microsoft.Network/virtualNetworks/subnets@2022-07-01' = {
  parent: virtualNetwork
  name: 'outbounddns'
  properties: {
    addressPrefix: '10.0.0.64/28'
    delegations: [
      {
        name: 'Microsoft.Network.dnsResolvers'
        properties: {
          serviceName: 'Microsoft.Network/dnsResolvers'
        }
      }
    ]
    privateEndpointNetworkPolicies: 'Enabled'
    privateLinkServiceNetworkPolicies: 'Enabled'
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

resource virtualNetworkLink 'Microsoft.Network/dnsForwardingRulesets/virtualNetworkLinks@2022-07-01' = {
  parent: dnsForwardingRuleset
  name: resource_name
  properties: {
    metadata: null
    virtualNetwork: {
      id: virtualNetwork.id
    }
  }
}

