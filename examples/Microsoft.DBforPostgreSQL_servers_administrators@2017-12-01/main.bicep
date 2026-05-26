param admin_login string = null
param administrator_login string = null
param administrator_login_password string = null
param location string = 'westeurope'
param resource_name string = 'acctest0001'

resource administrator 'Microsoft.DBforPostgreSQL/servers/administrators@2017-12-01' = {
  parent: server
  name: 'activeDirectory'
  properties: {
    administratorType: 'ActiveDirectory'
    login: admin_login
    sid: data.azurerm_client_config.current.client_id
    tenantId: data.azurerm_client_config.current.tenant_id
  }
}

resource server 'Microsoft.DBforPostgreSQL/servers@2017-12-01' = {
  location: location
  name: resource_name
  properties: {
    administratorLogin: administrator_login
    administratorLoginPassword: administrator_login_password
    createMode: 'Default'
    infrastructureEncryption: 'Disabled'
    minimalTlsVersion: 'TLS1_2'
    publicNetworkAccess: 'Enabled'
    sslEnforcement: 'Enabled'
    storageProfile: {
      backupRetentionDays: 7
      storageAutogrow: 'Enabled'
      storageMB: 51200
    }
    version: '9.6'
  }
  sku: {
    capacity: 2
    family: 'Gen5'
    name: 'GP_Gen5_2'
    tier: 'GeneralPurpose'
  }
}

