import {
  to = azapi_resource.example
  identity = {
    id   = "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/example-rg/providers/Microsoft.Network/virtualNetworks/example-vnet"
    type = "Microsoft.Network/virtualNetworks@2023-11-01"
  }
}

resource "azapi_resource" "example" {
  type      = "Microsoft.Network/virtualNetworks@2023-11-01"
  name      = "example-vnet"
  parent_id = "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/example-rg"
  location  = "westus"
  body = {
    properties = {
      addressSpace = {
        addressPrefixes = [
          "10.0.0.0/16"
        ]
      }
    }
  }
}
