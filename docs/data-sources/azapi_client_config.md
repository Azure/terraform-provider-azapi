---
subcategory: ""
layout: "azapi"
page_title: "Client Config Data Source: azapi_client_config"
description: |-
  Gets information about the configuration of the azapi provider.
---

# azapi_client_config

Use this data source to access the configuration of the azapi provider.

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

data "azapi_client_config" "current" {
}

output "subscription_id" {
  value = data.azapi_client_config.current.subscription_id
}

output "tenant_id" {
  value = data.azapi_client_config.current.tenant_id
}
```

## Arguments Reference

There are no arguments available for this data source.

## Attributes Reference

* `tenant_id` is set to the Azure Tenant ID.

* `subscription_id` is set to the Azure Subscription ID.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `read` - (Defaults to 30 minutes) Used when retrieving the azure resource.
