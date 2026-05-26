param location string = 'westeurope'
param resource_name string = 'acctest0001'

resource vault 'Microsoft.KeyVault/vaults@2021-10-01' = {
  location: location
  name: resource_name
  properties: {
    accessPolicies: [
      {
        objectId: data.azurerm_client_config.current.object_id
        permissions: {
          certificates: [
            'ManageContacts'
          ]
          keys: [
            'Create'
          ]
          secrets: [
            'Set'
          ]
          storage: []
        }
        tenantId: data.azurerm_client_config.current.tenant_id
      }
    ]
    createMode: 'default'
    enableRbacAuthorization: false
    enableSoftDelete: true
    enabledForDeployment: false
    enabledForDiskEncryption: false
    enabledForTemplateDeployment: false
    publicNetworkAccess: 'Enabled'
    sku: {
      family: 'A'
      name: 'standard'
    }
    softDeleteRetentionInDays: 7
    tenantId: data.azurerm_client_config.current.tenant_id
  }
}

