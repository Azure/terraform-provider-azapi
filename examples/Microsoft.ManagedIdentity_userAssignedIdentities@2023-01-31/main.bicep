param location string = 'westeurope'
param resource_name string = 'acctest0001'

resource userAssignedIdentity 'Microsoft.ManagedIdentity/userAssignedIdentities@2023-01-31' = {
  location: location
  name: resource_name
}

