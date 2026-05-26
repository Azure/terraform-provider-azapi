param admin_password string = null
param location string = 'westeurope'
param resource_name string = 'acctest0001'

resource lab 'Microsoft.LabServices/labs@2022-08-01' = {
  location: location
  name: resource_name
  properties: {
    autoShutdownProfile: {
      shutdownOnDisconnect: 'Disabled'
      shutdownOnIdle: 'None'
      shutdownWhenNotConnected: 'Disabled'
    }
    connectionProfile: {
      clientRdpAccess: 'None'
      clientSshAccess: 'None'
      webRdpAccess: 'None'
      webSshAccess: 'None'
    }
    securityProfile: {
      openAccess: 'Disabled'
    }
    title: 'Test Title'
    virtualMachineProfile: {
      additionalCapabilities: {
        installGpuDrivers: 'Disabled'
      }
      adminUser: {
        password: admin_password
        username: 'testadmin'
      }
      createOption: 'Image'
      imageReference: {
        offer: '0001-com-ubuntu-server-focal'
        publisher: 'canonical'
        sku: '20_04-lts'
        version: 'latest'
      }
      sku: {
        capacity: 1
        name: 'Classic_Fsv2_2_4GB_128_S_SSD'
      }
      usageQuota: 'PT0S'
      useSharedPassword: 'Disabled'
    }
  }
}

resource schedule 'Microsoft.LabServices/labs/schedules@2022-08-01' = {
  parent: lab
  name: resource_name
  properties: {
    stopAt: '2023-06-30T04:33:55Z'
    timeZoneId: 'America/Los_Angeles'
  }
}

