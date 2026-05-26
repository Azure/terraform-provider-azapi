param location string = 'westeurope'
param resource_name string = 'acctest0001'

resource authorizationRule 'Microsoft.EventHub/namespaces/authorizationRules@2021-11-01' = {
  parent: namespace
  name: 'example'
  properties: {
    rights: [
      'Listen'
      'Send'
      'Manage'
    ]
  }
}

resource diagnosticSetting 'Microsoft.Insights/diagnosticSettings@2021-05-01-preview' = {
  parent: vault
  name: resource_name
  properties: {
    eventHubAuthorizationRuleId: authorizationRule.id
    eventHubName: namespace.name
    logs: [
      {
        categoryGroup: 'Audit'
        enabled: true
        retentionPolicy: {
          days: 0
          enabled: false
        }
      }
    ]
    metrics: [
      {
        category: 'AllMetrics'
        enabled: true
        retentionPolicy: {
          days: 0
          enabled: false
        }
      }
    ]
  }
}

resource namespace 'Microsoft.EventHub/namespaces@2022-01-01-preview' = {
  location: location
  name: resource_name
  properties: {
    disableLocalAuth: false
    isAutoInflateEnabled: false
    publicNetworkAccess: 'Enabled'
    zoneRedundant: false
  }
  sku: {
    capacity: 1
    name: 'Basic'
    tier: 'Basic'
  }
}

resource vault 'Microsoft.KeyVault/vaults@2021-10-01' = {
  location: location
  name: resource_name
  properties: {
    accessPolicies: []
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
    tenantId: data.azurerm_client_config.current.tenant_id
  }
}

