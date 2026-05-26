param location string = 'westeurope'
param resource_name string = 'acctest0001'

resource provisioningService 'Microsoft.Devices/provisioningServices@2022-02-05' = {
  location: location
  name: resource_name
  properties: {
    allocationPolicy: 'Hashed'
    enableDataResidency: false
    iotHubs: []
    publicNetworkAccess: 'Enabled'
  }
  sku: {
    capacity: 1
    name: 'S1'
  }
}

