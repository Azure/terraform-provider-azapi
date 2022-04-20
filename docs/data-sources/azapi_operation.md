---
subcategory: ""
layout: "azapi"
page_title: "Generic Azure Resource Operation Data Source: azapi_operation"
description: |-
  Perform resource function to get information from an existing azure resource
---

# azapi_operation

This resource can perform resource function to get information form any existing Azure resource manager resource.

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

data "azapi_operation" "example" {
  type                   = "Microsoft.Automation/automationAccounts@2021-06-22"
  resource_id            = azurerm_automation_account.example.id
  operation              = "listKeys"
  response_export_values = ["*"]
}
```

## Arguments Reference

The following arguments are supported:

* `type` - (Required) It is in a format like `<resource-type>@<api-version>`. `<resource-type>` is the Azure resource type, for example, `Microsoft.Storage/storageAccounts`.
  `<api-version>` is version of the API used to manage this azure resource.

* `resource_id` - (Required) The ID of an existing azure source.

* `operation` - (Optional) The name of the resource operation. It's also possible to make Http requests towards the resource ID if leave this field empty.

---
* `body` - (Optional) A JSON object that contains the request body.

* `method` - (Optional) Specifies the Http method of the azure resource operation. Allowed values are `POST`, `PATCH`, `GET`, `PUT`, `DELETE`, `CONNECT`, `HEAD`, `OPTIONS` and `TRACE`. Defaults to `POST`.

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

* `id` - The ID of the azure resource operation.

* `output` - The output json containing the properties specified in `response_export_values`. Here are some examples to decode json and extract the value.
```hcl
// it will output "nHGYNd******i4wdug=="
output "primary_key" {
  value = jsondecode(azapi_operation.test.output).keys.0.Value
}

// it will output "6yoCad******SLzKzg=="
output "secondary_key" {
  value = jsondecode(azapi_operation.test.output).keys.1.Value
}
```

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `read` - (Defaults to 30 minutes) Used when retrieving the azure resource.
