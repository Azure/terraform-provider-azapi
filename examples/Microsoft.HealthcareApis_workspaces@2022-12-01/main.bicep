param location string = 'westeurope'
param resource_name string = 'acctest0001'

resource workspace 'Microsoft.HealthcareApis/workspaces@2022-12-01' = {
  location: location
  name: resource_name
}

