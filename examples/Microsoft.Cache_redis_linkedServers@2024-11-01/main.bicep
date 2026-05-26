param location string = 'westus'
param resource_name string = 'acctest0001'

resource linkedServer 'Microsoft.Cache/redis/linkedServers@2024-11-01' = {
  parent: redis_primary
  name: '${resource_name}-secondary'
  properties: {
    linkedRedisCacheId: redis_secondary.id
    linkedRedisCacheLocation: location
    serverRole: 'Secondary'
  }
}

resource redis_primary 'Microsoft.Cache/redis@2024-11-01' = {
  parent: resourceGroup_secondary
  location: location
  name: '${resource_name}-primary'
  properties: {
    disableAccessKeyAuthentication: false
    enableNonSslPort: false
    minimumTlsVersion: '1.2'
    publicNetworkAccess: 'Enabled'
    redisConfiguration: {
      'maxmemory-delta': '642'
      'maxmemory-policy': 'allkeys-lru'
      'maxmemory-reserved': '642'
      'preferred-data-persistence-auth-method': ''
    }
    redisVersion: '6'
    sku: {
      capacity: 1
      family: 'P'
      name: 'Premium'
    }
  }
}

resource redis_secondary 'Microsoft.Cache/redis@2024-11-01' = {
  location: location
  name: '${resource_name}-secondary'
  properties: {
    disableAccessKeyAuthentication: false
    enableNonSslPort: false
    minimumTlsVersion: '1.2'
    publicNetworkAccess: 'Enabled'
    redisConfiguration: {
      'maxmemory-delta': '642'
      'maxmemory-policy': 'allkeys-lru'
      'maxmemory-reserved': '642'
      'preferred-data-persistence-auth-method': ''
    }
    redisVersion: '6'
    sku: {
      capacity: 1
      family: 'P'
      name: 'Premium'
    }
  }
}

