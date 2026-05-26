param location string = 'westus'
param resource_name string = 'acctest0001'

resource authorizationRule 'Microsoft.EventHub/namespaces/authorizationRules@2024-01-01' = {
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

resource diagnosticSetting 'Microsoft.AADIAM/diagnosticSettings@2017-04-01' = {
  name: '${resource_name}-DS-unique'
  properties: {
    eventHubAuthorizationRuleId: authorizationRule.id
    eventHubName: eventhub.name
    logs: [
      {
        category: 'RiskyUsers'
        enabled: true
      }
      {
        category: 'ServicePrincipalSignInLogs'
        enabled: true
      }
      {
        category: 'SignInLogs'
        enabled: true
      }
      {
        category: 'B2CRequestLogs'
        enabled: true
      }
      {
        category: 'UserRiskEvents'
        enabled: true
      }
      {
        category: 'NonInteractiveUserSignInLogs'
        enabled: true
      }
      {
        category: 'AuditLogs'
        enabled: true
      }
    ]
  }
}

resource eventhub 'Microsoft.EventHub/namespaces/eventhubs@2024-01-01' = {
  parent: namespace
  name: '${resource_name}-EH-unique'
  properties: {
    messageRetentionInDays: 1
    partitionCount: 2
    status: 'Active'
  }
}

resource namespace 'Microsoft.EventHub/namespaces@2024-01-01' = {
  location: location
  name: '${resource_name}-EHN-unique'
  properties: {
    disableLocalAuth: false
    isAutoInflateEnabled: false
    minimumTlsVersion: '1.2'
    publicNetworkAccess: 'Enabled'
  }
  sku: {
    capacity: 1
    name: 'Basic'
    tier: 'Basic'
  }
}

