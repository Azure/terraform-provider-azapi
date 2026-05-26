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

