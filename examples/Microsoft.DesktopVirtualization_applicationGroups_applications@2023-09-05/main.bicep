param location string = 'westeurope'
param resource_name string = 'acctest0001'

resource application 'Microsoft.DesktopVirtualization/applicationGroups/applications@2023-09-05' = {
  parent: applicationGroup
  location: location
  name: resource_name
  properties: {
    commandLineSetting: 'DoNotAllow'
    filePath: 'C:\\Program Files\\Google\\Chrome\\Application\\chrome.exe'
    showInPortal: false
  }
}

resource applicationGroup 'Microsoft.DesktopVirtualization/applicationGroups@2023-09-05' = {
  location: location
  name: resource_name
  properties: {
    applicationGroupType: 'RemoteApp'
    hostPoolArmPath: hostPool.id
  }
}

resource hostPool 'Microsoft.DesktopVirtualization/hostPools@2023-09-05' = {
  location: location
  name: resource_name
  properties: {
    hostPoolType: 'Pooled'
    loadBalancerType: 'BreadthFirst'
    maxSessionLimit: 999999
    preferredAppGroupType: 'Desktop'
    publicNetworkAccess: 'Enabled'
    startVMOnConnect: false
    validationEnvironment: false
  }
}

