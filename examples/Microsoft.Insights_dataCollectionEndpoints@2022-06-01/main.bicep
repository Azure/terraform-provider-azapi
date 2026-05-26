param location string = 'westeurope'
param resource_name string = 'acctest0001'

resource dataCollectionEndpoint 'Microsoft.Insights/dataCollectionEndpoints@2022-06-01' = {
  location: location
  name: resource_name
  properties: {
    description: ''
    networkAcls: {
      publicNetworkAccess: 'Enabled'
    }
  }
}

