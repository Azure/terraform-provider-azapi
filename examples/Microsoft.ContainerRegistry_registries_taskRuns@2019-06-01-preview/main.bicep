param location string = 'westeurope'
param resource_name string = 'acctest0001'

resource registry 'Microsoft.ContainerRegistry/registries@2021-08-01-preview' = {
  location: location
  name: resource_name
  properties: {
    adminUserEnabled: false
    anonymousPullEnabled: false
    dataEndpointEnabled: false
    encryption: {
      status: 'disabled'
    }
    networkRuleBypassOptions: 'AzureServices'
    policies: {
      exportPolicy: {
        status: 'enabled'
      }
      quarantinePolicy: {
        status: 'disabled'
      }
      retentionPolicy: {
        status: 'disabled'
      }
      trustPolicy: {
        status: 'disabled'
      }
    }
    publicNetworkAccess: 'Enabled'
    zoneRedundancy: 'Disabled'
  }
  sku: {
    name: 'Standard'
    tier: 'Standard'
  }
}

resource taskRun 'Microsoft.ContainerRegistry/registries/taskRuns@2019-06-01-preview' = {
  parent: registry
  location: location
  name: resource_name
  properties: {
    runRequest: {
      dockerFilePath: 'Dockerfile'
      imageNames: [
        'helloworld:{{.Run.ID}}'
        'helloworld:latest'
      ]
      platform: {
        os: 'Linux'
      }
      sourceLocation: 'https://github.com/Azure-Samples/aci-helloworld.git#master'
      type: 'DockerBuildRequest'
    }
  }
}

