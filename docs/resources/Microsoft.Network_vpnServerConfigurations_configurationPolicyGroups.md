---
subcategory: "Microsoft.Network - Various Networking Services"
page_title: "vpnServerConfigurations/configurationPolicyGroups"
description: |-
  Manages a VPN Server Configuration Policy Group.
---

# Microsoft.Network/vpnServerConfigurations/configurationPolicyGroups - VPN Server Configuration Policy Group

This article demonstrates how to use `azapi` provider to manage the VPN Server Configuration Policy Group resource in Azure.

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

variable "resource_name" {
  type    = string
  default = "acctest0001"
}

variable "location" {
  type    = string
  default = "westeurope"
}

variable "radius_server_secret" {
  type        = string
  description = "The RADIUS server secret for VPN authentication"
  sensitive   = true
}

resource "azapi_resource" "resourceGroup" {
  type     = "Microsoft.Resources/resourceGroups@2020-06-01"
  name     = var.resource_name
  location = var.location
}

resource "azapi_resource" "vpnServerConfiguration" {
  type      = "Microsoft.Network/vpnServerConfigurations@2022-07-01"
  parent_id = azapi_resource.resourceGroup.id
  name      = var.resource_name
  location  = var.location
  body = {
    properties = {
      radiusClientRootCertificates = [
      ]
      radiusServerAddress = ""
      radiusServerRootCertificates = [
      ]
      radiusServerSecret = ""
      radiusServers = [
        {
          radiusServerAddress = "10.105.1.1"
          radiusServerScore   = 15
          radiusServerSecret  = var.radius_server_secret
        },
      ]
      vpnAuthenticationTypes = [
        "Radius",
      ]
      vpnClientIpsecPolicies = [
      ]
      vpnClientRevokedCertificates = [
      ]
      vpnClientRootCertificates = [
      ]
      vpnProtocols = [
        "OpenVPN",
        "IkeV2",
      ]
    }
  }
  schema_validation_enabled = false
  response_export_values    = ["*"]
}

resource "azapi_resource" "configurationPolicyGroup" {
  type      = "Microsoft.Network/vpnServerConfigurations/configurationPolicyGroups@2022-07-01"
  parent_id = azapi_resource.vpnServerConfiguration.id
  name      = var.resource_name
  body = {
    properties = {
      isDefault = false
      policyMembers = [
        {
          attributeType  = "RadiusAzureGroupId"
          attributeValue = "6ad1bd08"
          name           = "policy1"
        },
      ]
      priority = 0
    }
  }
  schema_validation_enabled = false
  response_export_values    = ["*"]
}


```



## Arguments Reference

The following arguments are supported:

* `type` - (Required) The type of the resource. This should be set to `Microsoft.Network/vpnServerConfigurations/configurationPolicyGroups@api-version`. The available api-versions for this resource are: [`2021-08-01`, `2022-01-01`, `2022-05-01`, `2022-07-01`, `2022-09-01`, `2022-11-01`, `2023-02-01`, `2023-04-01`, `2023-05-01`, `2023-06-01`, `2023-09-01`, `2023-11-01`, `2024-01-01`, `2024-03-01`, `2024-05-01`, `2024-07-01`, `2024-10-01`].

* `parent_id` - (Required) The ID of the azure resource in which this resource is created. The allowed values are:  
  `/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Network/vpnServerConfigurations/{resourceName}`

* `name` - (Required) Specifies the name of the azure resource. Changing this forces a new resource to be created.

* `body` - (Required) Specifies the configuration of the resource. More information about the arguments in `body` can be found in the [Microsoft documentation](https://learn.microsoft.com/en-us/azure/templates/Microsoft.Network/vpnServerConfigurations/configurationPolicyGroups?pivots=deployment-language-terraform).

For other arguments, please refer to the [azapi_resource](https://registry.terraform.io/providers/Azure/azapi/latest/docs/resources/resource) documentation.

## Import

 ```shell
 # Azure resource can be imported using the resource id, e.g.
 terraform import azapi_resource.example /subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Network/vpnServerConfigurations/{resourceName}/configurationPolicyGroups/{resourceName}
 
 # It also supports specifying API version by using the resource id with api-version as a query parameter, e.g.
 terraform import azapi_resource.example /subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Network/vpnServerConfigurations/{resourceName}/configurationPolicyGroups/{resourceName}?api-version=2024-10-01
 ```
