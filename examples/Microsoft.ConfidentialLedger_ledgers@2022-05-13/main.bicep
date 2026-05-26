param ledger_certificate string = null
param location string = 'westeurope'
param resource_name string = 'acctest0001'

resource ledger 'Microsoft.ConfidentialLedger/ledgers@2022-05-13' = {
  location: location
  name: resource_name
  properties: {
    aadBasedSecurityPrincipals: [
      {
        ledgerRoleName: 'Administrator'
        principalId: data.azurerm_client_config.current.object_id
        tenantId: data.azurerm_client_config.current.tenant_id
      }
    ]
    certBasedSecurityPrincipals: [
      {
        cert: ledger_certificate
        ledgerRoleName: 'Administrator'
      }
    ]
    ledgerType: 'Private'
  }
}

