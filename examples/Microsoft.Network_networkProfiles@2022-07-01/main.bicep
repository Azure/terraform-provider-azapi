param location string = 'westeurope'
param resource_name string = 'acctest0001'

resource networkProfile 'Microsoft.Network/networkProfiles@2022-07-01' = {
  location: location
  name: resource_name
  properties: {
    containerNetworkInterfaceConfigurations: [
      {
        name: 'acctesteth-230630033653886950'
        properties: {
          ipConfigurations: [
            {
              name: 'acctestipconfig-230630033653886950'
              properties: {
                subnet: {
                  id: subnet.id
                }
              }
            }
          ]
        }
      }
    ]
  }
}

resource subnet 'Microsoft.Network/virtualNetworks/subnets@2022-07-01' = {
  parent: virtualNetwork
  name: resource_name
  properties: {
    addressPrefix: '10.1.0.0/24'
    delegations: [
      {
        name: 'acctestdelegation-230630033653886950'
        properties: {
          serviceName: 'Microsoft.ContainerInstance/containerGroups'
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
        '10.1.0.0/16'
      ]
    }
    dhcpOptions: {
      dnsServers: []
    }
    subnets: []
  }
}

