param location string = 'eastus'
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

