param location string = 'westeurope'
param resource_name string = 'acctest0001'

resource systemTopic 'Microsoft.EventGrid/systemTopics@2021-12-01' = {
  location: 'global'
  name: resource_name
  properties: {
    source: resourceGroup.id
    topicType: 'Microsoft.Resources.ResourceGroups'
  }
}

