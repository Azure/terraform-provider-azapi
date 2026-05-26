param location string = 'eastus'
param resource_name string = 'acctest0001'

resource connection 'Microsoft.Web/connections@2016-06-01' = {
  location: location
  name: resource_name
  properties: {
    api: {
      id: data.azurerm_managed_api.test.id
    }
    displayName: 'Service Bus'
  }
}

resource namespaces 'Microsoft.ServiceBus/namespaces@2022-10-01-preview' = {
  location: location
  name: resource_name
  identity: {
    type: 'None'
    userAssignedIdentities: null
  }
  properties: {
    disableLocalAuth: false
    minimumTlsVersion: '1.2'
    premiumMessagingPartitions: 0
    publicNetworkAccess: 'Enabled'
  }
  sku: {
    capacity: 0
    name: 'Basic'
    tier: 'Basic'
  }
}

resource workflows 'Microsoft.Logic/workflows@2019-05-01' = {
  location: location
  name: resource_name
  identity: {
    type: 'None'
    userAssignedIdentities: null
  }
  properties: {
    definition: {
      '$schema': 'https://schema.management.azure.com/providers/Microsoft.Logic/schemas/2016-06-01/workflowdefinition.json#'
      contentVersion: '1.0.0.0'
    }
    state: 'Enabled'
  }
}

