---
layout: "azapi"
page_title: "AzAPI Provider: Authenticating via a Service Principal and OpenID Connect"
description: |-
  This guide will cover how to use a Service Principal (Shared Account) with OpenID Connect as authentication for the AzAPI Provider.

---

We recommend using either a Service Principal or Managed Service Identity when running Terraform non-interactively (such as when running Terraform in a CI server) - and authenticating using the Azure CLI when running Terraform locally.

## Setting up an Application and Service Principal in Azure

A Service Principal is a security principal within Azure Active Directory which can be granted access to resources within Azure Subscriptions. To authenticate with a Service Principal, you will need to create an Application object within Azure Active Directory, which you will use as a means of authentication, either [using a Client Secret](service_principal_client_secret.html), [a Client Certificate](service_principal_client_certificate.html), or OpenID Connect (which is documented in this guide). This can be done using the Azure Portal.

This guide will cover how to create an Application and linked Service Principal, and then how to assign federated identity credentials to the Application so that it can be used for authentication via OpenID Connect. Once that's done finally we're going to grant the Service Principal permission to manage resources in the Subscription - to do this we're going to assign `Contributor` rights to the Subscription - however, [it's possible to assign other permissions](https://azure.microsoft.com/documentation/articles/role-based-access-built-in-roles/) depending on your configuration.

### Creating the Application and Service Principal

We're going to create the Application in the Azure Portal - to do this navigate to [the **Azure Active Directory** overview](https://portal.azure.com/#blade/Microsoft_AAD_IAM/ActiveDirectoryMenuBlade/Overview) within the Azure Portal - [then select the **App Registration** blade](https://portal.azure.com/#blade/Microsoft_AAD_IAM/ActiveDirectoryMenuBlade/RegisteredApps/RegisteredApps/Overview). Click the **New registration** button at the top to add a new Application within Azure Active Directory. On this page, set the following values then press **Create**:

* **Name** - this is a friendly identifier and can be anything (e.g. "Terraform")
* **Supported Account Types** - this should be set to "Accounts in this organizational directory only (single-tenant)"
* **Redirect URI** - you should choose "Web" for the URI type. the actual value can be left blank

At this point the newly created Azure Active Directory application should be visible on-screen - if it's not, navigate to the [the **App Registration** blade](https://portal.azure.com/#blade/Microsoft_AAD_IAM/ActiveDirectoryMenuBlade/RegisteredApps/RegisteredApps/Overview) and select the Azure Active Directory application.

At the top of this page, you'll need to take note of the "Application (client) ID" and the "Directory (tenant) ID", which you can use for the values of `client_id` and `tenant_id` respectively.

### Configure Azure Active Directory Application to Trust a GitHub Repository

An application will need a federated credential specified for each GitHub Environment, Branch Name, Pull Request, or Tag based on your use case. For this example, we'll give permission to `main` branch workflow runs.

-> **Tip:** You can also configure the Application using the [azuread_application](https://registry.terraform.io/providers/hashicorp/azuread/latest/docs/resources/application) and [azuread_application_federated_identity_credential](https://registry.terraform.io/providers/hashicorp/azuread/latest/docs/resources/application_federated_identity_credential) resources in the AzureAD Terraform Provider.

#### Via the Portal

On the Azure Active Directory application page, go to **Certificates and secrets**.

In the Federated credentials tab, select Add credential. The Add a credential blade opens. In the **Federated credential scenario** drop-down box select **GitHub actions deploying Azure resources**.

Specify the **Organization** and **Repository** for your GitHub Actions workflow. For **Entity type**, select **Environment**, **Branch**, **Pull request**, or **Tag** and specify the value. The values must exactly match the configuration in the GitHub workflow. For our example, let's select **Branch** and specify `main`.

Add a **Name** for the federated credential.

The **Issuer**, **Audiences**, and **Subject identifier** fields autopopulate based on the values you entered.

Click **Add** to configure the federated credential.

### Via the Azure API

```sh
az rest --method POST \
        --uri https://graph.microsoft.com/beta/applications/${APP_OBJ_ID}/federatedIdentityCredentials \
        --headers Content-Type='application/json' \
        --body @body.json
```

Where the body is:

```json
{
  "name":"${REPO_NAME}-pull-request",
  "issuer":"https://token.actions.githubusercontent.com",
  "subject":"repo:${REPO_OWNER}/${REPO_NAME}:refs:refs/heads/main",
  "description":"${REPO_OWNER} PR",
  "audiences":["api://AzureADTokenExchange"]
}
```

See the [official documentation](https://docs.microsoft.com/en-us/azure/active-directory/develop/workload-identity-federation-create-trust-github) for more details.

### Granting the Application access to manage resources in your Azure Subscription

Once the Application exists in Azure Active Directory - we can grant it permissions to modify resources in the Subscription. To do this, [navigate to the **Subscriptions** blade within the Azure Portal](https://portal.azure.com/#blade/Microsoft_Azure_Billing/SubscriptionsBlade), then select the Subscription you wish to use, then click **Access Control (IAM)**, and finally **Add** > **Add role assignment**.

Firstly, specify a Role which grants the appropriate permissions needed for the Service Principal (for example, `Contributor` will grant Read/Write on all resources in the Subscription). There's more information about [the built in roles available here](https://azure.microsoft.com/en-gb/documentation/articles/role-based-access-built-in-roles/).

Secondly, search for and select the name of the Service Principal created in Azure Active Directory to assign it this role - then press **Save**.

### Configure Azure Active Directory Application to Trust a Generic Issuer

On the Azure Active Directory application page, go to **Certificates and secrets**.

In the Federated credentials tab, select **Add credential**. The 'Add a credential' blade opens. Refer to the instructions from your OIDC provider for completing the form, before choosing a **Name** for the federated credential and clicking the **Add** button.

## Configuring the Service Principal in Terraform

As we've obtained the credentials for this Service Principal - it's possible to configure them in a few different ways.

When storing the credentials as Environment Variables, for example:

```bash
export ARM_CLIENT_ID="00000000-0000-0000-0000-000000000000"
export ARM_SUBSCRIPTION_ID="00000000-0000-0000-0000-000000000000"
export ARM_TENANT_ID="00000000-0000-0000-0000-000000000000"
```

The provider will use the `ARM_OIDC_TOKEN` environment variable as an OIDC token. You can use this variable to specify the token provided by your OIDC provider.

When running Terraform in GitHub Actions, the provider will detect the `ACTIONS_ID_TOKEN_REQUEST_URL` and `ACTIONS_ID_TOKEN_REQUEST_TOKEN` environment variables set by the GitHub Actions runtime. You can also specify the `ARM_OIDC_REQUEST_TOKEN` and `ARM_OIDC_REQUEST_URL` environment variables.

For GitHub Actions workflows, you'll need to ensure the workflow has `write` permissions for the `id-token`.

```yaml
permissions:
  id-token: write
  contents: read
```

For more information about OIDC in GitHub Actions, see [official documentation](https://docs.github.com/en/actions/deployment/security-hardening-your-deployments/configuring-openid-connect-in-cloud-providers).

The following Terraform and Provider blocks can be specified - where `1.3.0` is the version of the AzAPI Provider that you'd like to use:

```hcl
terraform {
  required_providers {
    azapi = {
      source  = "azure/azapi"
      version = "=1.3.0"
    }
  }
}

provider "azapi" {
  use_oidc = true
}
```

When running Terraform in Azure Pipelines, there are two ways to authenticate using OIDC. 

The first way is to use the OIDC token.
You can specify the OIDC token using the `oidc_token` or `oidc_token_file_path` provider arguments. 
You can also specify the OIDC request token and URL using the environment variables `ARM_OIDC_TOKEN` and `ARM_OIDC_TOKEN_FILE_PATH`.

Here is an example of how to specify the OIDC token using the `oidc_token` provider argument:

```hcl
terraform {
  required_providers {
    azapi = {
      source = "azure/azapi"
    }
  }
}

provider "azapi" {
  oidc_token = "{OIDC Token}"

  // or use oidc_token_file_path
  // oidc_token_file_path = "{OIDC Token File Path}"

  use_oidc = true
}
```

And here is an example of azure-pipelines.yml file:

```yaml
  - task: AzureCLI@2
    displayName: Acc Tests with OIDC Token
    inputs:
      azureSubscription: 'azapi-oidc-test' // Azure Service Connection ID
      scriptType: 'pscore'
      scriptLocation: 'inlineScript'
      inlineScript: |
        $env:ARM_TENANT_ID = $env:tenantId
        $env:ARM_CLIENT_ID = $env:servicePrincipalId
        $env:ARM_OIDC_TOKEN = $env:idToken
        $env:ARM_USE_OIDC = 'true'
        terraform plan
      addSpnToEnvironment: true
```

The second way is to use the OIDC request token and URL. The provider will detect the `SYSTEM_OIDCREQUESTURI` environment variable set by the Azure Pipelines runtime and use it as the OIDC request URL.
You can specify the OIDC request token using the `oidc_request_token` provider argument or the environment variable `ARM_OIDC_REQUEST_TOKEN`.
And the Azure Service Connection ID must be specified using the `oidc_azure_service_connection_id` provider argument or the environment variable `ARM_OIDC_AZURE_SERVICE_CONNECTION_ID`.

Here is an example of how to specify the OIDC request token and URL using the `oidc_request_token` and `oidc_azure_service_connection_id` provider arguments:

```hcl
terraform {
  required_providers {
    azapi = {
      source = "azure/azapi"
    }
  }
}

provider "azapi" {
  oidc_request_token               = "{OIDC Request Token}"
  oidc_azure_service_connection_id = "{Azure Service Connection ID}"
  use_oidc                         = true
}
```

And here is an example of azure-pipelines.yml file:

```yaml
  - task: AzureCLI@2
    displayName: Acc Tests with OIDC Azure Pipeline
    inputs:
      azureSubscription: 'azapi-oidc-test' // Azure Service Connection ID
      scriptType: 'pscore'
      scriptLocation: 'inlineScript'
      inlineScript: |
        $env:ARM_TENANT_ID = $env:tenantId
        $env:ARM_CLIENT_ID = $env:servicePrincipalId
        $env:ARM_OIDC_REQUEST_TOKEN = "$(System.AccessToken)"
        $env:ARM_OIDC_AZURE_SERVICE_CONNECTION_ID = "azapi-oidc-test"
        $env:ARM_USE_OIDC = 'true'
        terraform plan
      addSpnToEnvironment: true
```

More information on [the fields supported in the Provider block can be found here](../index.html#argument-reference).

At this point running either `terraform plan` or `terraform apply` should allow Terraform to run using the Service Principal to authenticate.

---

It's also possible to configure these variables either in-line or from using variables in Terraform (as the `oidc_token`, `oidc_token_file_path`, or `oidc_request_token` and `oidc_request_url` are in this example), like so:

~> **NOTE:** We'd recommend not defining these variables in-line since they could easily be checked into Source Control.

```hcl
variable "oidc_token" {}
variable "oidc_token_file_path" {}
variable "oidc_request_token" {}
variable "oidc_request_url" {}

terraform {
  required_providers {
    azapi = {
      source  = "azure/azapi"
      version = "=1.3.0"
    }
  }
}

provider "azapi" {
  features {}

  subscription_id = "00000000-0000-0000-0000-000000000000"
  client_id       = "00000000-0000-0000-0000-000000000000"
  use_oidc        = true

  # for GitHub Actions
  oidc_request_token = var.oidc_request_token
  oidc_request_url   = var.oidc_request_url

  # for other generic OIDC providers, providing token directly
  oidc_token = var.oidc_token

  # for other generic OIDC providers, reading token from a file
  oidc_token_file_path = var.oidc_token_file_path

  tenant_id = "00000000-0000-0000-0000-000000000000"
}
```

More information on [the fields supported in the Provider block can be found here](../index.html#argument-reference).

At this point running either `terraform plan` or `terraform apply` should allow Terraform to run using the Service Principal to authenticate.
