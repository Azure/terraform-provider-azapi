param location string = 'westus'
param recommendation_id string = null
param resource_name string = 'acctest0001'

resource suppression 'Microsoft.Advisor/recommendations/suppressions@2023-01-01' = {
  name: resource_name
  properties: {
    suppressionId: ''
    ttl: '00:30:00'
  }
}

