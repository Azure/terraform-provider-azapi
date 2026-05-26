param location string = 'westeurope'
param resource_name string = 'acctest0001'

resource deployment 'Microsoft.Resources/deployments@2020-06-01' = {
  name: resource_name
  properties: {
    mode: 'Complete'
    template: {
      '$schema': 'https://schema.management.azure.com/schemas/2015-01-01/deploymentTemplate.json#'
      contentVersion: '1.0.0.0'
      parameters: {
        storageAccountType: {
          allowedValues: [
            'Standard_LRS'
            'Standard_GRS'
            'Standard_ZRS'
          ]
          defaultValue: 'Standard_LRS'
          metadata: {
            description: 'Storage Account type'
          }
          type: 'string'
        }
      }
      resources: [
        {
          apiVersion: '[variables(\'apiVersion\')]'
          location: '[variables(\'location\')]'
          name: '[variables(\'storageAccountName\')]'
          properties: {
            accountType: '[parameters(\'storageAccountType\')]'
          }
          type: 'Microsoft.Storage/storageAccounts'
        }
        {
          apiVersion: '[variables(\'apiVersion\')]'
          location: '[variables(\'location\')]'
          name: '[variables(\'publicIPAddressName\')]'
          properties: {
            dnsSettings: {
              domainNameLabel: '[variables(\'dnsLabelPrefix\')]'
            }
            publicIPAllocationMethod: '[variables(\'publicIPAddressType\')]'
          }
          type: 'Microsoft.Network/publicIPAddresses'
        }
      ]
      variables: {
        apiVersion: '2015-06-15'
        dnsLabelPrefix: '[concat(\'terraform-tdacctest\', uniquestring(resourceGroup().id))]'
        location: '[resourceGroup().location]'
        publicIPAddressName: '[concat(\'myPublicIp\', uniquestring(resourceGroup().id))]'
        publicIPAddressType: 'Dynamic'
        storageAccountName: '[concat(uniquestring(resourceGroup().id), \'storage\')]'
      }
    }
  }
}

