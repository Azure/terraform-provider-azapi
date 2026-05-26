param location string = 'westeurope'
param resource_name string = 'acctest0001'

resource healthBot 'Microsoft.HealthBot/healthBots@2022-08-08' = {
  location: location
  name: resource_name
  sku: {
    name: 'F0'
  }
}

