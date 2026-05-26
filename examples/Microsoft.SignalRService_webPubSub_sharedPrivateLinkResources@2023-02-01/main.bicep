param location string = 'westeurope'
param resource_name string = 'acctest0001'

resource sharedPrivateLinkResource 'Microsoft.SignalRService/webPubSub/sharedPrivateLinkResources@2023-02-01' = {
  parent: webPubSub
  name: resource_name
  properties: {
    groupId: 'vault'
    privateLinkResourceId: vault.id
  }
}

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

resource webPubSub 'Microsoft.SignalRService/webPubSub@2023-02-01' = {
  location: location
  name: resource_name
  properties: {
    disableAadAuth: false
    disableLocalAuth: false
    publicNetworkAccess: 'Enabled'
    tls: {
      clientCertEnabled: false
    }
  }
  sku: {
    capacity: 1
    name: 'Standard_S1'
  }
}

