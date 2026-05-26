param administrator_login_password string = null
param location string = 'westeurope'
param resource_name string = 'acctest0001'

resource administrator 'Microsoft.Sql/servers/administrators@2020-11-01-preview' = {
  parent: server
  name: 'ActiveDirectory'
  properties: {
    administratorType: 'ActiveDirectory'
    login: 'sqladmin'
    sid: data.azurerm_client_config.current.client_id
    tenantId: data.azurerm_client_config.current.tenant_id
  }
}

resource server 'Microsoft.Sql/servers@2015-05-01-preview' = {
  location: location
  name: resource_name
  properties: {
    administratorLogin: 'mradministrator'
    administratorLoginPassword: administrator_login_password
    version: '12.0'
  }
}

