param location string = 'westeurope'
param resource_name string = 'acctest0001'

resource component 'Microsoft.Insights/components@2020-02-02' = {
  location: location
  name: resource_name
  kind: 'web'
  properties: {
    Application_Type: 'web'
    DisableIpMasking: false
    DisableLocalAuth: false
    ForceCustomerStorageForProfiler: false
    RetentionInDays: 90
    SamplingPercentage: 100
    publicNetworkAccessForIngestion: 'Enabled'
    publicNetworkAccessForQuery: 'Enabled'
  }
}

resource privateLinkScope 'Microsoft.Insights/privateLinkScopes@2019-10-17-preview' = {
  location: 'Global'
  name: resource_name
  properties: {}
}

resource scopedResource 'Microsoft.Insights/privateLinkScopes/scopedResources@2019-10-17-preview' = {
  parent: privateLinkScope
  name: resource_name
  properties: {
    linkedResourceId: component.id
  }
}

