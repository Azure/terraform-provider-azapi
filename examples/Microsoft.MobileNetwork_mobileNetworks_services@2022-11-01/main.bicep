param location string = 'eastus'
param resource_name string = 'acctest0001'

resource mobileNetwork 'Microsoft.MobileNetwork/mobileNetworks@2022-11-01' = {
  location: location
  name: resource_name
  properties: {
    publicLandMobileNetworkIdentifier: {
      mcc: '001'
      mnc: '01'
    }
  }
}

resource service 'Microsoft.MobileNetwork/mobileNetworks/services@2022-11-01' = {
  parent: mobileNetwork
  location: location
  name: resource_name
  properties: {
    pccRules: [
      {
        ruleName: 'default-rule'
        rulePrecedence: 1
        serviceDataFlowTemplates: [
          {
            direction: 'Uplink'
            ports: []
            protocol: [
              'ip'
            ]
            remoteIpList: [
              '10.3.4.0/24'
            ]
            templateName: 'IP-to-server'
          }
        ]
        trafficControl: 'Enabled'
      }
    ]
    servicePrecedence: 0
  }
}

