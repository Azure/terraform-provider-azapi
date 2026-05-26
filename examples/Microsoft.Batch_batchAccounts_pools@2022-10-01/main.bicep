param location string = 'westeurope'
param resource_name string = 'acctest0001'

resource batchAccount 'Microsoft.Batch/batchAccounts@2022-10-01' = {
  location: location
  name: resource_name
  properties: {
    encryption: {
      keySource: 'Microsoft.Batch'
    }
    poolAllocationMode: 'BatchService'
    publicNetworkAccess: 'Enabled'
  }
}

resource pool 'Microsoft.Batch/batchAccounts/pools@2022-10-01' = {
  parent: batchAccount
  name: resource_name
  properties: {
    certificates: null
    deploymentConfiguration: {
      virtualMachineConfiguration: {
        imageReference: {
          offer: 'UbuntuServer'
          publisher: 'Canonical'
          sku: '18.04-lts'
          version: 'latest'
        }
        nodeAgentSkuId: 'batch.node.ubuntu 18.04'
        osDisk: {
          ephemeralOSDiskSettings: {
            placement: ''
          }
        }
      }
    }
    displayName: ''
    interNodeCommunication: 'Enabled'
    metadata: []
    scaleSettings: {
      fixedScale: {
        nodeDeallocationOption: ''
        resizeTimeout: 'PT15M'
        targetDedicatedNodes: 1
        targetLowPriorityNodes: 0
      }
    }
    taskSlotsPerNode: 1
    vmSize: 'STANDARD_A1'
  }
}

