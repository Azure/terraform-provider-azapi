import {
  to = azapi_data_plane_resource.example
  id = "exampleappconf.azconfig.io/kv/mykey|Microsoft.AppConfiguration/configurationStores/keyValues@1.0"
}

resource "azapi_data_plane_resource" "example" {
  type      = "Microsoft.AppConfiguration/configurationStores/keyValues@1.0"
  parent_id = "exampleappconf.azconfig.io"
  name      = "mykey"
  body = {
    content_type = ""
    value        = "myvalue"
  }
}
