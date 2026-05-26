param billing_account_id string = null
param payment_method_id string = null
param payment_sca_id string = null
param resource_name string = 'acctest0001'

resource billingProfile 'Microsoft.Billing/billingAccounts/billingProfiles@2024-04-01' = {
  name: resource_name
  properties: {
    billTo: {
      addressLine1: 'TestWay'
      city: 'Redmond'
      companyName: 'TestCompany'
      country: 'US'
      isValidAddress: true
      postalCode: '12345-1234'
      region: 'WA'
    }
    displayName: resource_name
    enabledAzurePlans: [
      {
        skuId: '0001'
      }
    ]
    shipTo: {
      addressLine1: 'TestWay'
      city: 'Redmond'
      companyName: 'TestCompany'
      country: 'US'
      isValidAddress: true
      postalCode: '12345-1234'
      region: 'WA'
    }
  }
}

