param location string = 'westus'
param resource_name string = 'acctest0001'

resource fleet 'Microsoft.ContainerService/fleets@2024-04-01' = {
  location: location
  name: resource_name
  properties: {}
}

