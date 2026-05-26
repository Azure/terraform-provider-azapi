param location string = 'westus'
param resource_name string = 'acctest0001'

resource fleet 'Microsoft.ContainerService/fleets@2024-04-01' = {
  location: location
  name: resource_name
  properties: {}
}

resource managedCluster 'Microsoft.ContainerService/managedClusters@2025-02-01' = {
  identity: [
    {
      identity_ids: []
      type: 'SystemAssigned'
    }
  ]
  location: location
  name: resource_name
  properties: {
    addonProfiles: {}
    agentPoolProfiles: [
      {
        count: 1
        enableAutoScaling: false
        enableEncryptionAtHost: false
        enableFIPS: false
        enableNodePublicIP: false
        enableUltraSSD: false
        kubeletDiskType: ''
        mode: 'System'
        name: 'default'
        nodeLabels: {}
        osDiskType: 'Managed'
        osType: 'Linux'
        scaleDownMode: 'Delete'
        tags: {}
        type: 'VirtualMachineScaleSets'
        upgradeSettings: {
          drainTimeoutInMinutes: 0
          maxSurge: '10%'
          nodeSoakDurationInMinutes: 0
        }
        vmSize: 'Standard_B2s'
      }
    ]
    apiServerAccessProfile: {
      disableRunCommand: false
      enablePrivateCluster: false
      enablePrivateClusterPublicFQDN: false
    }
    autoUpgradeProfile: {
      nodeOSUpgradeChannel: 'NodeImage'
      upgradeChannel: 'none'
    }
    azureMonitorProfile: {
      metrics: {
        enabled: false
      }
    }
    disableLocalAccounts: false
    dnsPrefix: resource_name
    enableRBAC: true
    kubernetesVersion: ''
    metricsProfile: {
      costAnalysis: {
        enabled: false
      }
    }
    nodeResourceGroup: ''
    securityProfile: {}
    servicePrincipalProfile: {
      clientId: 'msi'
    }
    supportPlan: 'KubernetesOfficial'
  }
  sku: {
    name: 'Base'
    tier: 'Free'
  }
}

resource member 'Microsoft.ContainerService/fleets/members@2024-04-01' = {
  parent: fleet
  name: resource_name
  properties: {
    clusterResourceId: managedCluster.id
    group: 'default'
  }
}

