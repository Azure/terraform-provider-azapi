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

resource account 'Microsoft.DeviceUpdate/accounts@2022-10-01' = {
  location: location
  name: resource_name
  properties: {
    publicNetworkAccess: 'Enabled'
    sku: 'Standard'
  }
}

resource instance 'Microsoft.DeviceUpdate/accounts/instances@2022-10-01' = {
  parent: account
  location: location
  name: resource_name
  properties: {
    accountName: account.name
    enableDiagnostics: false
    iotHubs: [
      {
        resourceId: IotHub.id
      }
    ]
  }
}

