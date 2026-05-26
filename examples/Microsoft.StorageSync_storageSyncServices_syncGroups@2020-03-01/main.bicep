param location string = 'westeurope'
param resource_name string = 'acctest0001'

resource storageSyncService 'Microsoft.StorageSync/storageSyncServices@2020-03-01' = {
  location: location
  name: resource_name
  properties: {
    incomingTrafficPolicy: 'AllowAllTraffic'
  }
}

resource syncGroup 'Microsoft.StorageSync/storageSyncServices/syncGroups@2020-03-01' = {
  parent: storageSyncService
  name: resource_name
}

