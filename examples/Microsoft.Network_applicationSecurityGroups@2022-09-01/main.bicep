param location string = 'westeurope'
param resource_name string = 'acctest0001'

resource applicationSecurityGroup 'Microsoft.Network/applicationSecurityGroups@2022-09-01' = {
  location: location
  name: resource_name
}

