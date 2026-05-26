param location string = 'westeurope'
param resource_name string = 'acctest0001'

resource maintenanceConfiguration 'Microsoft.Maintenance/maintenanceConfigurations@2022-07-01-preview' = {
  location: location
  name: resource_name
  properties: {
    extensionProperties: {}
    maintenanceScope: 'SQLDB'
    namespace: 'Microsoft.Maintenance'
    visibility: 'Custom'
  }
}

