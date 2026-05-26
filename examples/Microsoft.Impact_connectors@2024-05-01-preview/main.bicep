param location string = 'westeurope'
param resource_name string = 'acctest0001'

resource connector 'Microsoft.Impact/connectors@2024-05-01-preview' = {
  name: resource_name
  properties: {
    connectorType: 'AzureMonitor'
  }
}

