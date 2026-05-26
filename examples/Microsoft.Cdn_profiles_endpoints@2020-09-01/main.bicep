param location string = 'westeurope'
param resource_name string = 'acctest0001'

resource endpoint 'Microsoft.Cdn/profiles/endpoints@2020-09-01' = {
  parent: profile
  location: location
  name: resource_name
  properties: {
    isHttpAllowed: true
    isHttpsAllowed: true
    origins: [
      {
        name: 'acceptanceTestCdnOrigin1'
        properties: {
          hostName: 'www.contoso.com'
          httpPort: 80
          httpsPort: 443
        }
      }
    ]
    queryStringCachingBehavior: 'IgnoreQueryString'
  }
}

resource profile 'Microsoft.Cdn/profiles@2020-09-01' = {
  location: location
  name: resource_name
  sku: {
    name: 'Standard_Verizon'
  }
}

