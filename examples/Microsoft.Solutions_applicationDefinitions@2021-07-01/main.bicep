param location string = 'westeurope'
param resource_name string = 'acctest0001'

resource applicationDefinition 'Microsoft.Solutions/applicationDefinitions@2021-07-01' = {
  location: location
  name: resource_name
  properties: {
    authorizations: [
      {
        principalId: data.azurerm_client_config.current.object_id
        roleDefinitionId: data.azapi_resource_action.roleDefinitions.output.value[0].name
      }
    ]
    description: 'Test Managed App Definition'
    displayName: 'TestManagedAppDefinition'
    isEnabled: true
    lockLevel: 'ReadOnly'
    packageFileUri: 'https://github.com/Azure/azure-managedapp-samples/raw/master/Managed Application Sample Packages/201-managed-storage-account/managedstorage.zip'
  }
}

