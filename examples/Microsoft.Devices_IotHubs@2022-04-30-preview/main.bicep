param location string = 'westeurope'
param resource_name string = 'acctest0001'

resource IotHub 'Microsoft.Devices/IotHubs@2022-04-30-preview' = {
  location: location
  name: resource_name
  properties: {
    cloudToDevice: {}
    enableFileUploadNotifications: false
    messagingEndpoints: {}
    routing: {
      fallbackRoute: {
        condition: 'true'
        endpointNames: [
          'events'
        ]
        isEnabled: true
        source: 'DeviceMessages'
      }
    }
    storageEndpoints: {}
  }
  sku: {
    capacity: 1
    name: 'S1'
  }
}

