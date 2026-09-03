---
subcategory: "Microsoft.RedHatOpenShift - Azure Red Hat OpenShift"
page_title: "openShiftClusters"
description: |-
  Manages a fully managed Azure Red Hat OpenShift Cluster (also known as ARO).
---

# Microsoft.RedHatOpenShift/openShiftClusters - fully managed Azure Red Hat OpenShift Cluster (also known as ARO)

This article demonstrates how to use `azapi` provider to manage the fully managed Azure Red Hat OpenShift Cluster (also known as ARO) resource in Azure.



## Example Usage

### default

```hcl
terraform {
  required_providers {
    azapi = {
      source = "Azure/azapi"
    }
    azuread = {
      source = "hashicorp/azuread"
    }
    time = {
      source = "hashicorp/time"
    }
  }
}

provider "azapi" {
  skip_provider_registration = false
}

provider "azuread" {}

data "azapi_client_config" "current" {}

variable "resource_name" {
  type    = string
  default = "acctest0001"
}

variable "location" {
  type    = string
  default = "australiaeast"
}

// Look up the Azure Red Hat OpenShift RP first-party service principal
data "azuread_service_principal" "aroRP" {
  display_name = "Azure Red Hat OpenShift RP"
}

// Look up built-in role definitions by name
data "azapi_resource_list" "roleDefinitions" {
  type      = "Microsoft.Authorization/roleDefinitions@2022-04-01"
  parent_id = "/subscriptions/${data.azapi_client_config.current.subscription_id}"
  response_export_values = {
    managed_identity_operator    = "value[?properties.roleName == 'Managed Identity Operator'].id | [0]"
    aro_cloud_controller_manager = "value[?properties.roleName == 'Azure Red Hat OpenShift Cloud Controller Manager'].id | [0]"
    aro_cluster_ingress_operator = "value[?properties.roleName == 'Azure Red Hat OpenShift Cluster Ingress Operator'].id | [0]"
    aro_machine_api_operator     = "value[?properties.roleName == 'Azure Red Hat OpenShift Machine API Operator'].id | [0]"
    aro_network_operator         = "value[?properties.roleName == 'Azure Red Hat OpenShift Network Operator'].id | [0]"
    aro_file_storage_operator    = "value[?properties.roleName == 'Azure Red Hat OpenShift File Storage Operator'].id | [0]"
    aro_image_registry_operator  = "value[?properties.roleName == 'Azure Red Hat OpenShift Image Registry Operator'].id | [0]"
    aro_service_operator         = "value[?properties.roleName == 'Azure Red Hat OpenShift Service Operator'].id | [0]"
    aro_federated_credential     = "value[?properties.roleName == 'Azure Red Hat OpenShift Federated Credential'].id | [0]"
    aro_first_party_network      = "value[?properties.roleName == 'Azure Red Hat OpenShift First Party Network'].id | [0]"
  }
}

resource "azapi_resource" "resourceGroup" {
  type     = "Microsoft.Resources/resourceGroups@2020-06-01"
  name     = var.resource_name
  location = var.location
}

// Cluster identity
resource "azapi_resource" "clusterIdentity" {
  type                      = "Microsoft.ManagedIdentity/userAssignedIdentities@2023-01-31"
  parent_id                 = azapi_resource.resourceGroup.id
  name                      = "${var.resource_name}-aro-cluster"
  location                  = var.location
  schema_validation_enabled = false
  response_export_values    = ["*"]
}

// Platform workload identities
resource "azapi_resource" "cloudControllerManagerIdentity" {
  type                      = "Microsoft.ManagedIdentity/userAssignedIdentities@2023-01-31"
  parent_id                 = azapi_resource.resourceGroup.id
  name                      = "${var.resource_name}-cloud-controller-manager"
  location                  = var.location
  schema_validation_enabled = false
  response_export_values    = ["*"]
}

resource "azapi_resource" "ingressIdentity" {
  type                      = "Microsoft.ManagedIdentity/userAssignedIdentities@2023-01-31"
  parent_id                 = azapi_resource.resourceGroup.id
  name                      = "${var.resource_name}-ingress"
  location                  = var.location
  schema_validation_enabled = false
  response_export_values    = ["*"]
}

resource "azapi_resource" "machineApiIdentity" {
  type                      = "Microsoft.ManagedIdentity/userAssignedIdentities@2023-01-31"
  parent_id                 = azapi_resource.resourceGroup.id
  name                      = "${var.resource_name}-machine-api"
  location                  = var.location
  schema_validation_enabled = false
  response_export_values    = ["*"]
}

resource "azapi_resource" "diskCsiDriverIdentity" {
  type                      = "Microsoft.ManagedIdentity/userAssignedIdentities@2023-01-31"
  parent_id                 = azapi_resource.resourceGroup.id
  name                      = "${var.resource_name}-disk-csi-driver"
  location                  = var.location
  schema_validation_enabled = false
  response_export_values    = ["*"]
}

resource "azapi_resource" "cloudNetworkConfigIdentity" {
  type                      = "Microsoft.ManagedIdentity/userAssignedIdentities@2023-01-31"
  parent_id                 = azapi_resource.resourceGroup.id
  name                      = "${var.resource_name}-cloud-network-config"
  location                  = var.location
  schema_validation_enabled = false
  response_export_values    = ["*"]
}

resource "azapi_resource" "imageRegistryIdentity" {
  type                      = "Microsoft.ManagedIdentity/userAssignedIdentities@2023-01-31"
  parent_id                 = azapi_resource.resourceGroup.id
  name                      = "${var.resource_name}-image-registry"
  location                  = var.location
  schema_validation_enabled = false
  response_export_values    = ["*"]
}

resource "azapi_resource" "fileCsiDriverIdentity" {
  type                      = "Microsoft.ManagedIdentity/userAssignedIdentities@2023-01-31"
  parent_id                 = azapi_resource.resourceGroup.id
  name                      = "${var.resource_name}-file-csi-driver"
  location                  = var.location
  schema_validation_enabled = false
  response_export_values    = ["*"]
}

resource "azapi_resource" "aroOperatorIdentity" {
  type                      = "Microsoft.ManagedIdentity/userAssignedIdentities@2023-01-31"
  parent_id                 = azapi_resource.resourceGroup.id
  name                      = "${var.resource_name}-aro-operator"
  location                  = var.location
  schema_validation_enabled = false
  response_export_values    = ["*"]
}

// Networking
resource "azapi_resource" "virtualNetwork" {
  type      = "Microsoft.Network/virtualNetworks@2022-07-01"
  parent_id = azapi_resource.resourceGroup.id
  name      = var.resource_name
  location  = var.location
  body = {
    properties = {
      addressSpace = {
        addressPrefixes = ["10.0.0.0/22"]
      }
    }
  }
  schema_validation_enabled = false
  response_export_values    = ["*"]
  lifecycle {
    ignore_changes = [body.properties.subnets]
  }
}

resource "azapi_resource" "masterSubnet" {
  type      = "Microsoft.Network/virtualNetworks/subnets@2022-07-01"
  parent_id = azapi_resource.virtualNetwork.id
  name      = "master"
  body = {
    properties = {
      addressPrefix = "10.0.0.0/23"
      serviceEndpoints = [
        {
          service = "Microsoft.Storage"
        },
        {
          service = "Microsoft.ContainerRegistry"
        },
      ]
    }
  }
  schema_validation_enabled = false
  response_export_values    = ["*"]
}

resource "azapi_resource" "workerSubnet" {
  type      = "Microsoft.Network/virtualNetworks/subnets@2022-07-01"
  parent_id = azapi_resource.virtualNetwork.id
  name      = "worker"
  body = {
    properties = {
      addressPrefix = "10.0.2.0/23"
      serviceEndpoints = [
        {
          service = "Microsoft.Storage"
        },
        {
          service = "Microsoft.ContainerRegistry"
        },
      ]
    }
  }
  schema_validation_enabled = false
  response_export_values    = ["*"]
  depends_on                = [azapi_resource.masterSubnet]
}

// Role assignments: cluster identity -> Federated Credential over each platform workload identity
resource "azapi_resource" "roleAssignment_clusterIdentity_ccm" {
  type      = "Microsoft.Authorization/roleAssignments@2022-04-01"
  parent_id = azapi_resource.cloudControllerManagerIdentity.id
  name      = uuid()
  body = {
    properties = {
      principalId      = azapi_resource.clusterIdentity.output.properties.principalId
      principalType    = "ServicePrincipal"
      roleDefinitionId = data.azapi_resource_list.roleDefinitions.output.aro_federated_credential
    }
  }
  schema_validation_enabled = false
  lifecycle {
    ignore_changes = [name]
  }
}

resource "azapi_resource" "roleAssignment_clusterIdentity_ingress" {
  type      = "Microsoft.Authorization/roleAssignments@2022-04-01"
  parent_id = azapi_resource.ingressIdentity.id
  name      = uuid()
  body = {
    properties = {
      principalId      = azapi_resource.clusterIdentity.output.properties.principalId
      principalType    = "ServicePrincipal"
      roleDefinitionId = data.azapi_resource_list.roleDefinitions.output.aro_federated_credential
    }
  }
  schema_validation_enabled = false
  lifecycle {
    ignore_changes = [name]
  }
}

resource "azapi_resource" "roleAssignment_clusterIdentity_machineApi" {
  type      = "Microsoft.Authorization/roleAssignments@2022-04-01"
  parent_id = azapi_resource.machineApiIdentity.id
  name      = uuid()
  body = {
    properties = {
      principalId      = azapi_resource.clusterIdentity.output.properties.principalId
      principalType    = "ServicePrincipal"
      roleDefinitionId = data.azapi_resource_list.roleDefinitions.output.aro_federated_credential
    }
  }
  schema_validation_enabled = false
  lifecycle {
    ignore_changes = [name]
  }
}

resource "azapi_resource" "roleAssignment_clusterIdentity_diskCsiDriver" {
  type      = "Microsoft.Authorization/roleAssignments@2022-04-01"
  parent_id = azapi_resource.diskCsiDriverIdentity.id
  name      = uuid()
  body = {
    properties = {
      principalId      = azapi_resource.clusterIdentity.output.properties.principalId
      principalType    = "ServicePrincipal"
      roleDefinitionId = data.azapi_resource_list.roleDefinitions.output.aro_federated_credential
    }
  }
  schema_validation_enabled = false
  lifecycle {
    ignore_changes = [name]
  }
}

resource "azapi_resource" "roleAssignment_clusterIdentity_cloudNetworkConfig" {
  type      = "Microsoft.Authorization/roleAssignments@2022-04-01"
  parent_id = azapi_resource.cloudNetworkConfigIdentity.id
  name      = uuid()
  body = {
    properties = {
      principalId      = azapi_resource.clusterIdentity.output.properties.principalId
      principalType    = "ServicePrincipal"
      roleDefinitionId = data.azapi_resource_list.roleDefinitions.output.aro_federated_credential
    }
  }
  schema_validation_enabled = false
  lifecycle {
    ignore_changes = [name]
  }
}

resource "azapi_resource" "roleAssignment_clusterIdentity_imageRegistry" {
  type      = "Microsoft.Authorization/roleAssignments@2022-04-01"
  parent_id = azapi_resource.imageRegistryIdentity.id
  name      = uuid()
  body = {
    properties = {
      principalId      = azapi_resource.clusterIdentity.output.properties.principalId
      principalType    = "ServicePrincipal"
      roleDefinitionId = data.azapi_resource_list.roleDefinitions.output.aro_federated_credential
    }
  }
  schema_validation_enabled = false
  lifecycle {
    ignore_changes = [name]
  }
}

resource "azapi_resource" "roleAssignment_clusterIdentity_fileCsiDriver" {
  type      = "Microsoft.Authorization/roleAssignments@2022-04-01"
  parent_id = azapi_resource.fileCsiDriverIdentity.id
  name      = uuid()
  body = {
    properties = {
      principalId      = azapi_resource.clusterIdentity.output.properties.principalId
      principalType    = "ServicePrincipal"
      roleDefinitionId = data.azapi_resource_list.roleDefinitions.output.aro_federated_credential
    }
  }
  schema_validation_enabled = false
  lifecycle {
    ignore_changes = [name]
  }
}

resource "azapi_resource" "roleAssignment_clusterIdentity_aroOperator" {
  type      = "Microsoft.Authorization/roleAssignments@2022-04-01"
  parent_id = azapi_resource.aroOperatorIdentity.id
  name      = uuid()
  body = {
    properties = {
      principalId      = azapi_resource.clusterIdentity.output.properties.principalId
      principalType    = "ServicePrincipal"
      roleDefinitionId = data.azapi_resource_list.roleDefinitions.output.aro_federated_credential
    }
  }
  schema_validation_enabled = false
  lifecycle {
    ignore_changes = [name]
  }
}

// Role assignments: cloud-controller-manager on subnets
resource "azapi_resource" "roleAssignment_ccm_masterSubnet" {
  type      = "Microsoft.Authorization/roleAssignments@2022-04-01"
  parent_id = azapi_resource.masterSubnet.id
  name      = uuid()
  body = {
    properties = {
      principalId      = azapi_resource.cloudControllerManagerIdentity.output.properties.principalId
      principalType    = "ServicePrincipal"
      roleDefinitionId = data.azapi_resource_list.roleDefinitions.output.aro_cloud_controller_manager
    }
  }
  schema_validation_enabled = false
  lifecycle {
    ignore_changes = [name]
  }
}

resource "azapi_resource" "roleAssignment_ccm_workerSubnet" {
  type      = "Microsoft.Authorization/roleAssignments@2022-04-01"
  parent_id = azapi_resource.workerSubnet.id
  name      = uuid()
  body = {
    properties = {
      principalId      = azapi_resource.cloudControllerManagerIdentity.output.properties.principalId
      principalType    = "ServicePrincipal"
      roleDefinitionId = data.azapi_resource_list.roleDefinitions.output.aro_cloud_controller_manager
    }
  }
  schema_validation_enabled = false
  lifecycle {
    ignore_changes = [name]
  }
}

// Role assignments: ingress on subnets
resource "azapi_resource" "roleAssignment_ingress_masterSubnet" {
  type      = "Microsoft.Authorization/roleAssignments@2022-04-01"
  parent_id = azapi_resource.masterSubnet.id
  name      = uuid()
  body = {
    properties = {
      principalId      = azapi_resource.ingressIdentity.output.properties.principalId
      principalType    = "ServicePrincipal"
      roleDefinitionId = data.azapi_resource_list.roleDefinitions.output.aro_cluster_ingress_operator
    }
  }
  schema_validation_enabled = false
  lifecycle {
    ignore_changes = [name]
  }
}

resource "azapi_resource" "roleAssignment_ingress_workerSubnet" {
  type      = "Microsoft.Authorization/roleAssignments@2022-04-01"
  parent_id = azapi_resource.workerSubnet.id
  name      = uuid()
  body = {
    properties = {
      principalId      = azapi_resource.ingressIdentity.output.properties.principalId
      principalType    = "ServicePrincipal"
      roleDefinitionId = data.azapi_resource_list.roleDefinitions.output.aro_cluster_ingress_operator
    }
  }
  schema_validation_enabled = false
  lifecycle {
    ignore_changes = [name]
  }
}

// Role assignments: machine-api on subnets
resource "azapi_resource" "roleAssignment_machineApi_masterSubnet" {
  type      = "Microsoft.Authorization/roleAssignments@2022-04-01"
  parent_id = azapi_resource.masterSubnet.id
  name      = uuid()
  body = {
    properties = {
      principalId      = azapi_resource.machineApiIdentity.output.properties.principalId
      principalType    = "ServicePrincipal"
      roleDefinitionId = data.azapi_resource_list.roleDefinitions.output.aro_machine_api_operator
    }
  }
  schema_validation_enabled = false
  lifecycle {
    ignore_changes = [name]
  }
}

resource "azapi_resource" "roleAssignment_machineApi_workerSubnet" {
  type      = "Microsoft.Authorization/roleAssignments@2022-04-01"
  parent_id = azapi_resource.workerSubnet.id
  name      = uuid()
  body = {
    properties = {
      principalId      = azapi_resource.machineApiIdentity.output.properties.principalId
      principalType    = "ServicePrincipal"
      roleDefinitionId = data.azapi_resource_list.roleDefinitions.output.aro_machine_api_operator
    }
  }
  schema_validation_enabled = false
  lifecycle {
    ignore_changes = [name]
  }
}

// Role assignments: cloud-network-config on vnet
resource "azapi_resource" "roleAssignment_cloudNetworkConfig_vnet" {
  type      = "Microsoft.Authorization/roleAssignments@2022-04-01"
  parent_id = azapi_resource.virtualNetwork.id
  name      = uuid()
  body = {
    properties = {
      principalId      = azapi_resource.cloudNetworkConfigIdentity.output.properties.principalId
      principalType    = "ServicePrincipal"
      roleDefinitionId = data.azapi_resource_list.roleDefinitions.output.aro_network_operator
    }
  }
  schema_validation_enabled = false
  lifecycle {
    ignore_changes = [name]
  }
}

// Role assignments: file-csi-driver on vnet
resource "azapi_resource" "roleAssignment_fileCsiDriver_vnet" {
  type      = "Microsoft.Authorization/roleAssignments@2022-04-01"
  parent_id = azapi_resource.virtualNetwork.id
  name      = uuid()
  body = {
    properties = {
      principalId      = azapi_resource.fileCsiDriverIdentity.output.properties.principalId
      principalType    = "ServicePrincipal"
      roleDefinitionId = data.azapi_resource_list.roleDefinitions.output.aro_file_storage_operator
    }
  }
  schema_validation_enabled = false
  lifecycle {
    ignore_changes = [name]
  }
}

// Role assignments: image-registry on vnet
resource "azapi_resource" "roleAssignment_imageRegistry_vnet" {
  type      = "Microsoft.Authorization/roleAssignments@2022-04-01"
  parent_id = azapi_resource.virtualNetwork.id
  name      = uuid()
  body = {
    properties = {
      principalId      = azapi_resource.imageRegistryIdentity.output.properties.principalId
      principalType    = "ServicePrincipal"
      roleDefinitionId = data.azapi_resource_list.roleDefinitions.output.aro_image_registry_operator
    }
  }
  schema_validation_enabled = false
  lifecycle {
    ignore_changes = [name]
  }
}

// Role assignments: aro-operator on subnets
resource "azapi_resource" "roleAssignment_aroOperator_masterSubnet" {
  type      = "Microsoft.Authorization/roleAssignments@2022-04-01"
  parent_id = azapi_resource.masterSubnet.id
  name      = uuid()
  body = {
    properties = {
      principalId      = azapi_resource.aroOperatorIdentity.output.properties.principalId
      principalType    = "ServicePrincipal"
      roleDefinitionId = data.azapi_resource_list.roleDefinitions.output.aro_service_operator
    }
  }
  schema_validation_enabled = false
  lifecycle {
    ignore_changes = [name]
  }
}

resource "azapi_resource" "roleAssignment_aroOperator_workerSubnet" {
  type      = "Microsoft.Authorization/roleAssignments@2022-04-01"
  parent_id = azapi_resource.workerSubnet.id
  name      = uuid()
  body = {
    properties = {
      principalId      = azapi_resource.aroOperatorIdentity.output.properties.principalId
      principalType    = "ServicePrincipal"
      roleDefinitionId = data.azapi_resource_list.roleDefinitions.output.aro_service_operator
    }
  }
  schema_validation_enabled = false
  lifecycle {
    ignore_changes = [name]
  }
}

// Role assignment: ARO RP first-party service principal on vnet
resource "azapi_resource" "roleAssignment_aroRP_vnet" {
  type      = "Microsoft.Authorization/roleAssignments@2022-04-01"
  parent_id = azapi_resource.virtualNetwork.id
  name      = uuid()
  body = {
    properties = {
      principalId      = data.azuread_service_principal.aroRP.object_id
      principalType    = "ServicePrincipal"
      roleDefinitionId = data.azapi_resource_list.roleDefinitions.output.aro_first_party_network
    }
  }
  schema_validation_enabled = false
  lifecycle {
    ignore_changes = [name]
  }
}

locals {
  cluster_domain            = "${var.resource_name}.aro-example.com"
  cluster_resource_group_id = "/subscriptions/${data.azapi_client_config.current.subscription_id}/resourcegroups/aro-${local.cluster_domain}-${var.location}"
}

resource "azapi_resource" "OpenShiftCluster" {
  type      = "Microsoft.RedHatOpenShift/openShiftClusters@2025-07-25"
  parent_id = azapi_resource.resourceGroup.id
  name      = var.resource_name
  location  = var.location
  identity {
    type         = "UserAssigned"
    identity_ids = [azapi_resource.clusterIdentity.id]
  }
  body = {
    properties = {
      apiserverProfile = {
        visibility = "Public"
      }
      clusterProfile = {
        domain               = local.cluster_domain
        fipsValidatedModules = "Enabled"
        resourceGroupId      = local.cluster_resource_group_id
      }
      consoleProfile = {}
      networkProfile = {
        podCidr     = "10.128.0.0/14"
        serviceCidr = "172.30.0.0/16"
        loadBalancerProfile = {
          managedOutboundIps = {
            count = 1
          }
        }
        preconfiguredNSG = "Disabled"
      }
      masterProfile = {
        vmSize           = "Standard_D8s_v3"
        subnetId         = azapi_resource.masterSubnet.id
        encryptionAtHost = "Enabled"
      }
      workerProfiles = [
        {
          name             = "worker"
          vmSize           = "Standard_D4s_v3"
          diskSizeGB       = 128
          subnetId         = azapi_resource.workerSubnet.id
          count            = 3
          encryptionAtHost = "Enabled"
        },
      ]
      platformWorkloadIdentityProfile = {
        platformWorkloadIdentities = {
          cloud-controller-manager = {
            resourceId = azapi_resource.cloudControllerManagerIdentity.id
          }
          ingress = {
            resourceId = azapi_resource.ingressIdentity.id
          }
          machine-api = {
            resourceId = azapi_resource.machineApiIdentity.id
          }
          disk-csi-driver = {
            resourceId = azapi_resource.diskCsiDriverIdentity.id
          }
          cloud-network-config = {
            resourceId = azapi_resource.cloudNetworkConfigIdentity.id
          }
          image-registry = {
            resourceId = azapi_resource.imageRegistryIdentity.id
          }
          file-csi-driver = {
            resourceId = azapi_resource.fileCsiDriverIdentity.id
          }
          aro-operator = {
            resourceId = azapi_resource.aroOperatorIdentity.id
          }
        }
      }
      ingressProfiles = [
        {
          name       = "default"
          visibility = "Public"
        },
      ]
    }
    tags = {
      key = "value"
    }
  }
  schema_validation_enabled = false
  ignore_casing             = false
  ignore_missing_property   = false
  timeouts {
    create = "90m"
    update = "90m"
    delete = "90m"
  }
  depends_on = [
    azapi_resource.roleAssignment_aroOperator_masterSubnet,
    azapi_resource.roleAssignment_aroOperator_workerSubnet,
    azapi_resource.roleAssignment_aroRP_vnet,
    azapi_resource.roleAssignment_ccm_masterSubnet,
    azapi_resource.roleAssignment_ccm_workerSubnet,
    azapi_resource.roleAssignment_cloudNetworkConfig_vnet,
    azapi_resource.roleAssignment_clusterIdentity_ccm,
    azapi_resource.roleAssignment_clusterIdentity_aroOperator,
    azapi_resource.roleAssignment_clusterIdentity_cloudNetworkConfig,
    azapi_resource.roleAssignment_clusterIdentity_diskCsiDriver,
    azapi_resource.roleAssignment_clusterIdentity_fileCsiDriver,
    azapi_resource.roleAssignment_clusterIdentity_imageRegistry,
    azapi_resource.roleAssignment_clusterIdentity_ingress,
    azapi_resource.roleAssignment_clusterIdentity_machineApi,
    azapi_resource.roleAssignment_fileCsiDriver_vnet,
    azapi_resource.roleAssignment_imageRegistry_vnet,
    azapi_resource.roleAssignment_ingress_masterSubnet,
    azapi_resource.roleAssignment_ingress_workerSubnet,
    azapi_resource.roleAssignment_machineApi_masterSubnet,
    azapi_resource.roleAssignment_machineApi_workerSubnet,
  ]
}

// ARO can remain transiently unstable immediately after provisioning completes.
// This pause gives the control plane time to settle before the acceptance test moves on.
resource "time_sleep" "wait_for_cluster_stabilization" {
  depends_on = [azapi_resource.OpenShiftCluster]

  // On destroy, Terraform removes this resource first, so this also inserts a delay
  // before the OpenShift cluster delete starts.
  create_duration  = "180s"
  destroy_duration = "180s"
}

```



## Arguments Reference

The following arguments are supported:

* `type` - (Required) The type of the resource. This should be set to `Microsoft.RedHatOpenShift/openShiftClusters@api-version`. The available api-versions for this resource are: [`2020-04-30`, `2021-09-01-preview`, `2022-04-01`, `2022-09-04`, `2023-04-01`, `2023-07-01-preview`, `2023-09-04`, `2023-11-22`, `2024-08-12-preview`, `2025-07-25`].

* `parent_id` - (Required) The ID of the azure resource in which this resource is created. The allowed values are:  
  `/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}`

* `name` - (Required) Specifies the name of the azure resource. Changing this forces a new resource to be created.

* `body` - (Required) Specifies the configuration of the resource. More information about the arguments in `body` can be found in the [Microsoft documentation](https://learn.microsoft.com/en-us/azure/templates/Microsoft.RedHatOpenShift/openShiftClusters?pivots=deployment-language-terraform).

For other arguments, please refer to the [azapi_resource](https://registry.terraform.io/providers/Azure/azapi/latest/docs/resources/resource) documentation.

## Import

 ```shell
 # Azure resource can be imported using the resource id, e.g.
 terraform import azapi_resource.example /subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.RedHatOpenShift/openShiftClusters/{resourceName}
 
 # It also supports specifying API version by using the resource id with api-version as a query parameter, e.g.
 terraform import azapi_resource.example /subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.RedHatOpenShift/openShiftClusters/{resourceName}?api-version=2025-07-25
 ```
