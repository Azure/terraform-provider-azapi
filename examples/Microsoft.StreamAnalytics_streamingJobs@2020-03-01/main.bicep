param location string = 'westeurope'
param resource_name string = 'acctest0001'

resource streamingJob 'Microsoft.StreamAnalytics/streamingJobs@2020-03-01' = {
  location: location
  name: resource_name
  properties: {
    cluster: {}
    compatibilityLevel: '1.0'
    contentStoragePolicy: 'SystemAccount'
    dataLocale: 'en-GB'
    eventsLateArrivalMaxDelayInSeconds: 60
    eventsOutOfOrderMaxDelayInSeconds: 50
    eventsOutOfOrderPolicy: 'Adjust'
    jobType: 'Cloud'
    outputErrorPolicy: 'Drop'
    sku: {
      name: 'Standard'
    }
    transformation: {
      name: 'main'
      properties: {
        query: '   SELECT *\n   INTO [YourOutputAlias]\n   FROM [YourInputAlias]\n'
        streamingUnits: 3
      }
    }
  }
}

