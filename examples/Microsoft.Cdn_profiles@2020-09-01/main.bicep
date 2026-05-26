param location string = 'westeurope'
param resource_name string = 'acctest0001'

resource profile 'Microsoft.Cdn/profiles@2020-09-01' = {
  location: location
  name: resource_name
  sku: {
    name: 'Standard_Verizon'
  }
}

