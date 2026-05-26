param client_secret string = null
param location string = 'westeurope'
param resource_name string = 'acctest0001'

resource botService 'Microsoft.BotService/botServices@2021-05-01-preview' = {
  location: 'global'
  name: resource_name
  kind: 'bot'
  properties: {
    displayName: resource_name
    isCmekEnabled: false
    isStreamingSupported: false
    msaAppId: data.azurerm_client_config.current.tenant_id
  }
  sku: {
    name: 'F0'
  }
}

resource connection 'Microsoft.BotService/botServices/connections@2021-05-01-preview' = {
  parent: botService
  location: 'global'
  name: resource_name
  kind: 'bot'
  properties: {
    clientId: botService.properties.msaAppId
    clientSecret: client_secret
    scopes: ''
    serviceProviderId: data.azapi_resource_action.listAuthServiceProviders.output.value[36].properties.id
  }
}

