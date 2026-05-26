param location string = 'eastus2'
param resource_name string = 'acctest0001'

resource netAppAccount 'Microsoft.NetApp/netAppAccounts@2022-05-01' = {
  location: location
  name: resource_name
  properties: {
    activeDirectories: []
  }
}

resource snapshotPolicy 'Microsoft.NetApp/netAppAccounts/snapshotPolicies@2022-05-01' = {
  parent: netAppAccount
  location: location
  name: resource_name
  properties: {
    dailySchedule: {
      hour: 22
      minute: 15
      snapshotsToKeep: 1
    }
    enabled: true
    hourlySchedule: {
      minute: 15
      snapshotsToKeep: 1
    }
    monthlySchedule: {
      daysOfMonth: '30,15,1'
      hour: 5
      minute: 0
      snapshotsToKeep: 1
    }
    weeklySchedule: {
      day: 'Monday,Friday'
      hour: 23
      minute: 0
      snapshotsToKeep: 1
    }
  }
}

