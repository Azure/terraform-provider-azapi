param location string = 'eastus'
param resource_name string = 'acctest0001'

resource redisEnterprise 'Microsoft.Cache/redisEnterprise@2022-01-01' = {
  location: location
  name: resource_name
  properties: {
    minimumTlsVersion: '1.2'
  }
  sku: {
    capacity: 2
    name: 'Enterprise_E100'
  }
}

