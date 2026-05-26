param location string = 'westeurope'
param qumulo_password string = ')^X#ZX#JRyIY}t9'
param resource_name string = 'acctest0001'

resource qumuloFileSystem 'Qumulo.Storage/fileSystems@2024-06-19' = {
  location: location
  name: resource_name
  tags: {
    environment: 'terraform-acctests'
    some_key: 'some-value'
  }
  properties: {
    adminPassword: qumulo_password
    availabilityZone: '1'
    delegatedSubnetId: subnet.id
    marketplaceDetails: {
      offerId: 'qumulo-saas-mpp'
      planId: 'azure-native-qumulo-v3'
      publisherId: 'qumulo1584033880660'
    }
    storageSku: 'Cold_LRS'
    userDetails: {
      email: 'test@test.com'
    }
  }
}

resource subnet 'Microsoft.Network/virtualNetworks/subnets@2024-05-01' = {
  parent: vnet
  location: location
  name: resource_name
  properties: {
    addressPrefix: '10.0.1.0/24'
    defaultOutboundAccess: true
    delegations: [
      {
        name: 'delegation'
        properties: {
          actions: [
            'Microsoft.Network/virtualNetworks/subnets/join/action'
          ]
          serviceName: 'Qumulo.Storage/fileSystems'
        }
      }
    ]
    privateEndpointNetworkPolicies: 'Disabled'
    privateLinkServiceNetworkPolicies: 'Enabled'
  }
}

resource vnet 'Microsoft.Network/virtualNetworks@2024-05-01' = {
  location: location
  name: resource_name
  properties: {
    addressSpace: {
      addressPrefixes: [
        '10.0.0.0/16'
      ]
    }
    privateEndpointVNetPolicies: 'Disabled'
    subnets: []
  }
}

