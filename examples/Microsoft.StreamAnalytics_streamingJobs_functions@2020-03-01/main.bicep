param location string = 'westeurope'
param resource_name string = 'acctest0001'

resource function 'Microsoft.StreamAnalytics/streamingJobs/functions@2020-03-01' = {
  parent: streamingJob
  name: resource_name
  properties: {
    properties: {
      binding: {
        properties: {
          script: 'function getRandomNumber(in) {\n  return in;\n}\n'
        }
        type: 'Microsoft.StreamAnalytics/JavascriptUdf'
      }
      inputs: [
        {
          dataType: 'bigint'
          isConfigurationParameter: false
        }
      ]
      output: {
        dataType: 'bigint'
      }
    }
    type: 'Scalar'
  }
}

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
        query: '    SELECT *\n    INTO [YourOutputAlias]\n    FROM [YourInputAlias]\n'
        streamingUnits: 3
      }
    }
  }
}

