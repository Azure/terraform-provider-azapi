param location string = 'eastus'
param resource_name string = 'acctest0001'

resource policyAssignment 'Microsoft.Authorization/policyAssignments@2022-06-01' = {
  name: resource_name
  properties: {
    displayName: ''
    enforcementMode: 'Default'
    parameters: {
      listOfAllowedLocations: {
        value: [
          'West Europe'
          'West US 2'
          'East US 2'
        ]
      }
    }
    policyDefinitionId: '/providers/Microsoft.Authorization/policyDefinitions/e56962a6-4747-49cd-b67b-bf8b01975c4c'
    scope: data.azapi_resource.subscription.id
  }
}

resource remediation 'Microsoft.PolicyInsights/remediations@2021-10-01' = {
  name: resource_name
  properties: {
    filters: {
      locations: []
    }
    policyAssignmentId: policyAssignment.id
    policyDefinitionReferenceId: ''
    resourceDiscoveryMode: 'ExistingNonCompliant'
  }
}

