param location string = 'westeurope'
param resource_name string = 'acctest0001'

resource customDomain 'Microsoft.Cdn/profiles/customDomains@2021-06-01' = {
  parent: profile
  name: resource_name
  properties: {
    azureDnsZone: {
      id: dnsZone.id
    }
    hostName: 'fabrikam.${resource_name}.com'
    tlsSettings: {
      certificateType: 'ManagedCertificate'
      minimumTlsVersion: 'TLS12'
    }
  }
}

resource dnsZone 'Microsoft.Network/dnsZones@2018-05-01' = {
  location: 'global'
  name: '${resource_name}.com'
}

resource profile 'Microsoft.Cdn/profiles@2021-06-01' = {
  location: 'global'
  name: resource_name
  properties: {
    originResponseTimeoutSeconds: 120
  }
  sku: {
    name: 'Premium_AzureFrontDoor'
  }
}

