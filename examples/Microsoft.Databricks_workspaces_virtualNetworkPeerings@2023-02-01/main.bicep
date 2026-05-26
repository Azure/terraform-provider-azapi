param location string = 'westeurope'
param resource_name string = 'acctest0001'

resource virtualNetwork 'Microsoft.Network/virtualNetworks@2022-07-01' = {
  location: location
  name: resource_name
  properties: {
    addressSpace: {
      addressPrefixes: [
        '10.0.1.0/24'
      ]
    }
    dhcpOptions: {
      dnsServers: []
    }
    subnets: []
  }
}

resource virtualNetworkPeering 'Microsoft.Databricks/workspaces/virtualNetworkPeerings@2023-02-01' = {
  parent: workspace
  name: resource_name
  properties: {
    allowForwardedTraffic: false
    allowGatewayTransit: false
    allowVirtualNetworkAccess: true
    databricksAddressSpace: {
      addressPrefixes: [
        '10.139.0.0/16'
      ]
    }
    remoteAddressSpace: {
      addressPrefixes: [
        '10.0.1.0/24'
      ]
    }
    remoteVirtualNetwork: {
      id: virtualNetwork.id
    }
    useRemoteGateways: false
  }
}

resource workspace 'Microsoft.Databricks/workspaces@2023-02-01' = {
  location: location
  name: resource_name
  properties: {
    managedResourceGroupId: data.azapi_resource_id.workspace_resource_group.id
    publicNetworkAccess: 'Enabled'
  }
  sku: {
    name: 'standard'
  }
}

