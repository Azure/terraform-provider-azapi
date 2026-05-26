param location string = 'eastus'
param resource_name string = 'acctest0001'

resource scheduledAction 'Microsoft.CostManagement/scheduledActions@2022-06-01-preview' = {
  name: resource_name
  kind: 'InsightAlert'
  properties: {
    displayName: 'acctest 230630032939736168'
    fileDestination: {
      fileFormats: []
    }
    notification: {
      message: 'Oops, cost anomaly'
      subject: 'Hi'
      to: [
        'test@test.com'
        'test@hashicorp.developer'
      ]
    }
    schedule: {
      endDate: '2024-07-18T11:21:14+00:00'
      frequency: 'Daily'
      startDate: '2023-07-18T11:21:14+00:00'
    }
    status: 'Enabled'
    viewId: view.id
  }
}

resource view 'Microsoft.CostManagement/views@2022-10-01' = {
  name: resource_name
  properties: {
    accumulated: 'False'
    chart: 'StackedColumn'
    displayName: 'Test View wgvtl'
    kpis: [
      {
        enabled: true
        type: 'Forecast'
      }
    ]
    pivots: [
      {
        name: 'ServiceName'
        type: 'Dimension'
      }
      {
        name: 'ResourceLocation'
        type: 'Dimension'
      }
      {
        name: 'ResourceGroupName'
        type: 'Dimension'
      }
    ]
    query: {
      dataSet: {
        aggregation: {
          totalCost: {
            function: 'Sum'
            name: 'Cost'
          }
          totalCostUSD: {
            function: 'Sum'
            name: 'CostUSD'
          }
        }
        granularity: 'Monthly'
        grouping: [
          {
            name: 'ResourceGroupName'
            type: 'Dimension'
          }
        ]
        sorting: [
          {
            direction: 'Ascending'
            name: 'BillingMonth'
          }
        ]
      }
      timeframe: 'MonthToDate'
      type: 'Usage'
    }
  }
}

