param location string = 'westeurope'
param resource_name string = 'acctest0001'

resource managedHSM 'Microsoft.KeyVault/managedHSMs@2021-10-01' = {
  location: location
  name: 'kvHsm230630033342437496'
  properties: {
    createMode: 'default'
    enablePurgeProtection: false
    enableSoftDelete: true
    initialAdminObjectIds: [
      data.azurerm_client_config.current.object_id
    ]
    publicNetworkAccess: 'Enabled'
    softDeleteRetentionInDays: 90
    tenantId: data.azurerm_client_config.current.tenant_id
  }
  sku: {
    family: 'B'
    name: 'Standard_B1'
  }
}

