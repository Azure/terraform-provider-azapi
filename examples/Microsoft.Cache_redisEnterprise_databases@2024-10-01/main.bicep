param location string = 'westus'
param resource_name string = 'acctest0001'

resource databas 'Microsoft.Cache/redisEnterprise/databases@2024-10-01' = {
  parent: redisEnterprise
  name: 'default'
  properties: {
    clientProtocol: 'Encrypted'
    clusteringPolicy: 'OSSCluster'
    evictionPolicy: 'VolatileLRU'
    modules: []
    port: 10000
  }
}

resource redisEnterprise 'Microsoft.Cache/redisEnterprise@2024-10-01' = {
  location: location
  name: resource_name
  properties: {
    minimumTlsVersion: '1.2'
  }
  sku: {
    capacity: 4
    name: 'Enterprise_E20'
  }
}

