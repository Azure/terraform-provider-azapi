param location string = 'eastus'
param resource_name string = 'acctest0001'

resource roleDefinition 'Microsoft.Authorization/roleDefinitions@2018-01-01-preview' = {
  name: '6faae21a-0cd6-4536-8c23-a278823d12ed'
  properties: {
    assignableScopes: [
      data.azapi_resource.subscription.id
    ]
    description: ''
    permissions: [
      {
        actions: [
          '*'
        ]
        dataActions: []
        notActions: []
        notDataActions: []
      }
    ]
    roleName: resource_name
    type: 'CustomRole'
  }
}

