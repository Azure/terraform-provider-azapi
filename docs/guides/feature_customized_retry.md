---
layout: "azapi"
page_title: "Feature: Customized Retry Configuration"
description: |-
  This guide will cover how to use the customized retry configuration feature in the AzAPI provider.

---

The AzAPI provider can digest the intermittent API errors and retry the requests based on the customized retry configuration. This feature is useful when you need to handle the API errors gracefully and improve the reliability of the Terraform deployments.


## Prerequisites

- [Terraform AzAPI provider](https://registry.terraform.io/providers/azure/azapi) version 2.1.0 or later

## Customized Retry for Resource Creation

The virtual network link resource may not be available immediately after the virtual network is created. In this case, you can configure the customized retry configuration to handle the `ResourceNotFound` error and retry the request.

For example, the following configuration will create a virtual network link to the private DNS zone and retry the request when the `ResourceNotFound` error occurs:

```hcl
resource "azapi_resource" "privateDnsZoneLinkBlob" {
  type      = "Microsoft.Network/privateDnsZones/virtualNetworkLinks@2024-06-01"
  parent_id = azapi_resource.privateDnsZoneBlob.id
  name      = "blob"
  location  = "global"
  body = {
    properties = {
      registrationEnabled = false
      resolutionPolicy    = "Default"
      virtualNetwork = {
        id = azapi_resource.virtualNetwork.id
      }
    }
  }
  locks = [azapi_resource.virtualNetwork.id]
  retry = {
    error_message_regex = ["ResourceNotFound"]
  }
}
```

Above configuration is only used for demonstration purposes. From the `2.0.1` version, the AzAPI provider will automatically retry the GET requests when the `ResourceNotFound` error occurs after the resource creation. 

## Customized Retry for Resource Deletion

The private DNS zone may not be deleted immediately after the nested virtual network link is deleted. In this case, you can configure the customized retry configuration to handle the `CannotDeleteResource` error and retry the request.

```hcl
resource "azapi_resource" "privateDnsZoneQueue" {
  type      = "Microsoft.Network/privateDnsZones@2018-09-01"
  parent_id = azapi_resource.resourceGroup.id
  name      = "privatelink.queue.core.windows.net"
  location  = "global"
  body = {
    properties = {
    }
  }
  retry = {
    error_message_regex = ["CannotDeleteResource"]
  }
}
```

