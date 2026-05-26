param location string = 'westeurope'
param resource_name string = 'acctest0001'

resource lab 'Microsoft.DevTestLab/labs@2018-09-15' = {
  location: location
  name: resource_name
  properties: {
    labStorageType: 'Premium'
  }
}

