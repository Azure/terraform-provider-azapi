param location string = 'westeurope'
param resource_name string = 'acctest0001'

resource ApplicationGatewayWebApplicationFirewallPolicy 'Microsoft.Network/ApplicationGatewayWebApplicationFirewallPolicies@2022-07-01' = {
  location: location
  name: resource_name
  properties: {
    customRules: []
    managedRules: {
      exclusions: []
      managedRuleSets: [
        {
          ruleGroupOverrides: []
          ruleSetType: 'OWASP'
          ruleSetVersion: '3.1'
        }
      ]
    }
    policySettings: {
      fileUploadLimitInMb: 100
      maxRequestBodySizeInKb: 128
      mode: 'Detection'
      requestBodyCheck: true
      state: 'Enabled'
    }
  }
}

