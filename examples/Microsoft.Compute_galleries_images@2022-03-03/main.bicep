param location string = 'westeurope'
param resource_name string = 'acctest0001'

resource gallery 'Microsoft.Compute/galleries@2022-03-03' = {
  location: location
  name: resource_name
  properties: {
    description: ''
  }
}

resource image 'Microsoft.Compute/galleries/images@2022-03-03' = {
  parent: gallery
  location: location
  name: resource_name
  properties: {
    architecture: 'x64'
    description: ''
    disallowed: {
      diskTypes: []
    }
    features: null
    hyperVGeneration: 'V1'
    identifier: {
      offer: 'AccTesOffer230630032848825313'
      publisher: 'AccTesPublisher230630032848825313'
      sku: 'AccTesSku230630032848825313'
    }
    osState: 'Generalized'
    osType: 'Linux'
    privacyStatementUri: ''
    recommended: {
      memory: {}
      vCPUs: {}
    }
    releaseNoteUri: ''
  }
}

