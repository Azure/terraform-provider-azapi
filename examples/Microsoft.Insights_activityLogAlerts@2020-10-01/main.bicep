param location string = 'westeurope'
param resource_name string = 'acctest0001'

resource actionGroup 'Microsoft.Insights/actionGroups@2023-01-01' = {
  location: 'global'
  name: resource_name
  properties: {
    armRoleReceivers: []
    automationRunbookReceivers: []
    azureAppPushReceivers: []
    azureFunctionReceivers: []
    emailReceivers: []
    enabled: true
    eventHubReceivers: []
    groupShortName: 'acctestag1'
    itsmReceivers: []
    logicAppReceivers: []
    smsReceivers: []
    voiceReceivers: []
    webhookReceivers: []
  }
}

resource actionGroup2 'Microsoft.Insights/actionGroups@2023-01-01' = {
  location: 'global'
  name: resource_name
  properties: {
    armRoleReceivers: []
    automationRunbookReceivers: []
    azureAppPushReceivers: []
    azureFunctionReceivers: []
    emailReceivers: []
    enabled: true
    eventHubReceivers: []
    groupShortName: 'acctestag2'
    itsmReceivers: []
    logicAppReceivers: []
    smsReceivers: []
    voiceReceivers: []
    webhookReceivers: []
  }
}

resource activityLogAlert 'Microsoft.Insights/activityLogAlerts@2020-10-01' = {
  location: 'global'
  name: resource_name
  properties: {
    actions: {
      actionGroups: [
        {
          actionGroupId: actionGroup.id
          webhookProperties: {}
        }
        {
          actionGroupId: actionGroup2.id
          webhookProperties: {
            from: 'terraform test'
            to: 'microsoft azure'
          }
        }
      ]
    }
    condition: {
      allOf: [
        {
          equals: 'ResourceHealth'
          field: 'category'
        }
        {
          anyOf: [
            {
              equals: 'Unavailable'
              field: 'properties.currentHealthStatus'
            }
            {
              equals: 'Degraded'
              field: 'properties.currentHealthStatus'
            }
          ]
        }
        {
          anyOf: [
            {
              equals: 'Unknown'
              field: 'properties.previousHealthStatus'
            }
            {
              equals: 'Available'
              field: 'properties.previousHealthStatus'
            }
          ]
        }
        {
          anyOf: [
            {
              equals: 'PlatformInitiated'
              field: 'properties.cause'
            }
            {
              equals: 'UserInitiated'
              field: 'properties.cause'
            }
          ]
        }
      ]
    }
    description: 'This is just a test acceptance.'
    enabled: true
    scopes: [
      resourceGroup.id
      storageAccount.id
    ]
  }
}

resource storageAccount 'Microsoft.Storage/storageAccounts@2021-09-01' = {
  location: location
  name: resource_name
  kind: 'StorageV2'
  properties: {
    accessTier: 'Hot'
    allowBlobPublicAccess: true
    allowCrossTenantReplication: true
    allowSharedKeyAccess: true
    defaultToOAuthAuthentication: false
    encryption: {
      keySource: 'Microsoft.Storage'
      services: {
        queue: {
          keyType: 'Service'
        }
        table: {
          keyType: 'Service'
        }
      }
    }
    isHnsEnabled: false
    isNfsV3Enabled: false
    isSftpEnabled: false
    minimumTlsVersion: 'TLS1_2'
    networkAcls: {
      defaultAction: 'Allow'
    }
    publicNetworkAccess: 'Enabled'
    supportsHttpsTrafficOnly: true
  }
  sku: {
    name: 'Standard_LRS'
  }
}

