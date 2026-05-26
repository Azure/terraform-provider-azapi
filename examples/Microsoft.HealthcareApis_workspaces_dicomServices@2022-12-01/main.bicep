param location string = 'westeurope'
param resource_name string = 'acctest0001'

resource dicomService 'Microsoft.HealthcareApis/workspaces/dicomServices@2022-12-01' = {
  parent: workspace
  location: location
  name: resource_name
  properties: {
    publicNetworkAccess: 'Enabled'
  }
}

resource workspace 'Microsoft.HealthcareApis/workspaces@2022-12-01' = {
  location: location
  name: resource_name
}

