param location string = 'westeurope'
param resource_name string = 'acctest0001'

resource origin 'Microsoft.Cdn/profiles/originGroups/origins@2021-06-01' = {
  parent: originGroup
  name: resource_name
  properties: {
    enabledState: 'Enabled'
    enforceCertificateNameCheck: false
    hostName: 'contoso.com'
    httpPort: 80
    httpsPort: 443
    originHostHeader: 'www.contoso.com'
    priority: 1
    weight: 1
  }
}

resource originGroup 'Microsoft.Cdn/profiles/originGroups@2021-06-01' = {
  parent: profile
  name: resource_name
  properties: {
    loadBalancingSettings: {
      additionalLatencyInMilliseconds: 0
      sampleSize: 16
      successfulSamplesRequired: 3
    }
    sessionAffinityState: 'Enabled'
    trafficRestorationTimeToHealedOrNewEndpointsInMinutes: 10
  }
}

resource profile 'Microsoft.Cdn/profiles@2021-06-01' = {
  location: 'global'
  name: resource_name
  properties: {
    originResponseTimeoutSeconds: 120
  }
  sku: {
    name: 'Standard_AzureFrontDoor'
  }
}

