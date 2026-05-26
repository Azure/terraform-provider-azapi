param location string = 'westus'
param resource_name string = 'acctest0001'

resource budget 'Microsoft.Consumption/budgets@2019-10-01' = {
  name: resource_name
  properties: {
    amount: 1000
    category: 'Cost'
    filter: {
      tags: {
        name: 'foo'
        operator: 'In'
        values: [
          'bar'
        ]
      }
    }
    notifications: {
      'Actual_EqualTo_90.000000_Percent': {
        contactEmails: [
          'foo@example.com'
          'bar@example.com'
        ]
        contactGroups: []
        contactRoles: []
        enabled: true
        operator: 'EqualTo'
        threshold: 90
        thresholdType: 'Actual'
      }
    }
    timeGrain: 'Monthly'
    timePeriod: {
      startDate: '2025-08-01T00:00:00Z'
    }
  }
}

