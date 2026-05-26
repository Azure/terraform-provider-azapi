param location string = 'westus'
param resource_name string = 'acctest0001'

resource resourceManagementPrivateLink 'Microsoft.Authorization/resourceManagementPrivateLinks@2020-05-01' = {
  location: location
  name: resource_name
}

