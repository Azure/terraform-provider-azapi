param location string = 'westeurope'
param resource_name string = 'acctest0001'

resource workbook 'Microsoft.Insights/workbooks@2022-04-01' = {
  location: location
  name: 'be1ad266-d329-4454-b693-8287e4d3b35d'
  kind: 'shared'
  properties: {
    category: 'workbook'
    displayName: 'acctest-amw-230630032616547405'
    serializedData: '{"fallbackResourceIds":["Azure Monitor"],"isLocked":false,"items":[{"content":{"json":"Test2022"},"name":"text - 0","type":1}],"version":"Notebook/1.0"}'
    sourceId: 'azure monitor'
  }
}

