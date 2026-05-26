param location string = 'westeurope'
param resource_name string = 'acctest0001'

resource query 'Microsoft.OperationalInsights/queryPacks/queries@2019-09-01' = {
  parent: queryPack
  name: 'aca50e92-d3e6-8f7d-1f70-2ec7adc1a926'
  properties: {
    body: '    let newExceptionsTimeRange = 1d;\n    let timeRangeToCheckBefore = 7d;\n    exceptions\n    | where timestamp < ago(timeRangeToCheckBefore)\n    | summarize count() by problemId\n    | join kind= rightanti (\n        exceptions\n        | where timestamp >= ago(newExceptionsTimeRange)\n        | extend stack = tostring(details[0].rawStack)\n        | summarize count(), dcount(user_AuthenticatedId), min(timestamp), max(timestamp), any(stack) by problemId\n    ) on problemId\n    | order by count_ desc\n'
    displayName: 'Exceptions - New in the last 24 hours'
    related: {}
  }
}

resource queryPack 'Microsoft.OperationalInsights/queryPacks@2019-09-01' = {
  location: location
  name: resource_name
  properties: {}
}

