param location string = 'westeurope'
param resource_name string = 'acctest0001'

resource configurationStore 'Microsoft.AppConfiguration/configurationStores@2023-03-01' = {
  location: location
  name: resource_name
  properties: {
    disableLocalAuth: false
    enablePurgeProtection: false
  }
  sku: {
    name: 'standard'
  }
}

