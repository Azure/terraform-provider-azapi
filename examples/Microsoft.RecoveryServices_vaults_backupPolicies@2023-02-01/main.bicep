param location string = 'westeurope'
param resource_name string = 'acctest0001'

resource backupPolicy 'Microsoft.RecoveryServices/vaults/backupPolicies@2023-02-01' = {
  parent: vault
  name: resource_name
  properties: {
    backupManagementType: 'AzureStorage'
    retentionPolicy: {
      dailySchedule: {
        retentionDuration: {
          count: 10
          durationType: 'Days'
        }
        retentionTimes: [
          '2018-07-30T23:00:00Z'
        ]
      }
      retentionPolicyType: 'LongTermRetentionPolicy'
    }
    schedulePolicy: {
      schedulePolicyType: 'SimpleSchedulePolicy'
      scheduleRunFrequency: 'Daily'
      scheduleRunTimes: [
        '2018-07-30T23:00:00Z'
      ]
    }
    timeZone: 'UTC'
    workLoadType: 'AzureFileShare'
  }
}

resource vault 'Microsoft.RecoveryServices/vaults@2022-10-01' = {
  location: location
  name: resource_name
  properties: {
    publicNetworkAccess: 'Enabled'
  }
  sku: {
    name: 'Standard'
  }
}

