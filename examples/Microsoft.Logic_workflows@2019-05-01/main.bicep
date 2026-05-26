param location string = 'westeurope'
param resource_name string = 'acctest0001'

resource workflow 'Microsoft.Logic/workflows@2019-05-01' = {
  location: location
  name: resource_name
  properties: {
    definition: {
      '$schema': 'https://schema.management.azure.com/providers/Microsoft.Logic/schemas/2016-06-01/workflowdefinition.json#'
      actions: {}
      contentVersion: '1.0.0.0'
      parameters: null
      triggers: {}
    }
    parameters: {}
    state: 'Enabled'
  }
}

