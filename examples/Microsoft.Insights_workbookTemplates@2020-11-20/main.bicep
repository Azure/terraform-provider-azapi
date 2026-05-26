param location string = 'westeurope'
param resource_name string = 'acctest0001'

resource workbookTemplate 'Microsoft.Insights/workbookTemplates@2020-11-20' = {
  location: location
  name: resource_name
  properties: {
    galleries: [
      {
        category: 'workbook'
        name: 'test'
        order: 0
        resourceType: 'Azure Monitor'
        type: 'workbook'
      }
    ]
    priority: 0
    templateData: {
      '$schema': 'https://github.com/Microsoft/Application-Insights-Workbooks/blob/master/schema/workbook.json'
      items: [
        {
          content: {
            json: '## New workbook\n---\n\nWelcome to your new workbook.'
          }
          name: 'text - 2'
          type: 1
        }
      ]
      styleSettings: {}
      version: 'Notebook/1.0'
    }
  }
}

