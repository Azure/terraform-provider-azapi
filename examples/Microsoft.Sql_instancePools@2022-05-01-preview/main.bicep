param location string = 'westeurope'
param resource_name string = 'acctest0001'

resource instancePool 'Microsoft.Sql/instancePools@2022-05-01-preview' = {
  location: resourceGroup().location
  name: resource_name
  properties: {
    licenseType: 'LicenseIncluded'
    subnetId: data.azapi_resource.subnet.id
    vCores: 8
  }
  sku: {
    family: 'Gen5'
    name: 'GP_Gen5'
    tier: 'GeneralPurpose'
  }
}

resource networkSecurityGroup 'Microsoft.Network/networkSecurityGroups@2023-04-01' = {
  location: resourceGroup().location
  name: resource_name
  properties: {
    securityRules: [
      {
        name: 'allow_tds_inbound'
        properties: {
          access: 'Allow'
          description: 'Allow access to data'
          destinationAddressPrefix: '*'
          destinationPortRange: '1433'
          direction: 'Inbound'
          priority: 1000
          protocol: 'TCP'
          sourceAddressPrefix: 'VirtualNetwork'
          sourcePortRange: '*'
        }
      }
      {
        name: 'allow_redirect_inbound'
        properties: {
          access: 'Allow'
          description: 'Allow inbound redirect traffic to Managed Instance inside the virtual network'
          destinationAddressPrefix: '*'
          destinationPortRange: '11000-11999'
          direction: 'Inbound'
          priority: 1100
          protocol: 'Tcp'
          sourceAddressPrefix: 'VirtualNetwork'
          sourcePortRange: '*'
        }
      }
      {
        name: 'allow_geodr_inbound'
        properties: {
          access: 'Allow'
          description: 'Allow inbound geodr traffic inside the virtual network'
          destinationAddressPrefix: '*'
          destinationPortRange: '5022'
          direction: 'Inbound'
          priority: 1200
          protocol: 'Tcp'
          sourceAddressPrefix: 'VirtualNetwork'
          sourcePortRange: '*'
        }
      }
      {
        name: 'deny_all_inbound'
        properties: {
          access: 'Deny'
          description: 'Deny all other inbound traffic'
          destinationAddressPrefix: '*'
          destinationPortRange: '*'
          direction: 'Inbound'
          priority: 4096
          protocol: '*'
          sourceAddressPrefix: '*'
          sourcePortRange: '*'
        }
      }
      {
        name: 'allow_linkedserver_outbound'
        properties: {
          access: 'Allow'
          description: 'Allow outbound linkedserver traffic inside the virtual network'
          destinationAddressPrefix: 'VirtualNetwork'
          destinationPortRange: '1433'
          direction: 'Outbound'
          priority: 1000
          protocol: 'Tcp'
          sourceAddressPrefix: '*'
          sourcePortRange: '*'
        }
      }
      {
        name: 'allow_redirect_outbound'
        properties: {
          access: 'Allow'
          description: 'Allow outbound redirect traffic to Managed Instance inside the virtual network'
          destinationAddressPrefix: 'VirtualNetwork'
          destinationPortRange: '11000-11999'
          direction: 'Outbound'
          priority: 1100
          protocol: 'Tcp'
          sourceAddressPrefix: '*'
          sourcePortRange: '*'
        }
      }
      {
        name: 'allow_geodr_outbound'
        properties: {
          access: 'Allow'
          description: 'Allow outbound geodr traffic inside the virtual network'
          destinationAddressPrefix: 'VirtualNetwork'
          destinationPortRange: '5022'
          direction: 'Outbound'
          priority: 1200
          protocol: 'Tcp'
          sourceAddressPrefix: '*'
          sourcePortRange: '*'
        }
      }
      {
        name: 'deny_all_outbound'
        properties: {
          access: 'Deny'
          description: 'Deny all other outbound traffic'
          destinationAddressPrefix: '*'
          destinationPortRange: '*'
          direction: 'Outbound'
          priority: 4096
          protocol: '*'
          sourceAddressPrefix: '*'
          sourcePortRange: '*'
        }
      }
    ]
  }
}

resource routeTable 'Microsoft.Network/routeTables@2023-04-01' = {
  location: resourceGroup().location
  name: resource_name
  properties: {
    disableBgpRoutePropagation: false
  }
}

resource virtualNetwork 'Microsoft.Network/virtualNetworks@2023-04-01' = {
  location: resourceGroup().location
  name: resource_name
  properties: {
    addressSpace: {
      addressPrefixes: [
        '10.0.0.0/16'
      ]
    }
    subnets: [
      {
        name: 'Default'
        properties: {
          addressPrefix: '10.0.0.0/24'
        }
      }
      {
        name: resource_name
        properties: {
          addressPrefix: '10.0.1.0/24'
          delegations: [
            {
              name: 'miDelegation'
              properties: {
                serviceName: 'Microsoft.Sql/managedInstances'
              }
            }
          ]
          networkSecurityGroup: {
            id: networkSecurityGroup.id
          }
          routeTable: {
            id: routeTable.id
          }
        }
      }
    ]
  }
}

