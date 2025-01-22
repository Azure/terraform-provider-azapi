---
layout: "azapi"
page_title: "Authentication: Authenticating via the Azure CLI"
description: |-
  This guide will cover how to use the Azure CLI as authentication for the AzAPI Provider.

---

# AzAPI Provider: Authenticating using the Azure CLI

## Important Notes about Authenticating using the Azure CLI

* Terraform only supports authenticating using the `az` CLI (and this must be available on your PATH) - authenticating using the older `azure` CLI or PowerShell Cmdlets are not supported.
* Authenticating via the Azure CLI is only supported when using a User Account. If you're using a Service Principal (for example via `az login --service-principal`) you should instead authenticate via the Service Principal directly (either using a [Client Secret](service_principal_client_secret.md) or a [Client Certificate](service_principal_client_certificate.md)).

---

## Logging into the Azure CLI

~> **Note**: If you're using the **China**, **German** or **Government** Azure Clouds - you'll need to first configure the Azure CLI to work with that Cloud.  You can do this by running:

```shell
$ az cloud set --name AzureChinaCloud|AzureGermanCloud|AzureUSGovernment
```

---

Login to the Azure CLI using:

```shell
$ az login
```

Once logged in - it's possible to list the Subscriptions associated with the account via:

```shell
$ az account list --out table
```

The output (similar to below) will display one or more Subscriptions - with the `SubscriptionId` column being the `SUBSCRIPTION_ID` field referenced below.

```bash
Name                           CloudName    SubscriptionId                        State    IsDefault
-----------------------------  -----------  ------------------------------------  -------  -----------
PAYG Subscription              AzureCloud   00000000-0000-0000-0000-000000000000  Enabled  False
Contoso Sales                  AzureCloud   00000000-0000-1000-0000-000000000000  Enabled  False
Contoso Dev                    AzureCloud   00000000-0000-1000-2000-000000000000  Enabled  True
Contoso Dogfood                AzureCloud   00000000-3000-0000-0070-000000000000  Enabled  False
Contoso Prod                   AzureCloud   00000000-0400-0000-0070-000000000000  Enabled  False

```

Should you have more than one Subscription, you can specify the Subscription to use via the following command:

```bash
$ az account set --subscription="SUBSCRIPTION_ID"
```

---

## Configuring Azure CLI authentication in Terraform

Now that we're logged into the Azure CLI - we can configure Terraform to use these credentials.

To configure Terraform to use the Default Subscription defined in the Azure CLI - we can use the following Provider block:

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


At this point running either `terraform plan` or `terraform apply` should allow Terraform to run using the Azure CLI to authenticate.

---

It's also possible to configure Terraform to use a specific Subscription - for example:

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
  subscription_id = "00000000-0000-0000-0000-000000000000"
}
```

More information on [the fields supported in the Provider block can be found here](../index.html#argument-reference).

At this point running either `terraform plan` or `terraform apply` should allow Terraform to run using the Azure CLI to authenticate.

---

If you're looking to use Terraform across Tenants - it's possible to do this by configuring the Tenant ID field in the Provider block, as shown below:

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

  subscription_id = "00000000-0000-0000-0000-000000000000"
  tenant_id       = "11111111-1111-1111-1111-111111111111"
}
```

More information on [the fields supported in the Provider block can be found here](../index.html#argument-reference).

At this point running either `terraform plan` or `terraform apply` should allow Terraform to run using the Azure CLI to authenticate.
