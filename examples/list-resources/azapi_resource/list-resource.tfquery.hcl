list "azapi_resource" "example" {
  provider = azapi
  config {
    type = "Microsoft.Storage/storageAccounts@2021-04-01"
    parent_id = "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/example-rg"
  }
}