param location string = 'westeurope'
param resource_name string = 'acctest0001'

resource diskEncryptionSet 'Microsoft.Compute/diskEncryptionSets@2022-03-02' = {
  identity: [
    {
      identity_ids: []
      type: 'SystemAssigned'
    }
  ]
  location: location
  name: resource_name
  properties: {
    activeKey: {
      keyUrl: azapi_resource_action.key.output.properties.keyUriWithVersion
      sourceVault: {
        id: vault.id
      }
    }
    encryptionType: 'EncryptionAtRestWithCustomerKey'
    rotationToLatestKeyVersionEnabled: false
  }
}

resource vault 'Microsoft.KeyVault/vaults@2023-02-01' = {
  location: location
  name: resource_name
  properties: {
    accessPolicies: []
    enableSoftDelete: true
    sku: {
      family: 'A'
      name: 'standard'
    }
    tenantId: data.azurerm_client_config.current.tenant_id
  }
}

