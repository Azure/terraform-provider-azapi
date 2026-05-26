param location string = 'westeurope'
param resource_name string = 'acctest0001'

resource fluidRelayServer 'Microsoft.FluidRelay/fluidRelayServers@2022-05-26' = {
  location: location
  name: resource_name
  properties: {}
  tags: {
    foo: 'bar'
  }
}

