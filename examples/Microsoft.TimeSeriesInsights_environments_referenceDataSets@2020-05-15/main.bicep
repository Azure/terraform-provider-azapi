param location string = 'westeurope'
param resource_name string = 'acctest0001'

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

resource referenceDataSet 'Microsoft.TimeSeriesInsights/environments/referenceDataSets@2020-05-15' = {
  parent: environment
  location: location
  name: resource_name
  properties: {
    dataStringComparisonBehavior: 'Ordinal'
    keyProperties: [
      {
        name: 'keyProperty1'
        type: 'String'
      }
    ]
  }
}

