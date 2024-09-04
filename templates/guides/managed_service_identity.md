---
layout: "azapi"
page_title: "AzAPI Provider: Authenticating via Managed Identity"
description: |-
  This guide will cover how to use managed identity for Azure resources as authentication for the AzAPI Provider.
---

## What is a managed identity?

[Managed identities for Azure resources](https://docs.microsoft.com/en-us/azure/active-directory/managed-identities-azure-resources/overview) can be used to authenticate to services that support Azure Active Directory (Azure AD) authentication. There are two types of managed identities: system-assigned and user-assigned. This article is based on system-assigned managed identities.

Managed identities work in conjunction with Azure Resource Manager (ARM), Azure AD, and the Azure Instance Metadata Service (IMDS). Azure resources that support managed identities expose an internal IMDS endpoint that the client can use to request an access token. No credentials are stored on the VM, and the only additional information needed to bootstrap the Terraform connection to Azure is the subscription ID and tenant ID.

Azure AD creates an AD identity when you configure an Azure resource to use a system-assigned managed identity. The configuration process is described in more detail, below. Azure AD then creates a service principal to represent the resource for role-based access control (RBAC) and access control (IAM). The lifecycle of a system-assigned identity is tied to the resource it is enabled for: it is created when the resource is created and it is automatically removed when the resource is deleted.

Before you can use the managed identity, it has to be configured. There are two steps:

1. Assign a role for the identity, associating it with the subscription that will be used to run Terraform. This step gives the identity permission to access Azure Resource Manager (ARM) resources.
1. Configure access control for one or more Azure resources. For example, if you use a key vault and a storage account, you will need to configure the vault and container separately.

Before you can create a resource with a managed identity and then assign an RBAC role, your account needs sufficient permissions. You need to be a member of the account **Owner** role, or have **Contributor** plus **User Access Administrator** roles.

Not all Azure services support managed identities, and availability varies by region. Configuration details vary slightly among services. For more information, see [Services that support managed identities for Azure resources](https://docs.microsoft.com/en-us/azure/active-directory/managed-identities-azure-resources/services-support-managed-identities).

## Configuring a VM to use a system-assigned managed identity

The (simplified) Terraform configuration below provisions a virtual machine with a system-assigned managed identity, and then grants the Contributor role to the identity.

```hcl
data "azurerm_subscription" "current" {}

resource "azapi_resource" "example" {
  type = "Microsoft.Compute/virtualMachines@2021-11-01"
  # ...

  identity {
    type = "SystemAssigned"
  }
}

data "azurerm_role_definition" "contributor" {
  name = "Contributor"
}

resource "azurerm_role_assignment" "example" {
  name               = azapi_resource.example.name
  scope              = data.azurerm_subscription.primary.id
  role_definition_id = "${data.azurerm_subscription.subscription.id}${data.azurerm_role_definition.contributor.id}"
  principal_id       = azapi_resource.example.identity.principal_id
}
```

## Configuring Terraform to use a managed identity

At this point we assume that managed identity is configured on the resource (e.g. virtual machine) being used - and that permissions have been assigned via Azure's Identity and Access Management system.

Terraform can be configured to use managed identity for authentication in one of two ways: using environment variables, or by defining the fields within the provider block.

### Configuring with environment variables

Using managed identity for Azure resources as authentication is enabled in `azapi` provider by default.

By default, Terraform will use the system assigned identity for authentication. To use a user assigned identity instead, you will need to specify the `ARM_CLIENT_ID` environment variable (equivalent to provider block argument [`client_id`](https://registry.terraform.io/providers/azure/azapi/latest/docs#client_id)) to the [client id](https://registry.terraform.io/providers/hashicorp/azurerm/latest/docs/resources/user_assigned_identity#client_id) of the identity.

By default, Terraform will use a well-known MSI endpoint to get the authentication token, which covers most use cases.

!> **Note:** we recommend against running Terraform inside of a Function App as the low memory ceiling can lead to Terraform being terminated and data (including the State File) being lost. Instead we’d recommend considering triggering an external process, such as Terraform Cloud or a CI System to run these longer-running more intensive processes - see [Terraform in Automation](https://learn.hashicorp.com/tutorials/terraform/automate-terraform) for more details.

In addition to a properly-configured management identity, Terraform needs to know the subscription ID and tenant ID to identify the full context for the Azure provider.

```shell
$ export ARM_SUBSCRIPTION_ID=159f2485-xxxx-xxxx-xxxx-xxxxxxxxxxxx
$ export ARM_TENANT_ID=72f988bf-xxxx-xxxx-xxxx-xxxxxxxxxxxx
$ export ARM_CLIENT_ID=xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx # only necessary for user assigned identity
```

A provider block is _technically_ optional when using environment variables. Even so, we recommend defining provider blocks so that you can pin or constrain the version of the provider being used, and configure other optional settings:

```hcl
terraform {
  required_providers {
    azapi = {
      source  = "azure/azapi"
      version = "=0.1.0"
    }
  }
}

provider "azapi" {
}
```

### Configuring with the provider block

It's also possible to configure a managed identity within the provider block:

```hcl
terraform {
  required_providers {
    azapi = {
      source  = "azure/azapi"
      version = "=0.1.0"
    }
  }
}

provider "azapi" {
  #...
}
```

More information on [the fields supported in the provider block can be found here](../index.html#argument-reference).
