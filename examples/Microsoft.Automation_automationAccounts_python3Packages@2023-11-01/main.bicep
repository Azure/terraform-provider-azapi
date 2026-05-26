param location string = 'westus'
param resource_name string = 'acctest0001'

resource automationAccount 'Microsoft.Automation/automationAccounts@2023-11-01' = {
  location: location
  name: resource_name
  properties: {
    disableLocalAuth: false
    encryption: {
      keySource: 'Microsoft.Automation'
    }
    publicNetworkAccess: true
    sku: {
      name: 'Basic'
    }
  }
}

resource python3Package 'Microsoft.Automation/automationAccounts/python3Packages@2023-11-01' = {
  parent: automationAccount
  name: resource_name
  tags: {
    key: 'foo'
  }
  properties: {
    contentLink: {
      uri: 'https://files.pythonhosted.org/packages/py3/r/requests/requests-2.31.0-py3-none-any.whl'
      version: '2.31.0'
    }
  }
}

