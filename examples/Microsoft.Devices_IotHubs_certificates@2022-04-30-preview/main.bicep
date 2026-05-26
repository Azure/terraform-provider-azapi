param certificate_content string = null
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
    name: 'B1'
  }
}

resource certificate 'Microsoft.Devices/IotHubs/certificates@2022-04-30-preview' = {
  parent: IotHub
  name: resource_name
  properties: {
    certificate: certificate_content
    isVerified: false
  }
}

