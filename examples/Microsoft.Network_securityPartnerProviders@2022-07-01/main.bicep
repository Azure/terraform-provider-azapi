param location string = 'westeurope'
param resource_name string = 'acctest0001'

resource securityPartnerProvider 'Microsoft.Network/securityPartnerProviders@2022-07-01' = {
  location: location
  name: resource_name
  properties: {
    securityProviderName: 'ZScaler'
  }
}

