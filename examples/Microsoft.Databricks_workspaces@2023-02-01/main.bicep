param location string = 'eastus2'
param resource_name string = 'acctest0001'

resource workspace 'Microsoft.Databricks/workspaces@2023-02-01' = {
  location: location
  name: resource_name
  properties: {
    managedResourceGroupId: data.azapi_resource_id.workspace_resource_group.id
    parameters: {
      prepareEncryption: {
        value: true
      }
      requireInfrastructureEncryption: {
        value: true
      }
    }
    publicNetworkAccess: 'Enabled'
  }
  sku: {
    name: 'premium'
  }
}

