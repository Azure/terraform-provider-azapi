param certificate_data string = null
param certificate_thumbprint string = null
param location string = 'westeurope'
param resource_name string = 'acctest0001'

resource batchAccount 'Microsoft.Batch/batchAccounts@2022-10-01' = {
  location: location
  name: resource_name
  properties: {
    encryption: {
      keySource: 'Microsoft.Batch'
    }
    poolAllocationMode: 'BatchService'
    publicNetworkAccess: 'Enabled'
  }
}

resource certificate 'Microsoft.Batch/batchAccounts/certificates@2022-10-01' = {
  parent: batchAccount
  name: 'SHA1-${certificate_thumbprint}'
  properties: {
    data: certificate_data
    format: 'Cer'
    thumbprint: certificate_thumbprint
    thumbprintAlgorithm: 'sha1'
  }
}

