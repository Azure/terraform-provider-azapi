param location string = 'westeurope'
param resource_name string = 'acctest0001'

resource attachNetwork 'Microsoft.DevCenter/devcenters/attachednetworks@2023-04-01' = {
  parent: devCenter
  name: resource_name
  properties: {
    networkConnectionId: networkConnection.id
  }
}

resource devCenter 'Microsoft.DevCenter/devcenters@2023-04-01' = {
  location: location
  name: resource_name
  identity: {
    type: 'SystemAssigned'
    userAssignedIdentities: null
  }
}

resource networkConnection 'Microsoft.DevCenter/networkConnections@2023-04-01' = {
  name: resource_name
  location: location
  properties: {
    domainJoinType: 'AzureADJoin'
    subnetId: subnet.id
  }
}

resource subnet 'Microsoft.Network/virtualNetworks/subnets@2022-07-01' = {
  parent: virtualNetwork
  name: resource_name
  properties: {
    addressPrefix: '10.0.2.0/24'
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
  }
}

