param location string = 'westeurope'
param resource_name string = 'acctest0001'

resource containerApp 'Microsoft.App/containerApps@2022-03-01' = {
  location: location
  name: resource_name
  properties: {
    configuration: {
      activeRevisionsMode: 'Single'
    }
    managedEnvironmentId: managedEnvironment.id
    template: {
      containers: [
        {
          env: []
          image: 'jackofallops/azure-containerapps-python-acctest:v0.0.1'
          name: 'acctest-cont-230630032906865620'
          probes: []
          resources: {
            cpu: json('0.25')
            ephemeralStorage: '1Gi'
            memory: '0.5Gi'
          }
          volumeMounts: []
        }
      ]
      scale: {
        maxReplicas: 10
      }
      volumes: []
    }
  }
}

resource managedEnvironment 'Microsoft.App/managedEnvironments@2022-03-01' = {
  location: location
  name: resource_name
  properties: {
    appLogsConfiguration: {
      destination: 'log-analytics'
      logAnalyticsConfiguration: {
        customerId: workspace.properties.customerId
        sharedKey: data.azapi_resource_action.sharedKeys.output.primarySharedKey
      }
    }
    vnetConfiguration: {}
  }
}

resource workspace 'Microsoft.OperationalInsights/workspaces@2022-10-01' = {
  location: location
  name: resource_name
  properties: {
    features: {
      disableLocalAuth: false
      enableLogAccessUsingOnlyResourcePermissions: true
    }
    publicNetworkAccessForIngestion: 'Enabled'
    publicNetworkAccessForQuery: 'Enabled'
    retentionInDays: 30
    sku: {
      name: 'PerGB2018'
    }
    workspaceCapping: {
      dailyQuotaGb: -1
    }
  }
}

