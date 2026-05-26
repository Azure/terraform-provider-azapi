param location string = 'westeurope'
param resource_name string = 'acctest0001'

resource digitalTwinsInstance 'Microsoft.DigitalTwins/digitalTwinsInstances@2020-12-01' = {
  identity: [
    {
      identity_ids: []
      type: 'SystemAssigned'
    }
  ]
  location: location
  name: resource_name
}

