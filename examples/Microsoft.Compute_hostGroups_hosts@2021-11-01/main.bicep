param location string = 'westeurope'
param resource_name string = 'acctest0001'

resource host 'Microsoft.Compute/hostGroups/hosts@2021-11-01' = {
  parent: hostGroup
  location: location
  name: resource_name
  properties: {
    autoReplaceOnFailure: true
    licenseType: 'None'
    platformFaultDomain: 1
  }
  sku: {
    name: 'DSv3-Type1'
  }
}

resource hostGroup 'Microsoft.Compute/hostGroups@2021-11-01' = {
  location: location
  name: resource_name
  properties: {
    platformFaultDomainCount: 2
  }
}

