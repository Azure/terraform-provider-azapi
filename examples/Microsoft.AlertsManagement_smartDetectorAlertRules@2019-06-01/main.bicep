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
    groupShortName: 'acctestag'
    itsmReceivers: []
    logicAppReceivers: []
    smsReceivers: []
    voiceReceivers: []
    webhookReceivers: []
  }
}

resource component 'Microsoft.Insights/components@2020-02-02' = {
  location: location
  name: resource_name
  kind: 'web'
  properties: {
    Application_Type: 'web'
    DisableIpMasking: false
    DisableLocalAuth: false
    ForceCustomerStorageForProfiler: false
    RetentionInDays: 90
    SamplingPercentage: 100
    publicNetworkAccessForIngestion: 'Enabled'
    publicNetworkAccessForQuery: 'Enabled'
  }
}

resource smartDetectorAlertRule 'microsoft.alertsManagement/smartDetectorAlertRules@2019-06-01' = {
  location: 'global'
  name: resource_name
  properties: {
    actionGroups: {
      customEmailSubject: ''
      customWebhookPayload: ''
      groupIds: [
        actionGroup.id
      ]
    }
    description: ''
    detector: {
      id: 'FailureAnomaliesDetector'
    }
    frequency: 'PT1M'
    scope: [
      component.id
    ]
    severity: 'Sev0'
    state: 'Enabled'
  }
}

