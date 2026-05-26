param location string = 'westus'
param resource_name string = 'acctest0001'

resource firewallRule 'Microsoft.Cache/redis/firewallRules@2024-11-01' = {
  parent: redis
  name: '${resource_name}_fwrule'
  properties: {
    endIP: '2.3.4.5'
    startIP: '1.2.3.4'
  }
}

resource redis 'Microsoft.Cache/redis@2024-11-01' = {
  location: location
  name: resource_name
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
    redisVersion: '6.0'
    sku: {
      capacity: 1
      family: 'P'
      name: 'Premium'
    }
  }
}

