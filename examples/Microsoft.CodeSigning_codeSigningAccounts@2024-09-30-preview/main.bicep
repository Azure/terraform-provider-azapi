param location string = 'westus'
param resource_name string = 'acctest0001'

resource codeSigningAccount 'Microsoft.CodeSigning/codeSigningAccounts@2024-09-30-preview' = {
  location: location
  name: resource_name
  properties: {
    sku: {
      name: 'Basic'
    }
  }
}

