---
subcategory: ""
layout: "azapi"
page_title: "Azure Resource Action Data Source: azapi_resource_action"
description: |-
  Perform resource action which gets information from an existing resource.
---

# azapi_resource_action

This resource can perform resource action which gets information from an existing resource.
It's recommended to use `azapi_resource_action` data source to perform readonly action, please use `azapi_resource_action` resource,
if user wants to perform actions which change a resource's state.

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

resource "azurerm_resource_group" "example" {
  name     = "example-rg"
  location = "west europe"
}

resource "azurerm_automation_account" "example" {
  name                = "example-account"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location
  sku_name            = "Basic"
}

data "azapi_resource_action" "example" {
  type                   = "Microsoft.Automation/automationAccounts@2021-06-22"
  resource_id            = azurerm_automation_account.example.id
  action                 = "listKeys"
  response_export_values = ["*"]
}
```

Here's an example to use the `azapi_resource_action` data source to get a provider's permissions.

```hcl
provider "azurerm" {
  features {}
}

data "azurerm_client_config" "current" {}

data "azapi_resource_action" "test" {
  type        = "Microsoft.Resources/providers@2021-04-01"
  resource_id = "/subscriptions/${data.azurerm_client_config.current.subscription_id}/providers/Microsoft.Network"
  action      = "providerPermissions"
  method      = "GET"
}
```

Here's an example to use the `azapi_resource_action` data source to perform a provider action.

```hcl
terraform {
  required_providers {
    azapi = {
      source = "Azure/azapi"
    }
  }
}

resource "azapi_resource_action" "test" {
  type        = "Microsoft.Cache@2023-04-01"
  resource_id = "/subscriptions/00000000-0000-0000-0000-000000000000/providers/Microsoft.Cache"
  action      = "CheckNameAvailability"
  body = jsonencode({
    type = "Microsoft.Cache/Redis"
    name = "cacheName"
  })
}
```

## Arguments Reference

The following arguments are supported:

* `type` - (Required) It is in a format like `<resource-type>@<api-version>`. `<resource-type>` is the Azure resource type, for example, `Microsoft.Storage/storageAccounts`.
  `<api-version>` is version of the API used to manage this azure resource.

* `resource_id` - (Required) The ID of an existing azure source.

* `action` - (Optional) The name of the resource action. It's also possible to make Http requests towards the resource ID if leave this field empty.

---
* `body` - (Optional) A JSON object that contains the request body.

* `method` - (Optional) Specifies the Http method of the azure resource action. Allowed values are `POST` and `GET`. Defaults to `POST`.

* `response_export_values` - (Optional) A list of path that needs to be exported from response body.
  Setting it to `["*"]` will export the full response body.
  Here's an example. If it sets to `["keys"]`, it will set the following json to computed property `output`.
```
{
  "keys": [
    {
      "KeyName": "Primary",
      "Permissions": "Full",
      "Value": "nHGYNd******i4wdug=="
    },
    {
      "KeyName": "Secondary",
      "Permissions": "Full",
      "Value": "6yoCad******SLzKzg=="
    }
  ]
}
```

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the azure resource action.

* `output` - The output json containing the properties specified in `response_export_values`. Here are some examples to decode json and extract the value.
```hcl
// it will output "nHGYNd******i4wdug=="
output "primary_key" {
  value = jsondecode(azapi_resource_action.test.output).keys.0.Value
}

// it will output "6yoCad******SLzKzg=="
output "secondary_key" {
  value = jsondecode(azapi_resource_action.test.output).keys.1.Value
}
```

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `read` - (Defaults to 30 minutes) Used when retrieving the azure resource.
