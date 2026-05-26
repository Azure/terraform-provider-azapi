param location string = 'eastus'
param resource_name string = 'acctest0001'

resource credentialSet 'Microsoft.ContainerRegistry/registries/credentialSets@2023-07-01' = {
  parent: registry
  identity: [
    {
      identity_ids: []
      type: 'SystemAssigned'
    }
  ]
  name: '${resource_name}-acr-credential-set'
  properties: {
    authCredentials: [
      {
        name: 'Credential1'
        passwordSecretIdentifier: 'https://${resource_name}vault.vault.azure.net/secrets/password'
        usernameSecretIdentifier: 'https://${resource_name}vault.vault.azure.net/secrets/username'
      }
    ]
    loginServer: 'docker.io'
  }
}

resource passwordSecret 'Microsoft.KeyVault/vaults/secrets@2023-02-01' = {
  parent: vault
  name: 'password'
  properties: {
    value: 'testpassword'
  }
}

resource registry 'Microsoft.ContainerRegistry/registries@2023-11-01-preview' = {
  location: location
  name: resource_name
  properties: {
    adminUserEnabled: false
    anonymousPullEnabled: false
    dataEndpointEnabled: false
    networkRuleBypassOptions: 'AzureServices'
    policies: {
      exportPolicy: {
        status: 'enabled'
      }
      quarantinePolicy: {
        status: 'disabled'
      }
      retentionPolicy: {}
      trustPolicy: {}
    }
    publicNetworkAccess: 'Enabled'
    zoneRedundancy: 'Disabled'
  }
  sku: {
    name: 'Basic'
  }
}

resource usernameSecret 'Microsoft.KeyVault/vaults/secrets@2023-02-01' = {
  parent: vault
  name: 'username'
  properties: {
    value: 'testuser'
  }
}

resource vault 'Microsoft.KeyVault/vaults@2023-02-01' = {
  location: location
  name: '${resource_name}vault'
  properties: {
    accessPolicies: [
      {
        objectId: data.azapi_client_config.current.object_id
        permissions: {
          certificates: []
          keys: []
          secrets: [
            'Get'
            'Set'
            'Delete'
            'Purge'
          ]
          storage: []
        }
        tenantId: data.azapi_client_config.current.tenant_id
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
    tenantId: data.azapi_client_config.current.tenant_id
  }
}

