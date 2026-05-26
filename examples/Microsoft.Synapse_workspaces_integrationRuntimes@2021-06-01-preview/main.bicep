param location string = 'westeurope'
param resource_name string = 'acctest0001'
param sql_administrator_login string = null
param sql_administrator_login_password string = null

resource container 'Microsoft.Storage/storageAccounts/blobServices/containers@2022-09-01' = {
  name: resource_name
  properties: {
    metadata: {
      key: 'value'
    }
  }
}

resource integrationRuntime 'Microsoft.Synapse/workspaces/integrationRuntimes@2021-06-01-preview' = {
  parent: workspace
  name: resource_name
  properties: {
    description: 'test'
    type: 'SelfHosted'
  }
}

resource storageAccount 'Microsoft.Storage/storageAccounts@2021-09-01' = {
  location: location
  name: resource_name
  kind: 'StorageV2'
  properties: {}
  sku: {
    name: 'Standard_LRS'
  }
}

resource workspace 'Microsoft.Synapse/workspaces@2021-06-01' = {
  identity: [
    {
      identity_ids: []
      type: 'SystemAssigned'
    }
  ]
  location: location
  name: resource_name
  properties: {
    defaultDataLakeStorage: {
      accountUrl: storageAccount.properties.primaryEndpoints.dfs
      filesystem: container.name
    }
    managedVirtualNetwork: 'default'
    publicNetworkAccess: 'Enabled'
    sqlAdministratorLogin: sql_administrator_login
    sqlAdministratorLoginPassword: sql_administrator_login_password
  }
}

