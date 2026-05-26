param location string = 'westus'
param resource_name string = 'acctest0001'

resource policyDefinition 'Microsoft.Authorization/policyDefinitions@2021-06-01' = {
  name: resource_name
  properties: {
    description: ''
    displayName: 'my-policy-definition'
    mode: 'All'
    parameters: {
      allowedLocations: {
        metadata: {
          description: 'The list of allowed locations for resources.'
          displayName: 'Allowed locations'
          strongType: 'location'
        }
        type: 'Array'
      }
    }
    policyRule: {
      if: {
        not: {
          field: 'location'
          in: '[parameters(\'allowedLocations\')]'
        }
      }
      then: {
        effect: 'audit'
      }
    }
    policyType: 'Custom'
  }
}

resource policySetDefinition 'Microsoft.Authorization/policySetDefinitions@2025-01-01' = {
  name: 'acctestpolset-${resource_name}'
  properties: {
    description: ''
    displayName: 'acctestpolset-${resource_name}'
    parameters: {
      allowedLocations: {
        metadata: {
          description: 'The list of allowed locations for resources.'
          displayName: 'Allowed locations'
          strongType: 'location'
        }
        type: 'Array'
      }
    }
    policyDefinitions: [
      {
        groupNames: []
        parameters: {
          listOfAllowedLocations: {
            value: '[parameters(\'allowedLocations\')]'
          }
        }
        policyDefinitionId: policyDefinition.id
        policyDefinitionReferenceId: ''
      }
    ]
    policyType: 'Custom'
  }
}

