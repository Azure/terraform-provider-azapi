param location string = 'eastus'
param resource_name string = 'acctest0001'

resource accessPolicyAssignment 'Microsoft.Cache/redis/accessPolicyAssignments@2024-03-01' = {
  parent: redis
  name: resource_name
  properties: {
    accessPolicyName: 'Data Contributor'
    objectId: data.azurerm_client_config.test.object_id
    objectIdAlias: 'ServicePrincipal'
  }
}

resource redis 'Microsoft.Cache/redis@2023-04-01' = {
  location: location
  name: resource_name
  properties: {
    enableNonSslPort: true
    minimumTlsVersion: '1.2'
    sku: {
      capacity: 2
      family: 'C'
      name: 'Standard'
    }
  }
}

