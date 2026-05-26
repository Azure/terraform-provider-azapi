param location string = 'westeurope'
param resource_name string = 'acctest0001'

resource federatedIdentityCredential 'Microsoft.ManagedIdentity/userAssignedIdentities/federatedIdentityCredentials@2022-01-31-preview' = {
  parent: userAssignedIdentity
  location: location
  name: resource_name
  properties: {
    audiences: [
      'foo'
    ]
    issuer: 'https://foo'
    subject: 'foo'
  }
}

resource userAssignedIdentity 'Microsoft.ManagedIdentity/userAssignedIdentities@2023-01-31' = {
  location: location
  name: resource_name
}

