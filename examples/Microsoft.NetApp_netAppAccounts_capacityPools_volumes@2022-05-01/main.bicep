param location string = 'centralus'
param resource_name string = 'acctest0001'

resource capacityPool 'Microsoft.NetApp/netAppAccounts/capacityPools@2022-05-01' = {
  parent: netAppAccount
  location: location
  name: resource_name
  properties: {
    serviceLevel: 'Standard'
    size: 4398046511104
  }
  tags: {
    SkipASMAzSecPack: 'true'
  }
}

resource netAppAccount 'Microsoft.NetApp/netAppAccounts@2022-05-01' = {
  location: location
  name: resource_name
  properties: {
    activeDirectories: []
  }
  tags: {
    SkipASMAzSecPack: 'true'
  }
}

resource subnet 'Microsoft.Network/virtualNetworks/subnets@2022-07-01' = {
  parent: virtualNetwork
  name: 'GatewaySubnet'
  properties: {
    addressPrefix: '10.6.1.0/24'
    delegations: []
    privateEndpointNetworkPolicies: 'Enabled'
    privateLinkServiceNetworkPolicies: 'Enabled'
    serviceEndpointPolicies: []
    serviceEndpoints: []
  }
}

resource subnet2 'Microsoft.Network/virtualNetworks/subnets@2022-07-01' = {
  parent: virtualNetwork
  name: resource_name
  properties: {
    addressPrefix: '10.6.2.0/24'
    delegations: [
      {
        name: 'testdelegation'
        properties: {
          serviceName: 'Microsoft.Netapp/volumes'
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
        '10.6.0.0/16'
      ]
    }
    dhcpOptions: {
      dnsServers: []
    }
    subnets: []
  }
  tags: {
    SkipASMAzSecPack: 'true'
  }
}

resource volume 'Microsoft.NetApp/netAppAccounts/capacityPools/volumes@2022-05-01' = {
  parent: capacityPool
  location: location
  name: resource_name
  properties: {
    avsDataStore: 'Enabled'
    creationToken: 'my-unique-file-path-230630034120103726'
    dataProtection: {}
    exportPolicy: {
      rules: [
        {
          allowedClients: '0.0.0.0/0'
          cifs: false
          hasRootAccess: true
          nfsv3: true
          nfsv41: false
          ruleIndex: 1
          unixReadOnly: false
          unixReadWrite: true
        }
      ]
    }
    networkFeatures: 'Basic'
    protocolTypes: [
      'NFSv3'
    ]
    serviceLevel: 'Standard'
    snapshotDirectoryVisible: true
    subnetId: subnet2.id
    usageThreshold: 107374182400
    volumeType: ''
  }
  tags: {
    SkipASMAzSecPack: 'true'
  }
  zones: []
}

