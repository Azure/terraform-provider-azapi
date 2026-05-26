param location string = 'westeurope'
param resource_name string = 'acctest0001'

resource capacityPool 'Microsoft.NetApp/netAppAccounts/capacityPools@2022-05-01' = {
  parent: netAppAccount
  location: location
  name: resource_name
  properties: {
    serviceLevel: 'Premium'
    size: 4398046511104
  }
}

resource netAppAccount 'Microsoft.NetApp/netAppAccounts@2022-05-01' = {
  location: location
  name: resource_name
  properties: {
    activeDirectories: []
  }
}

resource snapshot 'Microsoft.NetApp/netAppAccounts/capacityPools/volumes/snapshots@2022-05-01' = {
  parent: volume
  location: location
  name: resource_name
}

resource subnet 'Microsoft.Network/virtualNetworks/subnets@2022-07-01' = {
  parent: virtualNetwork
  name: resource_name
  properties: {
    addressPrefix: '10.0.2.0/24'
    delegations: [
      {
        name: 'netapp'
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
        '10.0.0.0/16'
      ]
    }
    dhcpOptions: {
      dnsServers: []
    }
    subnets: []
  }
}

resource volume 'Microsoft.NetApp/netAppAccounts/capacityPools/volumes@2022-05-01' = {
  parent: capacityPool
  location: location
  name: resource_name
  properties: {
    avsDataStore: 'Disabled'
    creationToken: 'my-unique-file-path-230630033642692134'
    dataProtection: {}
    exportPolicy: {
      rules: []
    }
    networkFeatures: 'Basic'
    protocolTypes: [
      'NFSv3'
    ]
    securityStyle: 'Unix'
    serviceLevel: 'Premium'
    snapshotDirectoryVisible: false
    snapshotId: ''
    subnetId: subnet.id
    usageThreshold: 107374182400
    volumeType: ''
  }
  zones: []
}

