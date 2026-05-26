param location string = 'westus'
param resource_name string = 'acctest0001'

resource attestationProvider 'Microsoft.Attestation/attestationProviders@2020-10-01' = {
  location: location
  name: resource_name
  properties: {}
}

