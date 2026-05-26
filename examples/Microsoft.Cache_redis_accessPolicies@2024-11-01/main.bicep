param location string = 'westus'
param resource_name string = 'acctest0001'

resource accessPolicy 'Microsoft.Cache/redis/accessPolicies@2024-11-01' = {
  parent: redis
  name: '${resource_name}-accessPolicy'
  properties: {
    permissions: '+@read +@connection +cluster|info allkeys'
  }
}

resource redis 'Microsoft.Cache/redis@2024-11-01' = {
  location: location
  name: resource_name
  properties: {
    disableAccessKeyAuthentication: false
    enableNonSslPort: true
    minimumTlsVersion: '1.2'
    publicNetworkAccess: 'Enabled'
    redisConfiguration: {
      'maxmemory-policy': 'volatile-lru'
      'preferred-data-persistence-auth-method': ''
    }
    redisVersion: '6'
    sku: {
      capacity: 1
      family: 'C'
      name: 'Basic'
    }
  }
}

