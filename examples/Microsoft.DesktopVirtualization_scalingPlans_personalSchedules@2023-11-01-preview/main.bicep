param location string = 'westeurope'
param resource_name string = 'acctest0001'

resource hostPool 'Microsoft.DesktopVirtualization/hostPools@2023-09-05' = {
  location: location
  name: resource_name
  properties: {
    hostPoolType: 'Personal'
    loadBalancerType: 'Persistent'
    maxSessionLimit: 999999
    preferredAppGroupType: 'Desktop'
    publicNetworkAccess: 'Enabled'
    startVMOnConnect: true
    validationEnvironment: false
  }
}

resource personalSchedule 'Microsoft.DesktopVirtualization/scalingPlans/personalSchedules@2023-11-01-preview' = {
  parent: scalingPlan
  name: 'Weekdays'
  properties: {
    daysOfWeek: [
      'Monday'
      'Tuesday'
      'Wednesday'
      'Thursday'
      'Friday'
    ]
    offPeakActionOnDisconnect: 'Hibernate'
    offPeakActionOnLogoff: 'Hibernate'
    offPeakMinutesToWaitOnDisconnect: 20
    offPeakMinutesToWaitOnLogoff: 15
    offPeakStartTime: {
      hour: 17
      minute: 30
    }
    offPeakStartVMOnConnect: 'Enable'
    peakActionOnDisconnect: 'Hibernate'
    peakActionOnLogoff: 'Hibernate'
    peakMinutesToWaitOnDisconnect: 60
    peakMinutesToWaitOnLogoff: 60
    peakStartTime: {
      hour: 8
      minute: 0
    }
    peakStartVMOnConnect: 'Enable'
    rampDownActionOnDisconnect: 'Hibernate'
    rampDownActionOnLogoff: 'Hibernate'
    rampDownMinutesToWaitOnDisconnect: 45
    rampDownMinutesToWaitOnLogoff: 30
    rampDownStartTime: {
      hour: 16
      minute: 30
    }
    rampDownStartVMOnConnect: 'Enable'
    rampUpActionOnDisconnect: 'Hibernate'
    rampUpActionOnLogoff: 'Hibernate'
    rampUpAutoStartHosts: 'None'
    rampUpMinutesToWaitOnDisconnect: 45
    rampUpMinutesToWaitOnLogoff: 30
    rampUpStartTime: {
      hour: 7
      minute: 0
    }
    rampUpStartVMOnConnect: 'Enable'
  }
}

resource scalingPlan 'Microsoft.DesktopVirtualization/scalingPlans@2023-11-01-preview' = {
  location: location
  name: resource_name
  properties: {
    exclusionTag: 'no-schedule'
    hostPoolReferences: [
      {
        hostPoolArmPath: hostPool.id
        scalingPlanEnabled: true
      }
    ]
    hostPoolType: 'Personal'
    schedules: []
    timeZone: 'W. Europe Standard Time'
  }
}

