param cdn_location string = 'global'
param location string = 'westus'
param resource_name string = 'acctest0001'

resource origin 'Microsoft.Cdn/profiles/originGroups/origins@2024-09-01' = {
  parent: originGroup
  name: '${resource_name}-origin'
  properties: {
    enabledState: 'Enabled'
    enforceCertificateNameCheck: false
    hostName: 'contoso.com'
    httpPort: 80
    httpsPort: 443
    originHostHeader: 'www.contoso.com'
    priority: 1
    weight: 1
  }
}

resource originGroup 'Microsoft.Cdn/profiles/originGroups@2024-09-01' = {
  parent: profile
  name: '${resource_name}-origingroup'
  properties: {
    loadBalancingSettings: {
      additionalLatencyInMilliseconds: 0
      sampleSize: 16
      successfulSamplesRequired: 3
    }
    sessionAffinityState: 'Enabled'
    trafficRestorationTimeToHealedOrNewEndpointsInMinutes: 10
  }
}

resource profile 'Microsoft.Cdn/profiles@2024-09-01' = {
  location: cdn_location
  name: '${resource_name}-profile'
  properties: {
    originResponseTimeoutSeconds: 120
  }
  sku: {
    name: 'Standard_AzureFrontDoor'
  }
}

resource rule 'Microsoft.Cdn/profiles/ruleSets/rules@2024-09-01' = {
  parent: ruleSet
  name: 'rule${substr(var.resource_name, -4, -1)}'
  properties: {
    actions: [
      {
        name: 'RouteConfigurationOverride'
        parameters: {
          cacheConfiguration: {
            cacheBehavior: 'OverrideIfOriginMissing'
            cacheDuration: '23:59:59'
            isCompressionEnabled: 'Disabled'
            queryParameters: 'clientIp={client_ip}'
            queryStringCachingBehavior: 'IgnoreSpecifiedQueryStrings'
          }
          originGroupOverride: {
            forwardingProtocol: 'HttpsOnly'
            originGroup: {
              id: originGroup.id
            }
          }
          typeName: 'DeliveryRuleRouteConfigurationOverrideActionParameters'
        }
      }
    ]
    conditions: []
    matchProcessingBehavior: 'Continue'
    order: 1
  }
}

resource ruleSet 'Microsoft.Cdn/profiles/ruleSets@2024-09-01' = {
  parent: profile
  name: 'ruleSet${substr(var.resource_name, -4, -1)}'
}

