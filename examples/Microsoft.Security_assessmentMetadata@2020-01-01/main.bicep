param location string = 'eastus'
param resource_name string = 'acctest0001'

resource assessmentMetadatum 'Microsoft.Security/assessmentMetadata@2020-01-01' = {
  name: '95c7a001-d595-43af-9754-1310c740d34c'
  properties: {
    assessmentType: 'CustomerManaged'
    description: 'Test Description'
    displayName: 'Test Display Name'
    severity: 'Medium'
  }
}

