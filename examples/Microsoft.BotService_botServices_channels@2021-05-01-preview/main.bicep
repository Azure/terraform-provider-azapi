param location string = 'westus'
param resource_name string = 'acctest0001'

resource botService 'Microsoft.BotService/botServices@2021-05-01-preview' = {
  location: location
  name: resource_name
  tags: {
    environment: 'production'
  }
  kind: 'bot'
  properties: {
    cmekKeyVaultUrl: ''
    description: ''
    developerAppInsightKey: ''
    developerAppInsightsApiKey: ''
    developerAppInsightsApplicationId: ''
    displayName: resource_name
    endpoint: ''
    iconUrl: 'https://docs.botframework.com/static/devportal/client/images/bot-framework-default.png'
    isCmekEnabled: false
    isStreamingSupported: false
    msaAppId: '12345678-1234-1234-1234-123456789012'
  }
  sku: {
    name: 'F0'
  }
}

resource channel 'Microsoft.BotService/botServices/channels@2021-05-01-preview' = {
  parent: botService
  location: location
  name: 'AlexaChannel'
  kind: 'bot'
  properties: {
    channelName: 'AlexaChannel'
    properties: {
      alexaSkillId: 'amzn1.ask.skill.19126e57-867f-4553-b953-ad0a720dddec'
      isEnabled: true
    }
  }
}

