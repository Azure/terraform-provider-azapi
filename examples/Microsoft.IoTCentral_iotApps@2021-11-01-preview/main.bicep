param location string = 'westeurope'
param resource_name string = 'acctest0001'

resource iotApp 'Microsoft.IoTCentral/iotApps@2021-11-01-preview' = {
  location: location
  name: resource_name
  properties: {
    displayName: resource_name
    publicNetworkAccess: 'Enabled'
    subdomain: 'subdomain-2306300333537'
    template: 'iotc-pnp-preview@1.0.0'
  }
  sku: {
    name: 'ST1'
  }
}

