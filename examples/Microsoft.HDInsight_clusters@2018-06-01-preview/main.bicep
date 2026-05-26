param location string = 'westeurope'
param resource_name string = 'acctest0001'
param rest_credential_password string = null
param vm_password string = null
param vm_username string = null

resource cluster 'Microsoft.HDInsight/clusters@2018-06-01-preview' = {
  location: location
  name: resource_name
  properties: {
    clusterDefinition: {
      componentVersion: {
        Spark: '2.4'
      }
      configurations: {
        gateway: {
          'restAuthCredential.isEnabled': true
          'restAuthCredential.password': rest_credential_password
          'restAuthCredential.username': 'acctestusrgw'
        }
      }
      kind: 'Spark'
    }
    clusterVersion: '4.0.3000.1'
    computeProfile: {
      roles: [
        {
          hardwareProfile: {
            vmSize: 'standard_a4_v2'
          }
          name: 'headnode'
          osProfile: {
            linuxOperatingSystemProfile: {
              password: vm_password
              username: vm_username
            }
          }
          targetInstanceCount: 2
        }
        {
          hardwareProfile: {
            vmSize: 'standard_a4_v2'
          }
          name: 'workernode'
          osProfile: {
            linuxOperatingSystemProfile: {
              password: vm_password
              username: vm_username
            }
          }
          targetInstanceCount: 3
        }
        {
          hardwareProfile: {
            vmSize: 'standard_a2_v2'
          }
          name: 'zookeepernode'
          osProfile: {
            linuxOperatingSystemProfile: {
              password: vm_password
              username: vm_username
            }
          }
          targetInstanceCount: 3
        }
      ]
    }
    encryptionInTransitProperties: {
      isEncryptionInTransitEnabled: false
    }
    minSupportedTlsVersion: '1.2'
    osType: 'Linux'
    storageProfile: {
      storageaccounts: [
        {
          container: container.name
          isDefault: true
          key: data.azapi_resource_action.listKeys.output.keys[0].value
          name: '${storageAccount.name}.blob.core.windows.net'
          resourceId: storageAccount.id
        }
      ]
    }
    tier: 'standard'
  }
}

resource container 'Microsoft.Storage/storageAccounts/blobServices/containers@2022-09-01' = {
  name: resource_name
  properties: {
    metadata: {
      key: 'value'
    }
  }
}

resource storageAccount 'Microsoft.Storage/storageAccounts@2021-09-01' = {
  location: location
  name: resource_name
  kind: 'StorageV2'
  properties: {
    accessTier: 'Hot'
    allowBlobPublicAccess: true
    allowCrossTenantReplication: true
    allowSharedKeyAccess: true
    defaultToOAuthAuthentication: false
    encryption: {
      keySource: 'Microsoft.Storage'
      services: {
        queue: {
          keyType: 'Service'
        }
        table: {
          keyType: 'Service'
        }
      }
    }
    isHnsEnabled: false
    isNfsV3Enabled: false
    isSftpEnabled: false
    minimumTlsVersion: 'TLS1_2'
    networkAcls: {
      defaultAction: 'Allow'
    }
    publicNetworkAccess: 'Enabled'
    supportsHttpsTrafficOnly: true
  }
  sku: {
    name: 'Standard_LRS'
  }
}

