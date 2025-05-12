---
subcategory: "Microsoft.Media - Media Services"
page_title: "mediaServices/contentKeyPolicies"
description: |-
  Manages a Media Services Content Key Policies.
---

# Microsoft.Media/mediaServices/contentKeyPolicies - Media Services Content Key Policies

This article demonstrates how to use `azapi` provider to manage the Media Services Content Key Policies resource in Azure.

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

resource "azapi_resource" "resourceGroup" {
  type     = "Microsoft.Resources/resourceGroups@2020-06-01"
  name     = var.resource_name
  location = var.location
}

resource "azapi_resource" "storageAccount" {
  type      = "Microsoft.Storage/storageAccounts@2021-09-01"
  parent_id = azapi_resource.resourceGroup.id
  name      = var.resource_name
  location  = var.location
  body = {
    kind = "StorageV2"
    properties = {
      accessTier                   = "Hot"
      allowBlobPublicAccess        = true
      allowCrossTenantReplication  = true
      allowSharedKeyAccess         = true
      defaultToOAuthAuthentication = false
      encryption = {
        keySource = "Microsoft.Storage"
        services = {
          queue = {
            keyType = "Service"
          }
          table = {
            keyType = "Service"
          }
        }
      }
      isHnsEnabled      = false
      isNfsV3Enabled    = false
      isSftpEnabled     = false
      minimumTlsVersion = "TLS1_2"
      networkAcls = {
        defaultAction = "Allow"
      }
      publicNetworkAccess      = "Enabled"
      supportsHttpsTrafficOnly = true
    }
    sku = {
      name = "Standard_GRS"
    }
  }
  schema_validation_enabled = false
  response_export_values    = ["*"]
}

resource "azapi_resource" "mediaService" {
  type      = "Microsoft.Media/mediaServices@2021-11-01"
  parent_id = azapi_resource.resourceGroup.id
  name      = var.resource_name
  location  = var.location
  body = {
    properties = {
      publicNetworkAccess = "Enabled"
      storageAccounts = [
        {
          id   = azapi_resource.storageAccount.id
          type = "Primary"
        },
      ]
    }
  }
  schema_validation_enabled = false
  response_export_values    = ["*"]
}

resource "azapi_resource" "contentKeyPolicy" {
  type      = "Microsoft.Media/mediaServices/contentKeyPolicies@2022-08-01"
  parent_id = azapi_resource.mediaService.id
  name      = var.resource_name
  body = {
    properties = {
      description = "My Policy Description"
      options = [
        {
          configuration = {
            "@odata.type" = "#Microsoft.Media.ContentKeyPolicyClearKeyConfiguration"
          }
          name = "ClearKeyOption"
          restriction = {
            "@odata.type" = "#Microsoft.Media.ContentKeyPolicyTokenRestriction"
            audience      = "urn:audience"
            issuer        = "urn:issuer"
            primaryVerificationKey = {
              "@odata.type" = "#Microsoft.Media.ContentKeyPolicySymmetricTokenKey"
              keyValue      = "AAAAAAAAAAAAAAAAAAAAAA=="
            }
            requiredClaims = [
            ]
            restrictionTokenType = "Swt"
          }
        },
      ]
    }
  }
  schema_validation_enabled = false
  response_export_values    = ["*"]
}


```



## Arguments Reference

The following arguments are supported:

* `type` - (Required) The type of the resource. This should be set to `Microsoft.Media/mediaServices/contentKeyPolicies@api-version`. The available api-versions for this resource are: [`2018-03-30-preview`, `2018-06-01-preview`, `2018-07-01`, `2019-05-01-preview`, `2020-05-01`, `2021-06-01`, `2021-11-01`, `2022-08-01`, `2023-01-01`].

* `parent_id` - (Required) The ID of the azure resource in which this resource is created. The allowed values are:  
  `/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Media/mediaServices/{resourceName}`

* `name` - (Required) Specifies the name of the azure resource. Changing this forces a new resource to be created.

* `body` - (Required) Specifies the configuration of the resource. More information about the arguments in `body` can be found in the [Microsoft documentation](https://learn.microsoft.com/en-us/azure/templates/Microsoft.Media/mediaServices/contentKeyPolicies?pivots=deployment-language-terraform).

For other arguments, please refer to the [azapi_resource](https://registry.terraform.io/providers/Azure/azapi/latest/docs/resources/resource) documentation.

## Import

 ```shell
 # Azure resource can be imported using the resource id, e.g.
 terraform import azapi_resource.example /subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Media/mediaServices/{resourceName}/contentKeyPolicies/{resourceName}
 
 # It also supports specifying API version by using the resource id with api-version as a query parameter, e.g.
 terraform import azapi_resource.example /subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Media/mediaServices/{resourceName}/contentKeyPolicies/{resourceName}?api-version=2023-01-01
 ```
