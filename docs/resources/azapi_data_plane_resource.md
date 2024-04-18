---
subcategory: ""
layout: "azapi"
page_title: "Azure Resource: azapi_data_plane_resource"
description: |-
  Manages a Azure data plane resource
---

# azapi_data_plane_resource

This resource can manage some Azure data plane resource.

## Example Usage

```hcl
terraform {
  required_providers {
    azapi = {
      source = "Azure/azapi"
    }
  }
}

provider "azapi" {
}

provider "azurerm" {
  features {}
}

data "azurerm_synapse_workspace" "example" {
  name                = "example-workspace"
  resource_group_name = azurerm_resource_group.example.name
}

resource "azapi_data_plane_resource" "dataset" {
  type      = "Microsoft.Synapse/workspaces/datasets@2020-12-01"
  parent_id = trimprefix(data.azurerm_synapse_workspace.example.connectivity_endpoints.dev, "https://")
  name      = "example-dataset"
  body = {
    properties = {
      type = "AzureBlob",
      typeProperties = {
        folderPath = {
          value = "@dataset().MyFolderPath"
          type  = "Expression"
        }
        fileName = {
          value = "@dataset().MyFileName"
          type  = "Expression"
        }
        format = {
          type = "TextFormat"
        }
      }
      parameters = {
        MyFolderPath = {
          type = "String"
        }
        MyFileName = {
          type = "String"
        }
      }
    }
  }
}

```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) Specifies the name of the azure resource. Changing this forces a new resource to be created. 

* `parent_id` - (Required) The ID of the azure resource in which this resource is created. Changing this forces a new resource to be created. 

* `type` - (Required) It is in a format like `<resource-type>@<api-version>`. `<api-version>` is version of the API used to manage this azure data plane resource.

-> **Note** For the available resource types and parent IDs, please refer to the `Available Resources` section below.

* `body` - (Required) A dynamic attribute that contains the request body used to create and update data plane resource. 

---

* `response_export_values` - (Optional) A list of path that needs to be exported from response body.
  Setting it to `["*"]` will export the full response body.
  Here's an example. If it sets to `["properties.loginServer", "properties.policies.quarantinePolicy.status"]`, it will set the following HCL object to computed property `output`.
```
{
  "properties" : {
    "loginServer" : "registry1.azurecr.io"
    "policies" : {
      "quarantinePolicy" = {
        "status" = "disabled"
      }
    }
  }
}
```

* `locks` - (Optional) A list of ARM resource IDs which are used to avoid create/modify/delete azapi resources at the same time.

* `ignore_missing_property` - (Optional) Whether ignore not returned properties like credentials in `body` to suppress plan-diff. Defaults to `true`. 
It's recommend to enable this option when some sensitive properties are not returned in response body, instead of setting them in `lifecycle.ignore_changes` because it will make the sensitive fields unable to update.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the azure resource.

* `output` - The output HCL object containing the properties specified in `response_export_values`. Here are some examples to use the values.
```
// it will output "registry1.azurecr.io"
output "login_server" {
  value = azapi_data_plane_resource.example.output.properties.loginServer
}

// it will output "disabled"
output "quarantine_policy" {
  value = azapi_data_plane_resource.example.output.properties.policies.quarantinePolicy.status
}
```

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the azure resource.
* `read` - (Defaults to 5 minutes) Used when retrieving the azure resource.
* `delete` - (Defaults to 30 minutes) Used when deleting the azure resource.

## Available Resources
| Resource Type | URL | Parent ID Example                                                                           |
| --- | --- |---------------------------------------------------------------------------------------------|
| Microsoft.AppConfiguration/configurationStores/keyValues | /kv/{key} | {storeName}.azconfig.io                                                                     |
| Microsoft.DeviceUpdate/accounts/groups | /deviceupdate/{instanceId}/management/groups/{groupId} | {accountName}.api.adu.microsoft.com/deviceupdate/{instanceName}                             |
| Microsoft.DeviceUpdate/accounts/groups/deployments | /deviceUpdate/{instanceId}/management/groups/{groupId}/deployments/{deploymentId} | {accountName}.api.adu.microsoft.com/deviceupdate/{instanceName}/management/groups/{groupId} |
| Microsoft.DeviceUpdate/accounts/v2/deployments | /deviceupdate/{instanceId}/v2/management/deployments/{deploymentId} | {accountName}.api.adu.microsoft.com/deviceupdate/{instanceName}                             |
| Microsoft.DeviceUpdate/accounts/v2/groups | /deviceupdate/{instanceId}/v2/management/groups/{groupId} | {accountName}.api.adu.microsoft.com/deviceupdate/{instanceName}                             |
| Microsoft.DigitalTwins/digitalTwinsInstances/digitaltwins | /digitaltwins/{id} | {instanceName}.api.weu.digitaltwins.azure.net                                               |
| Microsoft.DigitalTwins/digitalTwinsInstances/digitaltwins/relationships | /digitaltwins/{id}/relationships/{relationshipId} | {instanceName}.api.weu.digitaltwins.azure.net/digitaltwins/{digitalTwinId}                  |
| Microsoft.DigitalTwins/digitalTwinsInstances/eventroutes | /eventroutes/{id} | {instanceName}.api.weu.digitaltwins.azure.net                                               |
| Microsoft.DigitalTwins/digitalTwinsInstances/jobs/imports | /jobs/imports/{id} | {instanceName}.api.weu.digitaltwins.azure.net                                               |
| Microsoft.IoTCentral/IoTApps/organizations | /organizations/{organizationId} | {appSubdomain}.azureiotcentral.com                                                          |
| Microsoft.IoTCentral/IoTApps/scheduledJobs | /scheduledJobs/{scheduledJobId} | {appSubdomain}.azureiotcentral.com                                                          |
| Microsoft.IoTCentral/IoTApps/users | /users/{userId} | {appSubdomain}.azureiotcentral.com                                                          |
| Microsoft.IoTCentral/iotApps/apiTokens | /apiTokens/{tokenId} | {appSubdomain}.azureiotcentral.com                                                          |
| Microsoft.IoTCentral/iotApps/continuousDataExports | /continuousDataExports/{exportId} | {appSubdomain}.azureiotcentral.com                                                          |
| Microsoft.IoTCentral/iotApps/dashboards | /dashboards/{dashboardId} | {appSubdomain}.azureiotcentral.com                                                          |
| Microsoft.IoTCentral/iotApps/dataExport/destinations | /dataExport/destinations/{destinationId} | {appSubdomain}.azureiotcentral.com                                                          |
| Microsoft.IoTCentral/iotApps/dataExport/exports | /dataExport/exports/{exportId} | {appSubdomain}.azureiotcentral.com                                                          |
| Microsoft.IoTCentral/iotApps/deploymentManifests | /deploymentManifests/{deploymentManifestId} | {appSubdomain}.azureiotcentral.com                                                          |
| Microsoft.IoTCentral/iotApps/deviceGroups | /deviceGroups/{deviceGroupId} | {appSubdomain}.azureiotcentral.com                                                          |
| Microsoft.IoTCentral/iotApps/deviceTemplates | /deviceTemplates/{deviceTemplateId} | {appSubdomain}.azureiotcentral.com                                                          |
| Microsoft.IoTCentral/iotApps/devices | /devices/{deviceId} | {appSubdomain}.azureiotcentral.com                                                          |
| Microsoft.IoTCentral/iotApps/devices/attestation | /devices/{deviceId}/attestation | {appSubdomain}.azureiotcentral.com                                                          |
| Microsoft.IoTCentral/iotApps/devices/relationships | /devices/{deviceId}/relationships/{relationshipId} | {appSubdomain}.azureiotcentral.com/devices/{deviceId}                                       |
| Microsoft.IoTCentral/iotApps/enrollmentGroups | /enrollmentGroups/{enrollmentGroupId} | {appSubdomain}.azureiotcentral.com                                                          |
| Microsoft.IoTCentral/iotApps/enrollmentGroups/certificates | /enrollmentGroups/{enrollmentGroupId}/certificates/{entry} | {appSubdomain}.azureiotcentral.com/enrollmentGroups/{enrollmentGroupId}                     |
| Microsoft.KeyVault/vaults/certificates/contacts | /certificates/contacts | {vaultName}.vault.azure.net                                                                 |
| Microsoft.KeyVault/vaults/certificates/issuers | /certificates/issuers/{issuer-name} | {vaultName}.vault.azure.net                                                                 |
| Microsoft.KeyVault/vaults/storage | /storage/{storage-account-name} | {vaultName}.vault.azure.net                                                                 |
| Microsoft.KeyVault/vaults/storage/sas | /storage/{storage-account-name}/sas/{sas-definition-name} | {vaultName}.vault.azure.net/storage/{storage-account-name}                                  |
| Microsoft.Purview/accounts/Account/collections | /collections/{collectionName} | {accountName}.purview.azure.com                                                             |
| Microsoft.Purview/accounts/Account/resourceSetRuleConfigs | /resourceSetRuleConfigs/defaultResourceSetRuleConfig | {accountName}.purview.azure.com                                                             |
| Microsoft.Purview/accounts/Scanning/azureKeyVaults | /azureKeyVaults/{azureKeyVaultName} | {accountName}.purview.azure.com/scan                                                        |
| Microsoft.Purview/accounts/Scanning/classificationrules | /classificationrules/{classificationRuleName} | {accountName}.purview.azure.com/scan                                                        |
| Microsoft.Purview/accounts/Scanning/credentials | /credentials/{credentialName} | {accountName}.purview.azure.com/scan                                                        |
| Microsoft.Purview/accounts/Scanning/datasources | /datasources/{dataSourceName} | {accountName}.purview.azure.com/scan                                                        |
| Microsoft.Purview/accounts/Scanning/datasources/scans | /datasources/{dataSourceName}/scans/{scanName} | {accountName}.purview.azure.com/scan/datasources/{dataSourceName}                           |
| Microsoft.Purview/accounts/Scanning/datasources/scans/triggers | /datasources/{dataSourceName}/scans/{scanName}/triggers/default | {accountName}.purview.azure.com/scan/datasources/{dataSourceName}/scans/{scanName}          |
| Microsoft.Purview/accounts/Scanning/integrationruntimes | /integrationruntimes/{integrationRuntimeName} | {accountName}.purview.azure.com/scan                                                        |
| Microsoft.Purview/accounts/Scanning/managedvirtualnetworks/managedprivateendpoints | /managedvirtualnetworks/{managedVirtualNetworkName}/managedprivateendpoints/{managedPrivateEndpointName} | {accountName}.purview.azure.com/scan/managedvirtualnetworks/{managedVirtualNetworkName}     |
| Microsoft.Purview/accounts/Workflow/workflows | /workflows/{workflowId} | {accountName}.purview.azure.com                                                             |
| Microsoft.Synapse/workspaces/databases | /databases/{databaseName} | {workspaceName}.dev.azuresynapse.net                                                        |
| Microsoft.Synapse/workspaces/dataflows | /dataflows/{dataFlowName} | {workspaceName}.dev.azuresynapse.net                                                        |
| Microsoft.Synapse/workspaces/datasets | /datasets/{datasetName} | {workspaceName}.dev.azuresynapse.net                                                        |
| Microsoft.Synapse/workspaces/kqlScripts | /kqlScripts/{kqlScriptName} | {workspaceName}.dev.azuresynapse.net                                                        |
| Microsoft.Synapse/workspaces/libraries | /libraries/{libraryName} | {workspaceName}.dev.azuresynapse.net                                                        |
| Microsoft.Synapse/workspaces/linkconnections | /linkconnections/{linkConnectionName} | {workspaceName}.dev.azuresynapse.net                                                        |
| Microsoft.Synapse/workspaces/linkedservices | /linkedservices/{linkedServiceName} | {workspaceName}.dev.azuresynapse.net                                                        |
| Microsoft.Synapse/workspaces/managedVirtualNetworks/managedPrivateEndpoints | /managedVirtualNetworks/{managedVirtualNetworkName}/managedPrivateEndpoints/{managedPrivateEndpointName} | {workspaceName}.dev.azuresynapse.net/managedVirtualNetworks/{managedVirtualNetworkName}     |
| Microsoft.Synapse/workspaces/notebooks | /notebooks/{notebookName} | {workspaceName}.dev.azuresynapse.net                                                        |
| Microsoft.Synapse/workspaces/pipelines | /pipelines/{pipelineName} | {workspaceName}.dev.azuresynapse.net                                                        |
| Microsoft.Synapse/workspaces/roleAssignments | /roleAssignments/{roleAssignmentId} | {workspaceName}.dev.azuresynapse.net                                                        |
| Microsoft.Synapse/workspaces/sparkJobDefinitions | /sparkJobDefinitions/{sparkJobDefinitionName} | {workspaceName}.dev.azuresynapse.net                                                        |
| Microsoft.Synapse/workspaces/sparkconfigurations | /sparkconfigurations/{sparkConfigurationName} | {workspaceName}.dev.azuresynapse.net                                                        |
| Microsoft.Synapse/workspaces/sqlScripts | /sqlScripts/{sqlScriptName} | {workspaceName}.dev.azuresynapse.net                                                        |
| Microsoft.Synapse/workspaces/triggers | /triggers/{triggerName} | {workspaceName}.dev.azuresynapse.net                                                        |
