param location string = 'westeurope'
param resource_name string = 'acctest0001'

resource capacity 'Microsoft.PowerBIDedicated/capacities@2021-01-01' = {
  location: location
  name: resource_name
  properties: {
    administration: {
      members: [
        data.azurerm_client_config.current.object_id
      ]
    }
    mode: 'Gen2'
  }
  sku: {
    name: 'A1'
  }
}

