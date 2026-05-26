param location string = 'westeurope'
param resource_name string = 'acctest0001'

resource emailService 'Microsoft.Communication/emailServices@2023-03-31' = {
  location: 'global'
  name: resource_name
  properties: {
    dataLocation: 'United States'
  }
}

