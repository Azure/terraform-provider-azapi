param location string = 'westeurope'
param resource_name string = 'acctest0001'
param subscription_id string = '00000000-0000-0000-0000-000000000000'
param vm_admin_password string = 'P@$$w0rd1234!'
param vm_admin_username string = 'adminuser'

resource storageAccount 'Microsoft.Storage/storageAccounts@2023-05-01' = {
  location: azurerm_resource_group.example.location
  name: '${resource_name}sa'
  kind: 'StorageV2'
  properties: {
    accessTier: 'Hot'
    allowBlobPublicAccess: false
    allowCrossTenantReplication: false
    allowSharedKeyAccess: true
    defaultToOAuthAuthentication: false
    dnsEndpointType: 'Standard'
    encryption: {
      keySource: 'Microsoft.Storage'
      services: {
        blob: {
          enabled: true
          keyType: 'Account'
        }
        file: {
          enabled: true
          keyType: 'Account'
        }
      }
    }
    isHnsEnabled: false
    isLocalUserEnabled: true
    isNfsV3Enabled: false
    isSftpEnabled: false
    minimumTlsVersion: 'TLS1_2'
    networkAcls: {
      bypass: 'AzureServices'
      defaultAction: 'Allow'
      ipRules: []
      resourceAccessRules: []
      virtualNetworkRules: []
    }
    publicNetworkAccess: 'Disabled'
    supportsHttpsTrafficOnly: true
  }
  sku: {
    name: 'Standard_LRS'
  }
}

