param location string = 'westeurope'
param resource_name string = 'acctest0001'

resource vault 'Microsoft.KeyVault/vaults@2023-02-01' = {
  location: location
  name: resource_name
  properties: {
    accessPolicies: []
    enableSoftDelete: true
    sku: {
      family: 'A'
      name: 'standard'
    }
    tenantId: data.azurerm_client_config.current.tenant_id
  }
}

