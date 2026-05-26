param location string = 'westeurope'
param resource_name string = 'acctest0001'

resource application 'Microsoft.Compute/galleries/applications@2022-03-03' = {
  parent: gallery
  location: location
  name: resource_name
  properties: {
    supportedOSType: 'Linux'
  }
}

resource gallery 'Microsoft.Compute/galleries@2022-03-03' = {
  location: location
  name: resource_name
  properties: {
    description: ''
  }
}

