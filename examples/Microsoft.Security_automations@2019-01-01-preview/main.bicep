param location string = 'westeurope'
param resource_name string = 'acctest0001'

resource automation 'Microsoft.Security/automations@2019-01-01-preview' = {
  location: resourceGroup().location
  name: 'ExportToWorkspace'
  properties: {
    actions: [
      {
        actionType: 'Workspace'
        workspaceResourceId: workspace.id
      }
    ]
    isEnabled: true
    scopes: [
      {
        description: 'Security Export for the subscription'
        scopePath: resourceGroup.id
      }
    ]
    sources: [
      {
        eventSource: 'Assessments'
        ruleSets: [
          {
            rules: [
              {
                expectedValue: 'Microsoft.Security/assessments'
                operator: 'Contains'
                propertyJPath: 'type'
                propertyType: 'String'
              }
            ]
          }
        ]
      }
      {
        eventSource: 'AssessmentsSnapshot'
        ruleSets: [
          {
            rules: [
              {
                expectedValue: 'Microsoft.Security/assessments'
                operator: 'Contains'
                propertyJPath: 'type'
                propertyType: 'String'
              }
            ]
          }
        ]
      }
      {
        eventSource: 'SubAssessments'
      }
      {
        eventSource: 'SubAssessmentsSnapshot'
      }
      {
        eventSource: 'Alerts'
        ruleSets: [
          {
            rules: [
              {
                expectedValue: 'low'
                operator: 'Equals'
                propertyJPath: 'Severity'
                propertyType: 'String'
              }
            ]
          }
          {
            rules: [
              {
                expectedValue: 'medium'
                operator: 'Equals'
                propertyJPath: 'Severity'
                propertyType: 'String'
              }
            ]
          }
          {
            rules: [
              {
                expectedValue: 'high'
                operator: 'Equals'
                propertyJPath: 'Severity'
                propertyType: 'String'
              }
            ]
          }
          {
            rules: [
              {
                expectedValue: 'informational'
                operator: 'Equals'
                propertyJPath: 'Severity'
                propertyType: 'String'
              }
            ]
          }
        ]
      }
      {
        eventSource: 'SecureScores'
      }
      {
        eventSource: 'SecureScoresSnapshot'
      }
      {
        eventSource: 'SecureScoreControls'
      }
      {
        eventSource: 'SecureScoreControlsSnapshot'
      }
      {
        eventSource: 'RegulatoryComplianceAssessment'
      }
      {
        eventSource: 'RegulatoryComplianceAssessmentSnapshot'
      }
    ]
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

