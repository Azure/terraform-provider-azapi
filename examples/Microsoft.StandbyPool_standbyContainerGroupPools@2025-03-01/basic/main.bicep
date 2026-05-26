param location string = 'eastus'
param resource_name string = 'acctest0001'

resource containerGroupProfile 'Microsoft.ContainerInstance/containerGroupProfiles@2024-05-01-preview' = {
  location: location
  name: '${resource_name}-contianerGroup'
  properties: {
    containers: [
      {
        name: 'mycontainergroupprofile'
        properties: {
          command: []
          environmentVariables: []
          image: 'mcr.microsoft.com/azuredocs/aci-helloworld:latest'
          ports: [
            {
              port: 8000
            }
          ]
          resources: {
            requests: {
              cpu: 1
              memoryInGB: json('1.5')
            }
          }
        }
      }
    ]
    imageRegistryCredentials: []
    ipAddress: {
      ports: [
        {
          port: 8000
          protocol: 'TCP'
        }
      ]
      type: 'Public'
    }
    osType: 'Linux'
    sku: 'Standard'
  }
}

resource standbyContainerGroupPool 'Microsoft.StandbyPool/standbyContainerGroupPools@2025-03-01' = {
  name: '${resource_name}-CGPool'
  location: location
  properties: {
    containerGroupProperties: {
      containerGroupProfile: {
        id: containerGroupProfile.id
        revision: 1
      }
      subnetIds: [
        {
          id: subnet.id
        }
      ]
    }
    elasticityProfile: {
      maxReadyCapacity: 5
      refillPolicy: 'always'
    }
    zones: [
      '1'
      '2'
      '3'
    ]
  }
}

resource subnet 'Microsoft.Network/virtualNetworks/subnets@2022-07-01' = {
  parent: virtualNetwork
  name: '${resource_name}-subnet'
  properties: {
    addressPrefix: '10.0.2.0/24'
    delegations: []
    privateEndpointNetworkPolicies: 'Enabled'
    privateLinkServiceNetworkPolicies: 'Enabled'
    serviceEndpointPolicies: []
    serviceEndpoints: []
  }
}

resource virtualNetwork 'Microsoft.Network/virtualNetworks@2022-07-01' = {
  location: location
  name: '${resource_name}-vnet'
  properties: {
    addressSpace: {
      addressPrefixes: [
        '10.0.0.0/16'
      ]
    }
    dhcpOptions: {
      dnsServers: []
    }
    subnets: []
  }
}

