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
