param location string = 'westeurope'
param radius_server_secret string = null
param resource_name string = 'acctest0001'

resource vpnServerConfiguration 'Microsoft.Network/vpnServerConfigurations@2022-07-01' = {
  location: location
  name: resource_name
  properties: {
    radiusClientRootCertificates: []
    radiusServerAddress: ''
    radiusServerRootCertificates: []
    radiusServerSecret: radius_server_secret
    radiusServers: [
      {
        radiusServerAddress: '10.105.1.1'
        radiusServerScore: 15
        radiusServerSecret: radius_server_secret
      }
    ]
    vpnAuthenticationTypes: [
      'Radius'
    ]
    vpnClientIpsecPolicies: []
    vpnClientRevokedCertificates: []
    vpnClientRootCertificates: []
    vpnProtocols: [
      'OpenVPN'
      'IkeV2'
    ]
  }
}

