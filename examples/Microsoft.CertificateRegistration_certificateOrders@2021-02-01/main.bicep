param location string = 'westeurope'
param resource_name string = 'acctest0001'

resource certificateOrder 'Microsoft.CertificateRegistration/certificateOrders@2021-02-01' = {
  location: 'global'
  name: resource_name
  properties: {
    autoRenew: true
    distinguishedName: 'CN=example.com'
    keySize: 2048
    productType: 'StandardDomainValidatedSsl'
    validityInYears: 1
  }
}

