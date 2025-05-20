---
subcategory: "Microsoft.Billing - Cost Management and Billing"
page_title: "billingAccounts/billingProfiles"
description: |-
  Manages a Billing Accounts Billing Profiles.
---

# Microsoft.Billing/billingAccounts/billingProfiles - Billing Accounts Billing Profiles

This article demonstrates how to use `azapi` provider to manage the Billing Accounts Billing Profiles resource in Azure.

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

variable "billing_account_id" {
  type        = string
  description = "Specify Billing Account Id for Billing Profile"
}

variable "payment_method_id" {
  type        = string
  description = "Specify Payment Method Id (For example: Credit Card and etc)"
}

variable "payment_sca_id" {
  type        = string
  description = "Specify Payment SCA Id for Payment Method Validation"
}

resource "azapi_resource" "billingProfile" {
  type      = "Microsoft.Billing/billingAccounts/billingProfiles@2024-04-01"
  parent_id = "/providers/Microsoft.Billing/billingAccounts/${var.billing_account_id}"
  name      = var.resource_name

  body = {
    properties = {
      billTo = {
        addressLine1   = "TestWay"
        city           = "Redmond"
        companyName    = "TestCompany"
        country        = "US"
        postalCode     = "12345-1234"
        region         = "WA"
        isValidAddress = true
      }
      displayName = var.resource_name
      enabledAzurePlans = [
        {
          skuId = "0001"
        }
      ]
      shipTo = {
        addressLine1   = "TestWay"
        city           = "Redmond"
        companyName    = "TestCompany"
        country        = "US"
        postalCode     = "12345-1234"
        region         = "WA"
        isValidAddress = true
      }
    }
  }

  create_headers = {
    "X-Ms-Payment-Method-Id" = var.payment_method_id
    "X-Ms-Payment-Sca-Id"    = var.payment_sca_id
  }

  schema_validation_enabled = false
  response_export_values    = ["*"]
}

```



## Arguments Reference

The following arguments are supported:

* `type` - (Required) The type of the resource. This should be set to `Microsoft.Billing/billingAccounts/billingProfiles@api-version`. The available api-versions for this resource are: [`2018-11-01-preview`, `2019-10-01-preview`, `2020-05-01`, `2024-04-01`].

* `parent_id` - (Required) The ID of the azure resource in which this resource is created. The allowed values are:  
  `//providers/Microsoft.Billing/billingAccounts/{resourceName}`

* `name` - (Required) Specifies the name of the azure resource. Changing this forces a new resource to be created.

* `body` - (Required) Specifies the configuration of the resource. More information about the arguments in `body` can be found in the [Microsoft documentation](https://learn.microsoft.com/en-us/azure/templates/Microsoft.Billing/billingAccounts/billingProfiles?pivots=deployment-language-terraform).

For other arguments, please refer to the [azapi_resource](https://registry.terraform.io/providers/Azure/azapi/latest/docs/resources/resource) documentation.

## Import

 ```shell
 # Azure resource can be imported using the resource id, e.g.
 terraform import azapi_resource.example //providers/Microsoft.Billing/billingAccounts/{resourceName}/billingProfiles/{resourceName}
 
 # It also supports specifying API version by using the resource id with api-version as a query parameter, e.g.
 terraform import azapi_resource.example //providers/Microsoft.Billing/billingAccounts/{resourceName}/billingProfiles/{resourceName}?api-version=2024-04-01
 ```
