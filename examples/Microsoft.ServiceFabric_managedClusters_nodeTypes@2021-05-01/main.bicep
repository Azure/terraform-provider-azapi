param admin_password string = null
param admin_username string = null
param location string = 'westeurope'
param resource_name string = 'acctest0001'

resource managedCluster 'Microsoft.ServiceFabric/managedClusters@2021-05-01' = {
  location: location
  name: resource_name
  properties: {
    addonFeatures: [
      'DnsService'
    ]
    adminPassword: admin_password
    adminUserName: admin_username
    clientConnectionPort: 12345
    clusterUpgradeCadence: 'Wave0'
    dnsName: resource_name
    httpGatewayConnectionPort: 23456
    loadBalancingRules: [
      {
        backendPort: 8000
        frontendPort: 443
        probeProtocol: 'http'
        probeRequestPath: '/'
        protocol: 'tcp'
      }
    ]
    networkSecurityRules: [
      {
        access: 'allow'
        destinationAddressPrefixes: [
          '0.0.0.0/0'
        ]
        destinationPortRanges: [
          '443'
        ]
        direction: 'inbound'
        name: 'rule443-allow-fe'
        priority: 1000
        protocol: 'tcp'
        sourceAddressPrefixes: [
          '0.0.0.0/0'
        ]
        sourcePortRanges: [
          '1-65535'
        ]
      }
    ]
  }
  sku: {
    name: 'Standard'
  }
  tags: {
    Test: 'value'
  }
}

resource nodeType 'Microsoft.ServiceFabric/managedClusters/nodeTypes@2021-05-01' = {
  parent: managedCluster
  name: resource_name
  properties: {
    applicationPorts: {
      endPort: 9000
      startPort: 7000
    }
    capacities: {}
    dataDiskSizeGB: 130
    dataDiskType: 'Standard_LRS'
    ephemeralPorts: {
      endPort: 20000
      startPort: 10000
    }
    isPrimary: true
    isStateless: false
    multiplePlacementGroups: false
    placementProperties: {}
    vmImageOffer: 'WindowsServer'
    vmImagePublisher: 'MicrosoftWindowsServer'
    vmImageSku: '2016-Datacenter'
    vmImageVersion: 'latest'
    vmInstanceCount: 5
    vmSecrets: []
    vmSize: 'Standard_DS2_v2'
  }
}

