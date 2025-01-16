---
layout: "azapi"
page_title: "Feature: Query Azure Resources with JMESPath"
description: |-
  This guide will cover how to use JMESPath queries to filter and format the output of Azure resources in Terraform.

---

The AzAPI resources and data sources use the `response_export_values` attribute to execute a [JMESPath query](http://jmespath.org/) on the response body. JMESPath is a query language for JSON, giving you the ability to select and modify data from the response. The result of the query is then stored in the `output` attribute.

This article covers how to use the features of JMESPath and gives examples of queries. 

## Prerequisites

- [Terraform AzAPI provider](https://registry.terraform.io/providers/azure/azapi) version 2.0.1 or later

Before introducing queries, take a look at the unmodified respone body of the `azapi_resource` data source, and we will use this response to demonstrate the queries.

```json
{
  "additionalCapabilities": null,
  "availabilitySet": null,
  "diagnosticsProfile": {
    "bootDiagnostics": {
      "enabled": true,
      "storageUri": "https://xxxxxx.blob.core.windows.net/"
    }
  },
  "osProfile": {
    "adminPassword": null,
    "adminUsername": "azureuser",
    "allowExtensionOperations": true,
    "computerName": "TestVM",
    "customData": null,
    "linuxConfiguration": {
      "disablePasswordAuthentication": true,
      "provisionVmAgent": true,
      "ssh": {
        "publicKeys": [
          {
            "keyData": "ssh-rsa AA***V stramer@contoso",
            "path": "/home/azureuser/.ssh/authorized_keys"
          }
        ]
      }
    },
    "secrets": [],
    "windowsConfiguration": null
  },
  // more properties
}
```

## Rename the properties

The `response_export_values` accepts a map where the key is the name for the result and the value is a JMESPath query string to filter the response. The following example renames the `osProfile.adminUsername` property to `admin`:

```hcl
data "azapi_resource" "virtualMachine" {
  type      = "Microsoft.Compute/virtualMachines@2020-06-01"
  name      = "myVirtualMachine"
  parent_id = "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/myResourceGroup"
  response_export_values = {
    admin   = "osProfile.adminUsername"
    ssh_key = "osProfile.linuxConfiguration.ssh.publicKeys[0].keyData"
  }
}
```

The `output` attribute contains the result of the queries:

```hcl
{
  admin = "azureuser"
  ssh_key = "ssh-rsa AA***V stramer@contoso"
}
```

Working with HCL object results, you can access properties from the top level with just the key. The `.` (subexpression) character is used to access properties of nested HCL objects. And the `[]` (index) character is used to access elements of an array.

For example, to use the `admin`:

```hcl
output "admin" {
  value = data.azapi_resource.virtualMachine.output.admin
}
```


## Export entire response

You can use the `@` character to export the entire response body. This is useful when you want to inspect the entire response or store it for later use.

Terraform configuration to export entire response:
```hcl
data "azapi_resource" "virtualMachine" {
  type      = "Microsoft.Compute/virtualMachines@2020-06-01"
  name      = "myVirtualMachine"
  parent_id = "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/myResourceGroup"
  response_export_values = {
    all = "@"
  }
}
```

The `output` attribute contains the result of the queries:

```
{
  all = {
    "additionalCapabilities": null,
    "availabilitySet": null,
    "diagnosticsProfile": {
      "bootDiagnostics": {
        "enabled": true,
        "storageUri": "https://xxxxxx.blob.core.windows.net/"
      }
    },
    "osProfile": {
      "adminPassword": null,
      "adminUsername": "azureuser",
      "allowExtensionOperations": true,
      "computerName": "TestVM",
      "customData": null,
      "linuxConfiguration": {
        "disablePasswordAuthentication": true,
        "provisionVmAgent": true,
        "ssh": {
          "publicKeys": [
            {
              "keyData": "ssh-rsa AA***V stramer@contoso",
              "path": "/home/azureuser/.ssh/authorized_keys"
            }
          ]
        }
      },
      "secrets": [],
      "windowsConfiguration": null
    },
    // more properties
  }
}
```

## Get properties from an array

Flattening an array is done with the `[]` JMESPath operator. All expressions after the `[]` operator are applied to each element in the current array.

For example, the following query gets the `Name` and `OS` properties from all virtual machines in a subscription:

```hcl
data "azapi_resource_list" "listVirtualMachinesBySubscription" {
  type      = "Microsoft.Compute/virtualMachines@2020-06-01"
  parent_id = "/subscriptions/00000000-0000-0000-0000-000000000000"
  response_export_values = {
    query = "value[].{Name:name, OS:properties.storageProfile.osDisk.osType, admin:properties.osProfile.adminUsername}"
  }
}
```

The `output` attribute contains the result of the query:

```hcl
{
  query = [
    {
      Name = "myVirtualMachine"
      OS = "Linux"
      admin = "azureuser"
    },
    {
      Name = "myVirtualMachine2"
      OS = "Windows"
      admin = "admin"
    }
  ]
}
```

## Filter arrays with boolean expressions

The other operation used to get data from an array is filtering. Filtering is done with the `[?...]` JMESPath operator. This operator takes a predicate as its contents. A predicate is any statement (including Boolean properties) that can be evaluated to either `true` or `false`. Expressions where the predicate evaluates to `true` are included in the output.

For example, the following query gets the `Name` and `admin` properties from all Linux virtual machines in a subscription:

```hcl
data "azapi_resource_list" "listLinuxVirtualMachinesBySubscription" {
  type      = "Microsoft.Compute/virtualMachines@2020-06-01"
  parent_id = "/subscriptions/00000000-0000-0000-0000-000000000000"
  response_export_values = {
    query = "value[?properties.storageProfile.osDisk.osType=='Linux'].{Name:name, admin:properties.osProfile.adminUsername}"
  }
}
```

The `output` attribute contains the result of the query:

```hcl
{
  query = [
    {
      Name = "myVirtualMachine"
      admin = "azureuser"
    }
  ]
}
```

Here is another example which demonstrates how to get the `id` of the `Contributor` role definition for a storage account, and then assign the role to a virtual machine:

```hcl
data "azapi_resource_list" "roleDefinitions" {
  type      = "Microsoft.Authorization/roleDefinitions@2022-04-01"
  parent_id = azapi_resource.storageAccount.id
  response_export_values = {
    contributorRoleDefinitionId = "value[?properties.roleName == 'Contributor'].id | [0]"
  }
}

resource "azapi_resource" "roleAssignment" {
  type      = "Microsoft.Authorization/roleAssignments@2022-04-01"
  parent_id = azapi_resource.storageAccount.id
  name      = "a08c80d9-771f-4769-bb4e-5c1a2873a189" // Random GUID
  body = {
    properties = {
      principalId      = azapi_resource.windowsVirtualMachine.identity[0].principal_id
      roleDefinitionId = data.azapi_resource_list.roleDefinitions.output.contributorRoleDefinitionId
    }
  }
}
```
