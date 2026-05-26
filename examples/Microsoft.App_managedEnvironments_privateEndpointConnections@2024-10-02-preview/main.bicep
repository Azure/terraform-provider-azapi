param location string = 'westeurope'
param resource_name string = 'acctest5925'

resource managedEnvironment 'Microsoft.App/managedEnvironments@2024-10-02-preview' = {
  location: location
  name: resource_name
  properties: {
    appLogsConfiguration: {
      destination: 'log-analytics'
      logAnalyticsConfiguration: {
        customerId: workspace.properties.customerId
        sharedKey: data.azapi_resource_action.sharedKeys.output.primarySharedKey
      }
    }
    publicNetworkAccess: 'Disabled'
    vnetConfiguration: {}
    workloadProfiles: [
      {
        name: 'Consumption'
        workloadProfileType: 'Consumption'
      }
    ]
  }
}

resource privateDnsZone 'Microsoft.Network/privateDnsZones@2020-06-01' = {
  location: 'global'
  name: 'acctestzone.azurecontainerapps.dev'
  properties: {}
}

resource privateDnsZoneGroup 'Microsoft.Network/privateEndpoints/privateDnsZoneGroups@2023-05-01' = {
  parent: privateEndpoint
  name: 'default'
  properties: {
    privateDnsZoneConfigs: [
      {
        name: 'config'
        properties: {
          privateDnsZoneId: privateDnsZone.id
        }
      }
    ]
  }
}

resource privateEndpoint 'Microsoft.Network/privateEndpoints@2023-05-01' = {
  location: location
  name: '${resource_name}-pe'
  properties: {
    privateLinkServiceConnections: [
      {
        name: '${resource_name}-connection'
        properties: {
          groupIds: [
            'managedEnvironments'
          ]
          privateLinkServiceId: managedEnvironment.id
        }
      }
    ]
    subnet: {
      id: subnet.id
    }
  }
}

resource subnet 'Microsoft.Network/virtualNetworks/subnets@2023-05-01' = {
  parent: virtualNetwork
  name: '${resource_name}-subnet'
  properties: {
    addressPrefix: '10.0.0.0/21'
    privateEndpointNetworkPolicies: 'Disabled'
  }
}

resource virtualNetwork 'Microsoft.Network/virtualNetworks@2023-05-01' = {
  location: location
  name: '${resource_name}-vnet'
  properties: {
    addressSpace: {
      addressPrefixes: [
        '10.0.0.0/16'
      ]
    }
  }
}

resource vnetLink 'Microsoft.Network/privateDnsZones/virtualNetworkLinks@2020-06-01' = {
  parent: privateDnsZone
  location: 'global'
  name: '${resource_name}-vnet-link'
  properties: {
    registrationEnabled: false
    virtualNetwork: {
      id: virtualNetwork.id
    }
  }
}

resource workspace 'Microsoft.OperationalInsights/workspaces@2022-10-01' = {
  location: location
  name: resource_name
  properties: {
    features: {
      disableLocalAuth: false
      enableLogAccessUsingOnlyResourcePermissions: true
    }
    publicNetworkAccessForIngestion: 'Enabled'
    publicNetworkAccessForQuery: 'Enabled'
    retentionInDays: 30
    sku: {
      name: 'PerGB2018'
    }
    workspaceCapping: {
      dailyQuotaGb: -1
    }
  }
}

