param location string = 'westeurope'
param resource_name string = 'acctest0001'

resource FrontDoorWebApplicationFirewallPolicy 'Microsoft.Network/FrontDoorWebApplicationFirewallPolicies@2020-11-01' = {
  location: 'global'
  name: resource_name
  properties: {
    customRules: {
      rules: [
        {
          action: 'Block'
          enabledState: 'Enabled'
          matchConditions: [
            {
              matchValue: [
                '192.168.1.0/24'
                '10.0.0.0/24'
              ]
              matchVariable: 'RemoteAddr'
              negateCondition: false
              operator: 'IPMatch'
            }
          ]
          name: 'Rule1'
          priority: 1
          rateLimitDurationInMinutes: 1
          rateLimitThreshold: 10
          ruleType: 'MatchRule'
        }
      ]
    }
    managedRules: {
      managedRuleSets: [
        {
          ruleGroupOverrides: [
            {
              ruleGroupName: 'PHP'
              rules: [
                {
                  action: 'Block'
                  enabledState: 'Disabled'
                  ruleId: '933111'
                }
              ]
            }
          ]
          ruleSetAction: 'Block'
          ruleSetType: 'DefaultRuleSet'
          ruleSetVersion: 'preview-0.1'
        }
        {
          ruleSetAction: 'Block'
          ruleSetType: 'BotProtection'
          ruleSetVersion: 'preview-0.1'
        }
      ]
    }
    policySettings: {
      customBlockResponseBody: 'PGh0bWw+CjxoZWFkZXI+PHRpdGxlPkhlbGxvPC90aXRsZT48L2hlYWRlcj4KPGJvZHk+CkhlbGxvIHdvcmxkCjwvYm9keT4KPC9odG1sPg=='
      customBlockResponseStatusCode: 403
      enabledState: 'Enabled'
      mode: 'Prevention'
      redirectUrl: 'https://www.fabrikam.com'
    }
  }
  sku: {
    name: 'Premium_AzureFrontDoor'
  }
}

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

resource securityPolicy 'Microsoft.Cdn/profiles/securityPolicies@2021-06-01' = {
  parent: profile
  name: resource_name
  properties: {
    parameters: {
      associations: [
        {
          domains: [
            {
              id: customDomain.id
            }
          ]
          patternsToMatch: [
            '/*'
          ]
        }
      ]
      type: 'WebApplicationFirewall'
      wafPolicy: {
        id: FrontDoorWebApplicationFirewallPolicy.id
      }
    }
  }
}

