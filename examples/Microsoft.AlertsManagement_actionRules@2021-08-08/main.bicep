param location string = 'westeurope'
param resource_name string = 'acctest0001'

resource actionRule 'Microsoft.AlertsManagement/actionRules@2021-08-08' = {
  location: 'global'
  name: resource_name
  properties: {
    actions: [
      {
        actionType: 'RemoveAllActionGroups'
      }
    ]
    description: ''
    enabled: true
    scopes: [
      resourceGroup.id
    ]
  }
}

