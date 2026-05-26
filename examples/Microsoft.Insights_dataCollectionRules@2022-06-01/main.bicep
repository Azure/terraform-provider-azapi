param location string = 'westeurope'
param resource_name string = 'acctest0001'

resource dataCollectionRule 'Microsoft.Insights/dataCollectionRules@2022-06-01' = {
  location: location
  name: resource_name
  properties: {
    dataFlows: [
      {
        destinations: [
          'test-destination-metrics'
        ]
        streams: [
          'Microsoft-InsightsMetrics'
        ]
      }
    ]
    description: ''
    destinations: {
      azureMonitorMetrics: {
        name: 'test-destination-metrics'
      }
    }
  }
}

