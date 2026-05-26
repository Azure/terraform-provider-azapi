param location string = 'centralus'
param resource_name string = 'acctest0001'

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

