param location string = 'westus'
param resource_name string = 'acctest0001'

resource job 'Microsoft.App/jobs@2025-01-01' = {
  location: location
  name: '${resource_name}-cajob'
  properties: {
    configuration: {
      manualTriggerConfig: {
        parallelism: 4
        replicaCompletionCount: 1
      }
      replicaRetryLimit: 10
      replicaTimeout: 10
      triggerType: 'Manual'
    }
    environmentId: managedEnvironment.id
    template: {
      containers: [
        {
          env: []
          image: 'jackofallops/azure-containerapps-python-acctest:v0.0.1'
          name: 'testcontainerappsjob0'
          probes: []
          resources: {
            cpu: json('0.5')
            memory: '1Gi'
          }
          volumeMounts: []
        }
      ]
      initContainers: []
      volumes: []
    }
  }
}

resource managedEnvironment 'Microsoft.App/managedEnvironments@2025-01-01' = {
  location: location
  name: '${resource_name}-env'
  properties: {
    appLogsConfiguration: {
      destination: 'log-analytics'
      logAnalyticsConfiguration: {
        customerId: workspace.properties.customerId
        sharedKey: data.azapi_resource_action.workspace_keys.output.primarySharedKey
      }
    }
  }
}

resource workspace 'Microsoft.OperationalInsights/workspaces@2023-09-01' = {
  location: location
  name: '${resource_name}-law'
  properties: {
    retentionInDays: 30
    sku: {
      name: 'PerGB2018'
    }
  }
}

