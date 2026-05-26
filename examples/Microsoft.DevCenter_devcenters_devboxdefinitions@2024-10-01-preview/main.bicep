param location string = 'westeurope'
param resource_name string = 'acctest0001'

resource devBoxDefinition 'Microsoft.DevCenter/devcenters/devboxdefinitions@2024-10-01-preview' = {
  parent: devCenter
  location: location
  name: resource_name
  properties: {
    hibernateSupport: 'Enabled'
    imageReference: {
      id: '${devCenter.id}/galleries/default/images/microsoftvisualstudio_visualstudioplustools_vs-2022-ent-general-win10-m365-gen2'
    }
    sku: {
      name: 'general_i_8c32gb256ssd_v2'
    }
  }
}

resource devCenter 'Microsoft.DevCenter/devcenters@2023-04-01' = {
  location: location
  name: resource_name
  identity: {
    type: 'SystemAssigned'
    userAssignedIdentities: null
  }
}

