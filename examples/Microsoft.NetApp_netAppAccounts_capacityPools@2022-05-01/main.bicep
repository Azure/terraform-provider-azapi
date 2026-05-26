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

