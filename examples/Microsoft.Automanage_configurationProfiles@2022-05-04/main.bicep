param location string = 'westeurope'
param resource_name string = 'acctest0001'

resource configurationProfile 'Microsoft.Automanage/configurationProfiles@2022-05-04' = {
  location: location
  name: resource_name
  properties: {
    configuration: {}
  }
}

