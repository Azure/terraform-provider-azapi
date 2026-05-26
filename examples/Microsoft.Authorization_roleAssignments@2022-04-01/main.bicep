param location string = 'eastus'
param resource_name string = 'acctest0001'

resource roleAssignments 'Microsoft.Authorization/roleAssignments@2022-04-01' = {
  name: '6faae21a-0cd6-4536-8c23-a278823d12ed'
  properties: {
    principalId: azurerm_user_assigned_identity.uai.principal_id
    principalType: 'ServicePrincipal'
    roleDefinitionId: data.azurerm_role_definition.roleAcrpull.id
  }
}

