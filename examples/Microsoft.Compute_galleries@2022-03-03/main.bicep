param location string = 'westeurope'
param resource_name string = 'acctest0001'

resource gallery 'Microsoft.Compute/galleries@2022-03-03' = {
  location: location
  name: resource_name
  properties: {
    description: ''
  }
}

