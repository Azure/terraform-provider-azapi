param location string = 'westeurope'
param resource_name string = 'acctest0001'

resource lab 'Microsoft.DevTestLab/labs@2018-09-15' = {
  location: location
  name: resource_name
  properties: {
    labStorageType: 'Premium'
  }
}

resource schedule 'Microsoft.DevTestLab/labs/schedules@2018-09-15' = {
  parent: lab
  location: location
  name: 'LabVmsShutdown'
  properties: {
    dailyRecurrence: {
      time: '0100'
    }
    notificationSettings: {
      status: 'Disabled'
      timeInMinutes: 0
      webhookUrl: ''
    }
    status: 'Disabled'
    taskType: 'LabVmsShutdownTask'
    timeZoneId: 'India Standard Time'
  }
  tags: {
    environment: 'Production'
  }
}

