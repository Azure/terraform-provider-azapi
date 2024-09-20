package parse

import (
	"encoding/json"
	"strings"
)

type ApiPath struct {
	UrlFormat       string
	ResourceType    string
	URL             string
	ParentIDExample string
}

var apiPaths = make([]ApiPath, 0)

func init() {
	err := json.Unmarshal([]byte(raw), &apiPaths)
	if err != nil {
		panic(err)
	}
}

func findApiPathByResourceType(resourceType string) *ApiPath {
	for _, apiPath := range apiPaths {
		if strings.EqualFold(apiPath.ResourceType, resourceType) {
			return &apiPath
		}
	}
	return nil
}

const raw = `
[
  {
    "UrlFormat": "{parentId}/kv/{name}",
    "ResourceType": "Microsoft.AppConfiguration/configurationStores/keyValues",
    "ParentIDExample": "{storeName}.azconfig.io",
    "Url": "/kv/{key}"
  },
  {
    "UrlFormat": "{parentId}/management/groups/{name}",
    "ResourceType": "Microsoft.DeviceUpdate/accounts/groups",
    "ParentIDExample": "{accountName}.api.adu.microsoft.com/deviceupdate/{instanceName}",
    "Url": "/deviceupdate/{instanceId}/management/groups/{groupId}"
  },
  {
    "UrlFormat": "{parentId}/deployments/{name}",
    "ResourceType": "Microsoft.DeviceUpdate/accounts/groups/deployments",
    "ParentIDExample": "{accountName}.api.adu.microsoft.com/deviceupdate/{instanceName}/management/groups/{groupId}",
    "Url": "/deviceUpdate/{instanceId}/management/groups/{groupId}/deployments/{deploymentId}"
  },
  {
    "UrlFormat": "{parentId}/v2/management/deployments/{name}",
    "ResourceType": "Microsoft.DeviceUpdate/accounts/v2/deployments",
    "ParentIDExample": "{accountName}.api.adu.microsoft.com/deviceupdate/{instanceName}",
    "Url": "/deviceupdate/{instanceId}/v2/management/deployments/{deploymentId}"
  },
  {
    "UrlFormat": "{parentId}/v2/management/groups/{name}",
    "ResourceType": "Microsoft.DeviceUpdate/accounts/v2/groups",
    "ParentIDExample": "{accountName}.api.adu.microsoft.com/deviceupdate/{instanceName}",
    "Url": "/deviceupdate/{instanceId}/v2/management/groups/{groupId}"
  },
  {
    "UrlFormat": "{parentId}/digitaltwins/{name}",
    "ResourceType": "Microsoft.DigitalTwins/digitalTwinsInstances/digitaltwins",
    "ParentIDExample": "{instanceName}.api.weu.digitaltwins.azure.net",
    "Url": "/digitaltwins/{id}"
  },
  {
    "UrlFormat": "{parentId}/relationships/{name}",
    "ResourceType": "Microsoft.DigitalTwins/digitalTwinsInstances/digitaltwins/relationships",
    "ParentIDExample": "{instanceName}.api.weu.digitaltwins.azure.net/digitaltwins/{digitalTwinId}",
    "Url": "/digitaltwins/{id}/relationships/{relationshipId}"
  },
  {
    "UrlFormat": "{parentId}/eventroutes/{name}",
    "ResourceType": "Microsoft.DigitalTwins/digitalTwinsInstances/eventroutes",
    "ParentIDExample": "{instanceName}.api.weu.digitaltwins.azure.net",
    "Url": "/eventroutes/{id}"
  },
  {
    "UrlFormat": "{parentId}/jobs/imports/{name}",
    "ResourceType": "Microsoft.DigitalTwins/digitalTwinsInstances/jobs/imports",
    "ParentIDExample": "{instanceName}.api.weu.digitaltwins.azure.net",
    "Url": "/jobs/imports/{id}"
  },
  {
    "UrlFormat": "{parentId}/api/organizations/{name}",
    "ResourceType": "Microsoft.IoTCentral/IoTApps/organizations",
    "ParentIDExample": "{appSubdomain}.azureiotcentral.com",
    "Url": "/organizations/{organizationId}"
  },
  {
    "UrlFormat": "{parentId}/api/scheduledJobs/{name}",
    "ResourceType": "Microsoft.IoTCentral/IoTApps/scheduledJobs",
    "ParentIDExample": "{appSubdomain}.azureiotcentral.com",
    "Url": "/scheduledJobs/{scheduledJobId}"
  },
  {
    "UrlFormat": "{parentId}/api/users/{name}",
    "ResourceType": "Microsoft.IoTCentral/IoTApps/users",
    "ParentIDExample": "{appSubdomain}.azureiotcentral.com",
    "Url": "/users/{userId}"
  },
  {
    "UrlFormat": "{parentId}/api/apiTokens/{tokenId}",
    "ResourceType": "Microsoft.IoTCentral/iotApps/apiTokens",
    "ParentIDExample": "{appSubdomain}.azureiotcentral.com",
    "Url": "/apiTokens/{tokenId}"
  },
  {
    "UrlFormat": "{parentId}/api/continuousDataExports/{name}",
    "ResourceType": "Microsoft.IoTCentral/iotApps/continuousDataExports",
    "ParentIDExample": "{appSubdomain}.azureiotcentral.com",
    "Url": "/continuousDataExports/{exportId}"
  },
  {
    "UrlFormat": "{parentId}/api/dashboards/{name}",
    "ResourceType": "Microsoft.IoTCentral/iotApps/dashboards",
    "ParentIDExample": "{appSubdomain}.azureiotcentral.com",
    "Url": "/dashboards/{dashboardId}"
  },
  {
    "UrlFormat": "{parentId}/api/dataExport/destinations/{name}",
    "ResourceType": "Microsoft.IoTCentral/iotApps/dataExport/destinations",
    "ParentIDExample": "{appSubdomain}.azureiotcentral.com",
    "Url": "/dataExport/destinations/{destinationId}"
  },
  {
    "UrlFormat": "{parentId}/api/dataExport/exports/{name}",
    "ResourceType": "Microsoft.IoTCentral/iotApps/dataExport/exports",
    "ParentIDExample": "{appSubdomain}.azureiotcentral.com",
    "Url": "/dataExport/exports/{exportId}"
  },
  {
    "UrlFormat": "{parentId}/api/deploymentManifests/{name}",
    "ResourceType": "Microsoft.IoTCentral/iotApps/deploymentManifests",
    "ParentIDExample": "{appSubdomain}.azureiotcentral.com",
    "Url": "/deploymentManifests/{deploymentManifestId}"
  },
  {
    "UrlFormat": "{parentId}/api/deviceGroups/{name}",
    "ResourceType": "Microsoft.IoTCentral/iotApps/deviceGroups",
    "ParentIDExample": "{appSubdomain}.azureiotcentral.com",
    "Url": "/deviceGroups/{deviceGroupId}"
  },
  {
    "UrlFormat": "{parentId}/api/deviceTemplates/{name}",
    "ResourceType": "Microsoft.IoTCentral/iotApps/deviceTemplates",
    "ParentIDExample": "{appSubdomain}.azureiotcentral.com",
    "Url": "/deviceTemplates/{deviceTemplateId}"
  },
  {
    "UrlFormat": "{parentId}/api/devices/{name}",
    "ResourceType": "Microsoft.IoTCentral/iotApps/devices",
    "ParentIDExample": "{appSubdomain}.azureiotcentral.com",
    "Url": "/devices/{deviceId}"
  },
  {
    "UrlFormat": "{parentId}/api/devices/{name}/attestation",
    "ResourceType": "Microsoft.IoTCentral/iotApps/devices/attestation",
    "ParentIDExample": "{appSubdomain}.azureiotcentral.com",
    "Url": "/devices/{deviceId}/attestation"
  },
  {
    "UrlFormat": "{parentId}/relationships/{name}",
    "ResourceType": "Microsoft.IoTCentral/iotApps/devices/relationships",
    "ParentIDExample": "{appSubdomain}.azureiotcentral.com/devices/{deviceId}",
    "Url": "/devices/{deviceId}/relationships/{relationshipId}"
  },
  {
    "UrlFormat": "{parentId}/api/enrollmentGroups/{name}",
    "ResourceType": "Microsoft.IoTCentral/iotApps/enrollmentGroups",
    "ParentIDExample": "{appSubdomain}.azureiotcentral.com",
    "Url": "/enrollmentGroups/{enrollmentGroupId}"
  },
  {
    "UrlFormat": "{parentId}/certificates/{name}",
    "ResourceType": "Microsoft.IoTCentral/iotApps/enrollmentGroups/certificates",
    "ParentIDExample": "{appSubdomain}.azureiotcentral.com/enrollmentGroups/{enrollmentGroupId}",
    "Url": "/enrollmentGroups/{enrollmentGroupId}/certificates/{entry}"
  },
  {
    "UrlFormat": "{parentId}/certificates/contacts",
    "ResourceType": "Microsoft.KeyVault/vaults/certificates/contacts",
    "ParentIDExample": "{vaultName}.vault.azure.net",
    "Url": "/certificates/contacts"
  },
  {
    "UrlFormat": "{parentId}/certificates/issuers/{name}",
    "ResourceType": "Microsoft.KeyVault/vaults/certificates/issuers",
    "ParentIDExample": "{vaultName}.vault.azure.net",
    "Url": "/certificates/issuers/{issuer-name}"
  },
  {
    "UrlFormat": "{parentId}/storage/{name}",
    "ResourceType": "Microsoft.KeyVault/vaults/storage",
    "ParentIDExample": "{vaultName}.vault.azure.net",
    "Url": "/storage/{storage-account-name}"
  },
  {
    "UrlFormat": "{parentId}/sas/{name}",
    "ResourceType": "Microsoft.KeyVault/vaults/storage/sas",
    "ParentIDExample": "{vaultName}.vault.azure.net/storage/{storage-account-name}",
    "Url": "/storage/{storage-account-name}/sas/{sas-definition-name}"
  },
  {
    "UrlFormat": "{parentId}/collections/{name}",
    "ResourceType": "Microsoft.Purview/accounts/Account/collections",
    "ParentIDExample": "{accountName}.purview.azure.com",
    "Url": "/collections/{collectionName}"
  },
  {
    "UrlFormat": "{parentId}/resourceSetRuleConfigs/{name=defaultResourceSetRuleConfig}",
    "ResourceType": "Microsoft.Purview/accounts/Account/resourceSetRuleConfigs",
    "ParentIDExample": "{accountName}.purview.azure.com",
    "Url": "/resourceSetRuleConfigs/defaultResourceSetRuleConfig"
  },
  {
    "UrlFormat": "{parentId}/azureKeyVaults/{name}",
    "ResourceType": "Microsoft.Purview/accounts/Scanning/azureKeyVaults",
    "ParentIDExample": "{accountName}.purview.azure.com/scan",
    "Url": "/azureKeyVaults/{azureKeyVaultName}"
  },
  {
    "UrlFormat": "{parentId}/classificationrules/{name}",
    "ResourceType": "Microsoft.Purview/accounts/Scanning/classificationrules",
    "ParentIDExample": "{accountName}.purview.azure.com/scan",
    "Url": "/classificationrules/{classificationRuleName}"
  },
  {
    "UrlFormat": "{parentId}/credentials/{name}",
    "ResourceType": "Microsoft.Purview/accounts/Scanning/credentials",
    "ParentIDExample": "{accountName}.purview.azure.com/scan",
    "Url": "/credentials/{credentialName}"
  },
  {
    "UrlFormat": "{parentId}/datasources/{name}",
    "ResourceType": "Microsoft.Purview/accounts/Scanning/datasources",
    "ParentIDExample": "{accountName}.purview.azure.com/scan",
    "Url": "/datasources/{dataSourceName}"
  },
  {
    "UrlFormat": "{parentId}/scans/{name}",
    "ResourceType": "Microsoft.Purview/accounts/Scanning/datasources/scans",
    "ParentIDExample": "{accountName}.purview.azure.com/scan/datasources/{dataSourceName}",
    "Url": "/datasources/{dataSourceName}/scans/{scanName}"
  },
  {
    "UrlFormat": "{parentId}/triggers/{name=default}",
    "ResourceType": "Microsoft.Purview/accounts/Scanning/datasources/scans/triggers",
    "ParentIDExample": "{accountName}.purview.azure.com/scan/datasources/{dataSourceName}/scans/{scanName}",
    "Url": "/datasources/{dataSourceName}/scans/{scanName}/triggers/default"
  },
  {
    "UrlFormat": "{parentId}/integrationruntimes/{name}",
    "ResourceType": "Microsoft.Purview/accounts/Scanning/integrationruntimes",
    "ParentIDExample": "{accountName}.purview.azure.com/scan",
    "Url": "/integrationruntimes/{integrationRuntimeName}"
  },
  {
    "UrlFormat": "{parentId}/managedvirtualnetworks/{name}",
    "ResourceType": "Microsoft.Purview/accounts/Scanning/managedvirtualnetworks",
    "ParentIDExample": "{accountName}.purview.azure.com/scan",
    "Url": "/managedvirtualnetworks/{managedVirtualNetworkName}"
  },
  {
    "UrlFormat": "{parentId}/managedprivateendpoints/{name}",
    "ResourceType": "Microsoft.Purview/accounts/Scanning/managedvirtualnetworks/managedprivateendpoints",
    "ParentIDExample": "{accountName}.purview.azure.com/scan/managedvirtualnetworks/{managedVirtualNetworkName}",
    "Url": "/managedvirtualnetworks/{managedVirtualNetworkName}/managedprivateendpoints/{managedPrivateEndpointName}"
  },
  {
    "UrlFormat": "{parentId}/workflows/{name}",
    "ResourceType": "Microsoft.Purview/accounts/Workflow/workflows",
    "ParentIDExample": "{accountName}.purview.azure.com",
    "Url": "/workflows/{workflowId}"
  },
  {
    "UrlFormat": "{parentId}/databases/{name}",
    "ResourceType": "Microsoft.Synapse/workspaces/databases",
    "ParentIDExample": "{workspaceName}.dev.azuresynapse.net",
    "Url": "/databases/{databaseName}"
  },
  {
    "UrlFormat": "{parentId}/dataflows/{name}",
    "ResourceType": "Microsoft.Synapse/workspaces/dataflows",
    "ParentIDExample": "{workspaceName}.dev.azuresynapse.net",
    "Url": "/dataflows/{dataFlowName}"
  },
  {
    "UrlFormat": "{parentId}/datasets/{name}",
    "ResourceType": "Microsoft.Synapse/workspaces/datasets",
    "ParentIDExample": "{workspaceName}.dev.azuresynapse.net",
    "Url": "/datasets/{datasetName}"
  },
  {
    "UrlFormat": "{parentId}/kqlScripts/{name}",
    "ResourceType": "Microsoft.Synapse/workspaces/kqlScripts",
    "ParentIDExample": "{workspaceName}.dev.azuresynapse.net",
    "Url": "/kqlScripts/{kqlScriptName}"
  },
  {
    "UrlFormat": "{parentId}/libraries/{name}",
    "ResourceType": "Microsoft.Synapse/workspaces/libraries",
    "ParentIDExample": "{workspaceName}.dev.azuresynapse.net",
    "Url": "/libraries/{libraryName}"
  },
  {
    "UrlFormat": "{parentId}/linkconnections/{name}",
    "ResourceType": "Microsoft.Synapse/workspaces/linkconnections",
    "ParentIDExample": "{workspaceName}.dev.azuresynapse.net",
    "Url": "/linkconnections/{linkConnectionName}"
  },
  {
    "UrlFormat": "{parentId}/linkedservices/{name}",
    "ResourceType": "Microsoft.Synapse/workspaces/linkedservices",
    "ParentIDExample": "{workspaceName}.dev.azuresynapse.net",
    "Url": "/linkedservices/{linkedServiceName}"
  },
  {
    "UrlFormat": "{parentId}/managedPrivateEndpoints/{name}",
    "ResourceType": "Microsoft.Synapse/workspaces/managedVirtualNetworks/managedPrivateEndpoints",
    "ParentIDExample": "{workspaceName}.dev.azuresynapse.net/managedVirtualNetworks/{managedVirtualNetworkName}",
    "Url": "/managedVirtualNetworks/{managedVirtualNetworkName}/managedPrivateEndpoints/{managedPrivateEndpointName}"
  },
  {
    "UrlFormat": "{parentId}/notebooks/{name}",
    "ResourceType": "Microsoft.Synapse/workspaces/notebooks",
    "ParentIDExample": "{workspaceName}.dev.azuresynapse.net",
    "Url": "/notebooks/{notebookName}"
  },
  {
    "UrlFormat": "{parentId}/pipelines/{name}",
    "ResourceType": "Microsoft.Synapse/workspaces/pipelines",
    "ParentIDExample": "{workspaceName}.dev.azuresynapse.net",
    "Url": "/pipelines/{pipelineName}"
  },
  {
    "UrlFormat": "{parentId}/roleAssignments/{name}",
    "ResourceType": "Microsoft.Synapse/workspaces/roleAssignments",
    "ParentIDExample": "{workspaceName}.dev.azuresynapse.net",
    "Url": "/roleAssignments/{roleAssignmentId}"
  },
  {
    "UrlFormat": "{parentId}/sparkJobDefinitions/{name}",
    "ResourceType": "Microsoft.Synapse/workspaces/sparkJobDefinitions",
    "ParentIDExample": "{workspaceName}.dev.azuresynapse.net",
    "Url": "/sparkJobDefinitions/{sparkJobDefinitionName}"
  },
  {
    "UrlFormat": "{parentId}/sparkconfigurations/{name}",
    "ResourceType": "Microsoft.Synapse/workspaces/sparkconfigurations",
    "ParentIDExample": "{workspaceName}.dev.azuresynapse.net",
    "Url": "/sparkconfigurations/{sparkConfigurationName}"
  },
  {
    "UrlFormat": "{parentId}/sqlScripts/{name}",
    "ResourceType": "Microsoft.Synapse/workspaces/sqlScripts",
    "ParentIDExample": "{workspaceName}.dev.azuresynapse.net",
    "Url": "/sqlScripts/{sqlScriptName}"
  },
  {
    "UrlFormat": "{parentId}/triggers/{name}",
    "ResourceType": "Microsoft.Synapse/workspaces/triggers",
    "ParentIDExample": "{workspaceName}.dev.azuresynapse.net",
    "Url": "/triggers/{triggerName}"
  }
]`
