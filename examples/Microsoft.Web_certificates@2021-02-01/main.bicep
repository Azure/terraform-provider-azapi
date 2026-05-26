param certificate_password string = null
param location string = 'westeurope'
param resource_name string = 'acctest0001'

resource certificate 'Microsoft.Web/certificates@2021-02-01' = {
  location: location
  name: resource_name
  properties: {
    password: certificate_password
    pfxBlob: filebase64("testdata/app_service_certificate.pfx")
  }
}

