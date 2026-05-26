param location string = 'westeurope'
param resource_name string = 'acctest0001'

resource privateLinkHub 'Microsoft.Synapse/privateLinkHubs@2021-06-01' = {
  location: location
  name: resource_name
}

