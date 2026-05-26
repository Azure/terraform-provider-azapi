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

