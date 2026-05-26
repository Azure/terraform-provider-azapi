param location string = 'westeurope'
param resource_name string = 'acctest0001'

resource serverfarm 'Microsoft.Web/serverfarms@2022-09-01' = {
  location: location
  name: resource_name
  properties: {
    hyperV: false
    perSiteScaling: false
    reserved: false
    zoneRedundant: false
  }
  sku: {
    name: 'S1'
  }
}

resource site 'Microsoft.Web/sites@2022-09-01' = {
  location: location
  name: resource_name
  properties: {
    clientAffinityEnabled: false
    clientCertEnabled: false
    clientCertMode: 'Required'
    enabled: true
    httpsOnly: false
    publicNetworkAccess: 'Enabled'
    serverFarmId: serverfarm.id
    siteConfig: {
      acrUseManagedIdentityCreds: false
      alwaysOn: true
      autoHealEnabled: false
      ftpsState: 'Disabled'
      http20Enabled: false
      loadBalancing: 'LeastRequests'
      localMySqlEnabled: false
      managedPipelineMode: 'Integrated'
      minTlsVersion: '1.2'
      publicNetworkAccess: 'Enabled'
      remoteDebuggingEnabled: false
      scmIpSecurityRestrictionsUseMain: false
      scmMinTlsVersion: '1.2'
      use32BitWorkerProcess: true
      vnetRouteAllEnabled: false
      webSocketsEnabled: false
      windowsFxVersion: ''
    }
    vnetRouteAllEnabled: false
  }
}

