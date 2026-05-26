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

resource firewallRule 'Microsoft.Synapse/workspaces/firewallRules@2021-06-01' = {
  parent: workspace
  name: 'AllowAll'
  properties: {
    endIpAddress: '255.255.255.255'
    startIpAddress: '0.0.0.0'
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
    managedVirtualNetwork: ''
    publicNetworkAccess: 'Enabled'
    sqlAdministratorLogin: sql_administrator_login
    sqlAdministratorLoginPassword: sql_administrator_login_password
  }
}

