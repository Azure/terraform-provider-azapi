param location string = 'westeurope'
param resource_name string = 'acctest0001'

resource accessPolicy 'Microsoft.TimeSeriesInsights/environments/accessPolicies@2020-05-15' = {
  parent: environment
  name: resource_name
  properties: {
    description: ''
    principalObjectId: 'aGUID'
    roles: [
      'Reader'
    ]
  }
}

resource environment 'Microsoft.TimeSeriesInsights/environments@2020-05-15' = {
  location: location
  name: resource_name
  kind: 'Gen1'
  properties: {
    dataRetentionTime: 'P30D'
    storageLimitExceededBehavior: 'PurgeOldData'
  }
  sku: {
    capacity: 1
    name: 'S1'
  }
}

