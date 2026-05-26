param location string = 'westeurope'
param resource_name string = 'acctest0001'

resource domain 'Microsoft.Communication/emailServices/domains@2023-04-01-preview' = {
  parent: emailService
  location: 'global'
  name: 'example.com'
  tags: {
    env: 'Test'
  }
  properties: {
    domainManagement: 'CustomerManaged'
    userEngagementTracking: 'Disabled'
  }
}

resource emailService 'Microsoft.Communication/emailServices@2023-04-01-preview' = {
  location: 'global'
  name: resource_name
  properties: {
    dataLocation: 'United States'
  }
}

resource senderUsername 'Microsoft.Communication/emailServices/domains/senderUsernames@2023-04-01-preview' = {
  parent: domain
  name: 'TestSenderUserName'
  properties: {
    displayName: 'TestDisplayName'
    username: 'TestSenderUserName'
  }
}

