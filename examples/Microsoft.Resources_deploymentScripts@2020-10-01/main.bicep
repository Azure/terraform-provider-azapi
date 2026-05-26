param location string = 'westeurope'
param resource_name string = 'acctest0001'

resource deploymentScript 'Microsoft.Resources/deploymentScripts@2020-10-01' = {
  location: location
  name: resource_name
  kind: 'AzurePowerShell'
  properties: {
    azPowerShellVersion: '8.3'
    cleanupPreference: 'Always'
    environmentVariables: null
    retentionInterval: 'P1D'
    scriptContent: '\t\t$output = \'Hello\'\n\t\tWrite-Output $output\n\t\t$DeploymentScriptOutputs = @{}\n\t\t$DeploymentScriptOutputs[\'text\'] = $output\n'
    supportingScriptUris: null
    timeout: 'P1D'
  }
}

