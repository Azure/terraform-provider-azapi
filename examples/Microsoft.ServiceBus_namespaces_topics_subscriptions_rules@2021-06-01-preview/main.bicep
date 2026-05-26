param location string = 'westeurope'
param resource_name string = 'acctest0001'

resource namespace 'Microsoft.ServiceBus/namespaces@2022-01-01-preview' = {
  location: location
  name: resource_name
  properties: {
    disableLocalAuth: false
    publicNetworkAccess: 'Enabled'
    zoneRedundant: false
  }
  sku: {
    capacity: 0
    name: 'Standard'
    tier: 'Standard'
  }
}

resource rule 'Microsoft.ServiceBus/namespaces/topics/subscriptions/rules@2021-06-01-preview' = {
  parent: subscription
  name: resource_name
  properties: {
    correlationFilter: {
      contentType: 'test_content_type'
      correlationId: 'test_correlation_id'
      label: 'test_label'
      messageId: 'test_message_id'
      replyTo: 'test_reply_to'
      replyToSessionId: 'test_reply_to_session_id'
      sessionId: 'test_session_id'
      to: 'test_to'
    }
    filterType: 'CorrelationFilter'
  }
}

resource subscription 'Microsoft.ServiceBus/namespaces/topics/subscriptions@2021-06-01-preview' = {
  parent: topic
  name: resource_name
  properties: {
    clientAffineProperties: {}
    deadLetteringOnFilterEvaluationExceptions: true
    deadLetteringOnMessageExpiration: false
    enableBatchedOperations: false
    isClientAffine: false
    maxDeliveryCount: 10
    requiresSession: false
    status: 'Active'
  }
}

resource topic 'Microsoft.ServiceBus/namespaces/topics@2021-06-01-preview' = {
  parent: namespace
  name: resource_name
  properties: {
    enableBatchedOperations: false
    enableExpress: false
    enablePartitioning: false
    maxSizeInMegabytes: 5120
    requiresDuplicateDetection: false
    status: 'Active'
    supportOrdering: false
  }
}

