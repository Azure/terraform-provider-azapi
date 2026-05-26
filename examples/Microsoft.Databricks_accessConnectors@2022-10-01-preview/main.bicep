param location string = 'westeurope'
param resource_name string = 'acctest0001'

resource accessConnector 'Microsoft.Databricks/accessConnectors@2022-10-01-preview' = {
  location: location
  name: resource_name
}

