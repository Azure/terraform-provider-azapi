param location string = 'westeurope'
param resource_name string = 'acctest0001'

resource component 'Microsoft.Insights/components@2020-02-02' = {
  location: location
  name: resource_name
  kind: 'web'
  properties: {
    Application_Type: 'web'
    DisableIpMasking: false
    DisableLocalAuth: false
    ForceCustomerStorageForProfiler: false
    RetentionInDays: 90
    SamplingPercentage: 100
    publicNetworkAccessForIngestion: 'Enabled'
    publicNetworkAccessForQuery: 'Enabled'
  }
}

resource compute 'Microsoft.MachineLearningServices/workspaces/computes@2022-05-01' = {
  parent: workspace
  location: location
  name: resource_name
  properties: {
    computeLocation: 'westeurope'
    computeType: 'ComputeInstance'
    description: ''
    disableLocalAuth: true
    properties: {
      vmSize: 'STANDARD_D2_V2'
    }
  }
}

resource storageAccount 'Microsoft.Storage/storageAccounts@2021-09-01' = {
  location: location
  name: resource_name
  kind: 'StorageV2'
  sku: {
    name: 'Standard_LRS'
  }
}

resource vault 'Microsoft.KeyVault/vaults@2021-10-01' = {
  location: location
  name: resource_name
  properties: {
    accessPolicies: []
    createMode: 'default'
    enablePurgeProtection: true
    enableRbacAuthorization: false
    enableSoftDelete: true
    enabledForDeployment: false
    enabledForDiskEncryption: false
    enabledForTemplateDeployment: false
    publicNetworkAccess: 'Enabled'
    sku: {
      family: 'A'
      name: 'standard'
    }
    tenantId: data.azurerm_client_config.current.tenant_id
  }
}

resource workspace 'Microsoft.MachineLearningServices/workspaces@2022-05-01' = {
  identity: [
    {
      identity_ids: []
      type: 'SystemAssigned'
    }
  ]
  location: location
  name: resource_name
  properties: {
    applicationInsights: component.id
    keyVault: vault.id
    publicNetworkAccess: 'Enabled'
    storageAccount: storageAccount.id
    v1LegacyMode: false
  }
  sku: {
    name: 'Basic'
    tier: 'Basic'
  }
}

