param location string = 'westeurope'
param resource_name string = 'acctest0001'

resource diskAccess 'Microsoft.Compute/diskAccesses@2022-03-02' = {
  location: location
  name: resource_name
  tags: {
    'cost-center': 'ops'
    environment: 'acctest'
  }
}

