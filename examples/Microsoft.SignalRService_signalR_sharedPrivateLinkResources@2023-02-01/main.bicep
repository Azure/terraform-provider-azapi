param location string = 'westeurope'
param resource_name string = 'acctest0001'

resource sharedPrivateLinkResource 'Microsoft.SignalRService/signalR/sharedPrivateLinkResources@2023-02-01' = {
  parent: signalR
  name: resource_name
  properties: {
    groupId: 'vault'
    privateLinkResourceId: vault.id
    requestMessage: 'please approve'
  }
}

resource signalR 'Microsoft.SignalRService/signalR@2023-02-01' = {
  location: location
  name: resource_name
  properties: {
    cors: {}
    disableAadAuth: false
    disableLocalAuth: false
    features: [
      {
        flag: 'ServiceMode'
        value: 'Default'
      }
      {
        flag: 'EnableConnectivityLogs'
        value: 'False'
      }
      {
        flag: 'EnableMessagingLogs'
        value: 'False'
      }
      {
        flag: 'EnableLiveTrace'
        value: 'False'
      }
    ]
    publicNetworkAccess: 'Enabled'
    resourceLogConfiguration: {
      categories: [
        {
          enabled: 'false'
          name: 'MessagingLogs'
        }
        {
          enabled: 'false'
          name: 'ConnectivityLogs'
        }
        {
          enabled: 'false'
          name: 'HttpRequestLogs'
        }
      ]
    }
    serverless: {
      connectionTimeoutInSeconds: 30
    }
    tls: {
      clientCertEnabled: false
    }
    upstream: {
      templates: []
    }
  }
  sku: {
    capacity: 1
    name: 'Standard_S1'
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

