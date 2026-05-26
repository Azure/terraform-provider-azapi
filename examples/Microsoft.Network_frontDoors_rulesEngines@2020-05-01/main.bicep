param location string = 'westeurope'
param resource_name string = 'acctest0001'

resource rulesEngine 'Microsoft.Network/frontDoors/rulesEngines@2020-05-01' = {
  name: resource_name
  properties: {
    rules: [
      {
        action: {
          routeConfigurationOverride: {
            '@odata.type': '#Microsoft.Azure.FrontDoor.Models.FrontdoorRedirectConfiguration'
            customHost: 'customhost.org'
            redirectProtocol: 'HttpsOnly'
            redirectType: 'Found'
          }
        }
        matchProcessingBehavior: 'Continue'
        name: resource_name
        priority: 0
      }
    ]
  }
}

