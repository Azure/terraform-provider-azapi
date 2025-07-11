---
layout: "azapi"
page_title: "Feature: Customized Retry Configuration"
description: |-
  This guide will cover how to use the customized retry configuration feature in the AzAPI provider.

---

## Why retry?

Sometimes, when managing cloud infrastructure, requests to the cloud provider may fail due to transient issues such as network problems, timeouts, eventual consistency, or rate limiting. In these cases, it can be beneficial to retry the request a few times before giving up.

The AzAPI provider can digest these intermittent API errors and retry the requests based on the customized retry configuration. This feature is useful when you need to handle the API errors gracefully and improve the reliability of the Terraform deployments.

## Prerequisites

- [Terraform AzAPI provider](https://registry.terraform.io/providers/azure/azapi) version 2.1.0 or later. Some features are only available from version 2.3.0.

## Retry configuration

There are two types of retry configurations available in the AzAPI provider:

1. Provider retry configuration
2. Resource-specific retry configuration

### Provider retry configuration

The provider retry configuration is a global configuration that applies to all resources managed by the provider. You can configure the provider retry behavior by setting the following provider values:

- `maximum_busy_retry_attempts`

This value controls the number of times the provider will retry a failed request. The default value is `3`.
A retry will be triggered if the request fails with HTTP 408, 429, 500, 502, 503, or 504.

In the case that the response header contains a `Retry-After` value, the provider will wait for the specified duration before retrying the request.

### Resource-specific retry configuration

In addition to the provider retry configuration, you can also configure the retry behavior for individual resources. This allows you to fine-tune the retry behavior for specific resources.

The resource-specific retry comes after the provider retry, that is to say that the provider retry will be attempted first, and if it fails or exceeds the maximum attempts, the resource-specific retry will be attempted.
Note that the resource-specific retry does not honour the `Retry-After` header and is exponential backoff based.

Resource specific retry is configured using the `retry` attribute.

If you configure a retry configuration, the maximum elapsed time for the retry will be set to the resource's timeout value for that operation (create, update, read, delete).

With `azapi_resource` and `azapi_data_plane_resource`, the provider performs a read operation after the resource has been created so that we can store the read-only values.

The schema of these retry attributes is as follows:

- `error_message_regex` - A list of regular expressions to match against error messages. If any of the regular expressions match, the request will be retried.
- `interval_seconds` - The initial number of seconds to wait before the 1st retry. The default value is `10`.
- `max_interval_seconds` - The maximum number of seconds to wait before retrying a request. The default value is `180`.

## Default resource-specific retry configuration

If you do not configure any retry values, the provider will use the following:

For the initial create/read/update/delete operation we will only retry on the provider's `maximum_busy_retry_attempts` value.

For the read-after-create, the provider will retry on HTTP 404 and 403 status codes up to the operation timeout. This is logical as, if we have just successfully created the resource, we should not be getting a 404 or 403 on any subsequent GET request.

## Example - Customized Retry for Resource Creation

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

## Example - Customized Retry for Resource Deletion

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
