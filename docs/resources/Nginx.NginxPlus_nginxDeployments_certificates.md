---
subcategory: "Nginx.NginxPlus - NGINX Plus"
page_title: "nginxDeployments/certificates"
description: |-
  Manages a Certificate for an NGINX Deployment.
---

# Nginx.NginxPlus/nginxDeployments/certificates - Certificate for an NGINX Deployment

This article demonstrates how to use `azapi` provider to manage the Certificate for an NGINX Deployment resource in Azure.



## Example Usage

### default

```hcl
terraform {
  required_providers {
    azapi = {
      source = "Azure/azapi"
    }
  }
}

provider "azapi" {
  skip_provider_registration = false
}

data "azapi_client_config" "current" {}

variable "resource_name" {
  type    = string
  default = "acctest0001"
}

variable "location" {
  type    = string
  default = "westus"
}

resource "azapi_resource" "resourceGroup" {
  type     = "Microsoft.Resources/resourceGroups@2020-06-01"
  name     = var.resource_name
  location = var.location
}

resource "azapi_resource" "publicIPAddress" {
  type      = "Microsoft.Network/publicIPAddresses@2024-05-01"
  parent_id = azapi_resource.resourceGroup.id
  name      = "${var.resource_name}-pip"
  location  = var.location
  body = {
    properties = {
      ddosSettings = {
        protectionMode = "VirtualNetworkInherited"
      }
      idleTimeoutInMinutes     = 4
      publicIPAddressVersion   = "IPv4"
      publicIPAllocationMethod = "Static"
    }
    sku = {
      name = "Standard"
      tier = "Regional"
    }
  }
}

resource "azapi_resource" "userAssignedIdentity" {
  type      = "Microsoft.ManagedIdentity/userAssignedIdentities@2023-01-31"
  parent_id = azapi_resource.resourceGroup.id
  name      = "${var.resource_name}-identity"
  location  = var.location
  body      = {}
}

resource "azapi_resource" "virtualNetwork" {
  type      = "Microsoft.Network/virtualNetworks@2024-05-01"
  parent_id = azapi_resource.resourceGroup.id
  name      = "${var.resource_name}-vnet"
  location  = var.location
  body = {
    properties = {
      addressSpace = {
        addressPrefixes = ["10.0.0.0/16"]
      }
      dhcpOptions = {
        dnsServers = []
      }
      privateEndpointVNetPolicies = "Disabled"
      subnets                     = []
    }
  }
}

resource "azapi_resource" "vault" {
  type      = "Microsoft.KeyVault/vaults@2023-02-01"
  parent_id = azapi_resource.resourceGroup.id
  name      = "${substr(var.resource_name, 0, 12)}kv${substr(md5(var.resource_name), 0, 6)}"
  location  = var.location
  body = {
    properties = {
      accessPolicies = [{
        objectId = data.azapi_client_config.current.object_id
        permissions = {
          certificates = ["Get", "List"]
          keys         = ["Get"]
          secrets      = ["Get", "List"]
          storage      = []
        }
        tenantId = data.azapi_client_config.current.tenant_id
        }, {
        objectId = data.azapi_client_config.current.object_id
        permissions = {
          certificates = ["Get", "Create", "Delete", "List", "ManageContacts", "Purge", "Recover"]
          keys         = ["Get"]
          secrets      = ["Get", "Delete", "List", "Purge", "Recover", "Set"]
          storage      = []
        }
        tenantId = data.azapi_client_config.current.tenant_id
      }]
      createMode                   = "default"
      enableRbacAuthorization      = false
      enableSoftDelete             = true
      enabledForDeployment         = false
      enabledForDiskEncryption     = false
      enabledForTemplateDeployment = false
      publicNetworkAccess          = "Enabled"
      sku = {
        family = "A"
        name   = "standard"
      }
      softDeleteRetentionInDays = 7
      tenantId                  = data.azapi_client_config.current.tenant_id
    }
  }
}

resource "azapi_resource" "subnet" {
  type      = "Microsoft.Network/virtualNetworks/subnets@2024-05-01"
  parent_id = azapi_resource.virtualNetwork.id
  name      = "${var.resource_name}-subnet"
  body = {
    properties = {
      addressPrefix         = "10.0.2.0/24"
      defaultOutboundAccess = true
      delegations = [{
        name = "delegation"
        properties = {
          serviceName = "NGINX.NGINXPLUS/nginxDeployments"
        }
      }]
      privateEndpointNetworkPolicies    = "Disabled"
      privateLinkServiceNetworkPolicies = "Enabled"
      serviceEndpointPolicies           = []
      serviceEndpoints                  = []
    }
  }
}

resource "azapi_resource" "nginxDeployment" {
  type      = "Nginx.NginxPlus/nginxDeployments@2024-11-01-preview"
  parent_id = azapi_resource.resourceGroup.id
  name      = "${var.resource_name}-nginx"
  location  = var.location
  identity {
    type         = "UserAssigned"
    identity_ids = [azapi_resource.userAssignedIdentity.id]
  }
  body = {
    properties = {
      autoUpgradeProfile = {
        upgradeChannel = "stable"
      }
      enableDiagnosticsSupport = false
      networkProfile = {
        frontEndIPConfiguration = {
          publicIPAddresses = [{
            id = azapi_resource.publicIPAddress.id
          }]
        }
        networkInterfaceConfiguration = {
          subnetId = azapi_resource.subnet.id
        }
      }
      scalingProperties = {
        capacity = 10
      }
    }
    sku = {
      name = "standardv2_Monthly"
    }
  }
}

resource "azapi_resource" "secret" {
  type      = "Microsoft.KeyVault/vaults/secrets@2023-02-01"
  parent_id = azapi_resource.vault.id
  name      = "${var.resource_name}-cert"
  body = {
    properties = {
      value = "dummy-certificate-content"
    }
  }
  depends_on = [azapi_resource.vault]
}

resource "azapi_resource" "certificate" {
  type      = "Nginx.NginxPlus/nginxDeployments/certificates@2024-11-01-preview"
  parent_id = azapi_resource.nginxDeployment.id
  name      = "${var.resource_name}-cert"
  body = {
    properties = {
      certificateVirtualPath = "/opt/cert/server.cert"
      keyVaultSecretId       = azapi_resource.secret.output.properties.secretUriWithVersion
      keyVirtualPath         = "/opt/cert/server.key"
    }
  }
}


```



## Arguments Reference

The following arguments are supported:

* `type` - (Required) The type of the resource. This should be set to `Nginx.NginxPlus/nginxDeployments/certificates@api-version`. The available api-versions for this resource are: [`2021-05-01-preview`, `2022-08-01`, `2023-04-01`, `2023-09-01`, `2024-01-01-preview`, `2024-06-01-preview`, `2024-09-01-preview`, `2024-11-01-preview`].

* `parent_id` - (Required) The ID of the azure resource in which this resource is created. The allowed values are:  
  `/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Nginx.NginxPlus/nginxDeployments/{resourceName}`

* `name` - (Required) Specifies the name of the azure resource. Changing this forces a new resource to be created.

* `body` - (Required) Specifies the configuration of the resource. More information about the arguments in `body` can be found in the [Microsoft documentation](https://learn.microsoft.com/en-us/azure/templates/Nginx.NginxPlus/nginxDeployments/certificates?pivots=deployment-language-terraform).

For other arguments, please refer to the [azapi_resource](https://registry.terraform.io/providers/Azure/azapi/latest/docs/resources/resource) documentation.

## Import

 ```shell
 # Azure resource can be imported using the resource id, e.g.
 terraform import azapi_resource.example /subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Nginx.NginxPlus/nginxDeployments/{resourceName}/certificates/{resourceName}
 
 # It also supports specifying API version by using the resource id with api-version as a query parameter, e.g.
 terraform import azapi_resource.example /subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Nginx.NginxPlus/nginxDeployments/{resourceName}/certificates/{resourceName}?api-version=2024-11-01-preview
 ```
