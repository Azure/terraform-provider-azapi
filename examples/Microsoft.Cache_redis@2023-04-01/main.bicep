param location string = 'eastus'
param resource_name string = 'acctest0001'

resource redis 'Microsoft.Cache/redis@2023-04-01' = {
  location: location
  name: resource_name
  properties: {
    enableNonSslPort: true
    minimumTlsVersion: '1.2'
    sku: {
      capacity: 2
      family: 'C'
      name: 'Standard'
    }
  }
}

