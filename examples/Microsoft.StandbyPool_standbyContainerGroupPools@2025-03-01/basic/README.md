# Use Standby Pool for Container Group 

## Overview

This example shows how to setup Standby Pool for Container Group. Standby Pools allow you to create a pool of pre-provisioned, pre-initialized container groups that are ready to receive workloads instantly. This capability significantly reduces cold start latency and enhances responsiveness for bursty, event-driven, or latency-sensitive applications.

In this example, we will create a Container Group and a Standby Pool for it.

## Prerequisites

1. You need an Azure account with an active subscription. [Create an account for free.](https://azure.microsoft.com/free/?WT.mc_id=A261C142F)
2. [Install and configure Terraform.](https://learn.microsoft.com/en-us/azure/developer/terraform/quickstart-configure)
3. [Configure Permission for Standby Pool.](https://learn.microsoft.com/en-us/azure/container-instances/container-instances-standby-pool-configure-permissions)

## Implement the Terraform code

1. Create a directory in which to test and run the sample Terraform code and make it the current directory.

2. Create a `main.tf` file and copy the content of the [main.tf](./main.tf) file into it.

3. Update the below variables in the `main.tf` file:
    - `resource_name`: The name of the resource group and resources.
    - `location`: The location for the resource group and resources.
    - `subscription_id`: The subscription ID to use.

4. Initialize Terraform  
Run `terraform init` to initialize the Terraform deployment. This command downloads the Azure provider required to manage your Azure resources.

5. Login to Azure  
Run `az login` to log in to your Azure account.

6. Create a Terraform execution plan  
Run `terraform plan` to create an execution plan.

6. Apply a Terraform execution plan  
Run `terraform apply` to apply the execution plan to your cloud infrastructure.
