param location string = 'eastus'
param resource_name string = 'acctest0001'

resource dataNetwork 'Microsoft.MobileNetwork/mobileNetworks/dataNetworks@2022-11-01' = {
  parent: mobileNetwork
  location: location
  name: resource_name
  properties: {}
}

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

resource simPolicy 'Microsoft.MobileNetwork/mobileNetworks/simPolicies@2022-11-01' = {
  parent: mobileNetwork
  location: location
  name: resource_name
  properties: {
    defaultSlice: {
      id: slice.id
    }
    registrationTimer: 3240
    sliceConfigurations: [
      {
        dataNetworkConfigurations: [
          {
            '5qi': 9
            additionalAllowedSessionTypes: null
            allocationAndRetentionPriorityLevel: 9
            allowedServices: [
              {
                id: service.id
              }
            ]
            dataNetwork: {
              id: dataNetwork.id
            }
            defaultSessionType: 'IPv4'
            maximumNumberOfBufferedPackets: 10
            preemptionCapability: 'NotPreempt'
            preemptionVulnerability: 'Preemptable'
            sessionAmbr: {
              downlink: '1 Gbps'
              uplink: '500 Mbps'
            }
          }
        ]
        defaultDataNetwork: {
          id: dataNetwork.id
        }
        slice: {
          id: slice.id
        }
      }
    ]
    ueAmbr: {
      downlink: '1 Gbps'
      uplink: '500 Mbps'
    }
  }
  tags: {
    key: 'value'
  }
}

resource slice 'Microsoft.MobileNetwork/mobileNetworks/slices@2022-11-01' = {
  parent: mobileNetwork
  location: location
  name: resource_name
  properties: {
    snssai: {
      sst: 1
    }
  }
}

