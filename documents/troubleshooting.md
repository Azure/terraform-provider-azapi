# Troubleshooting

### Q1 After apply, plan found there's still a change?

More context:
```
Terraform used the selected providers to generate the following execution plan. Resource actions are indicated with the following symbols:
  ~ update in-place

Terraform will perform the following actions:

  # azapi_resource.test will be updated in-place
  ~ resource "azapi_resource" "test" {
      ~ body                            = jsonencode(
          ~ {
              ~ properties = {
                  + accountKey  = "TOI************QqA=="
                    # (2 unchanged elements hidden)
      ~ output                          = jsonencode({}) -> (known after apply)
        tags                            = {}
        # (5 unchanged attributes hidden)
    }

Plan: 0 to add, 1 to change, 0 to destroy.
```

This happens when a property contains sensitive credential like a storage account's access key, the value won't be returned by design. 

Please add `ignore_missing_property = true` to the resource block and apply it.

Similarly, values in incorrect casing might be returned and cause a diff. Please use `ignore_casing = true` to suppress it.

### Q2 How to manage resources which can be managed by both parent resource API and its own API like `Microsoft.Network/virtualNetworks/subnets`?

More conext: Microsoft.Network/virtualNetworks/subnets has its own API and it can also be managed by Microsoft.Network/virtualNetworks API.

It's recommendded to manage Microsoft.Network/virtualNetworks/subnets in Microsoft.Network/virtualNetworks instead of using its own API.

The reason is, when update the vnet, its request body won't contain any subnets definitions, so existing subnets will be removed.


### Q3 How to use API/properties which is not in embeded schema?

Please add `schema_validation_enabled = false` to the resource block.


### Q4 How to manage properties which can only be set during update?

More context: Properties like `CMK` can't be set when create this resource. Because `CMK` depends on key vault key which depends on this resource's identity.

It can be solved by using `azapi_update_resource` to perform a multi-steps deployment, here's an [example](https://github.com/Azure/terraform-provider-azapi/tree/main/examples/Microsoft.ServiceBus/ServiceBusNamespace-CMK/main.tf).