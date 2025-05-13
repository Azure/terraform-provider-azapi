---
subcategory: "Microsoft.ConfidentialLedger - Confidential Ledger"
page_title: "ledgers"
description: |-
  Manages a Confidential Ledger.
---

# Microsoft.ConfidentialLedger/ledgers - Confidential Ledger

This article demonstrates how to use `azapi` provider to manage the Confidential Ledger resource in Azure.

## Example Usage

### default

```hcl
terraform {
  required_providers {
    azapi = {
      source = "Azure/azapi"
    }
    azurerm = {
      source = "hashicorp/azurerm"
    }
  }
}

provider "azurerm" {
  features {
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

data "azurerm_client_config" "current" {
}

resource "azapi_resource" "resourceGroup" {
  type     = "Microsoft.Resources/resourceGroups@2020-06-01"
  name     = var.resource_name
  location = var.location
}

resource "azapi_resource" "ledger" {
  type      = "Microsoft.ConfidentialLedger/ledgers@2022-05-13"
  parent_id = azapi_resource.resourceGroup.id
  name      = var.resource_name
  location  = var.location
  body = {
    properties = {
      aadBasedSecurityPrincipals = [
        {
          ledgerRoleName = "Administrator"
          principalId    = data.azurerm_client_config.current.object_id
          tenantId       = data.azurerm_client_config.current.tenant_id
        },
      ]
      certBasedSecurityPrincipals = [
        {
          cert           = "-----BEGIN CERTIFICATE-----MIIBsjCCATigAwIBAgIUZWIbyG79TniQLd2UxJuU74tqrKcwCgYIKoZIzj0EAwMwEDEOMAwGA1UEAwwFdXNlcjAwHhcNMjEwMzE2MTgwNjExWhcNMjIwMzE2MTgwNjExWjAQMQ4wDAYDVQQDDAV1c2VyMDB2MBAGByqGSM49AgEGBSuBBAAiA2IABBiWSo/j8EFit7aUMm5lF+lUmCu+IgfnpFD+7QMgLKtxRJ3aGSqgS/GpqcYVGddnODtSarNE/HyGKUFUolLPQ5ybHcouUk0kyfA7XMeSoUA4lBz63Wha8wmXo+NdBRo39qNTMFEwHQYDVR0OBBYEFPtuhrwgGjDFHeUUT4nGsXaZn69KMB8GA1UdIwQYMBaAFPtuhrwgGjDFHeUUT4nGsXaZn69KMA8GA1UdEwEB/wQFMAMBAf8wCgYIKoZIzj0EAwMDaAAwZQIxAOnozm2CyqRwSSQLls5r+mUHRGRyXHXwYtM4Dcst/VEZdmS9fqvHRCHbjUlO/+HNfgIwMWZ4FmsjD3wnPxONOm9YdVn/PRD7SsPRPbOjwBiE4EBGaHDsLjYAGDSGi7NJnSkA-----END CERTIFICATE-----"
          ledgerRoleName = "Administrator"
        },
      ]
      ledgerType = "Private"
    }
  }
  schema_validation_enabled = false
  response_export_values    = ["*"]
}


```



## Arguments Reference

The following arguments are supported:

* `type` - (Required) The type of the resource. This should be set to `Microsoft.ConfidentialLedger/ledgers@api-version`. The available api-versions for this resource are: [`2020-12-01-preview`, `2021-05-13-preview`, `2022-05-13`, `2022-09-08-preview`, `2023-01-26-preview`, `2023-06-28-preview`, `2024-07-09-preview`, `2024-09-19-preview`].

* `parent_id` - (Required) The ID of the azure resource in which this resource is created. The allowed values are:  
  `/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}`

* `name` - (Required) Specifies the name of the azure resource. Changing this forces a new resource to be created.

* `body` - (Required) Specifies the configuration of the resource. More information about the arguments in `body` can be found in the [Microsoft documentation](https://learn.microsoft.com/en-us/azure/templates/Microsoft.ConfidentialLedger/ledgers?pivots=deployment-language-terraform).

For other arguments, please refer to the [azapi_resource](https://registry.terraform.io/providers/Azure/azapi/latest/docs/resources/resource) documentation.

## Import

 ```shell
 # Azure resource can be imported using the resource id, e.g.
 terraform import azapi_resource.example /subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.ConfidentialLedger/ledgers/{resourceName}
 
 # It also supports specifying API version by using the resource id with api-version as a query parameter, e.g.
 terraform import azapi_resource.example /subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.ConfidentialLedger/ledgers/{resourceName}?api-version=2024-09-19-preview
 ```
