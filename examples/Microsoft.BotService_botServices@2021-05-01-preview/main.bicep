param location string = 'westeurope'
param resource_name string = 'acctest0001'

resource botService 'Microsoft.BotService/botServices@2021-05-01-preview' = {
  location: 'global'
  name: resource_name
  kind: 'sdk'
  properties: {
    developerAppInsightKey: ''
    developerAppInsightsApiKey: ''
    developerAppInsightsApplicationId: ''
    displayName: resource_name
    endpoint: ''
    luisAppIds: []
    luisKey: ''
    msaAppId: data.azurerm_client_config.current.client_id
  }
  sku: {
    name: 'F0'
  }
  tags: {
    environment: 'production'
  }
}

