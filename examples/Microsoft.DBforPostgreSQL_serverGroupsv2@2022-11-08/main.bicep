param administrator_login_password string = null
param location string = 'westeurope'
param resource_name string = 'acctest0001'

resource serverGroupsv2 'Microsoft.DBforPostgreSQL/serverGroupsv2@2022-11-08' = {
  location: location
  name: resource_name
  properties: {
    administratorLoginPassword: administrator_login_password
    coordinatorEnablePublicIpAccess: true
    coordinatorServerEdition: 'GeneralPurpose'
    coordinatorStorageQuotaInMb: 131072
    coordinatorVCores: 2
    enableHa: false
    nodeCount: 0
    nodeEnablePublicIpAccess: false
    nodeServerEdition: 'MemoryOptimized'
  }
}

