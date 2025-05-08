# Introduction

This directory contains examples for the Azure provider. The examples are categorized by Azure resource type and scenario. Each example contains a `main.tf` file that demonstrates how to use the resource type in a specific scenario. The `README.md` file in each example directory provides additional information about the example.


## File Structure

The file structure is as follows:

```
AzureResourceType@APIVersion
├── senario_name
│   ├── main.tf (required)
│   ├── README.md (optional)


# Example
# Notice: The `/` in the directory name is replaced with `_` in the resource type name
Microsoft.Storage_storageAccounts@2021-04-01
├── basic
│   ├── main.tf
├── with_private_endpoint
│   ├── main.tf
│   ├── README.md
```


## Contributing

If you have an example that you would like to contribute, please follow the guidelines below:

1. Create a new directory for the resource type and API version if it does not already exist. The directory name should be in the format `AzureResourceType@APIVersion` and `/` in the resource type name should be replaced with `_`. For example, `Microsoft.Storage_storageAccounts@2021-04-01`. Please try to use the latest API version available.

2. Create a new directory for the scenario if it does not already exist. The directory name should be descriptive of the scenario. For example, `basic`, `with_private_endpoint`.

3. Create a `main.tf` file in the scenario directory that demonstrates how to use the resource type in the scenario.

    1. Template: Here is a template for the `main.tf` file:

    ```hcl
    terraform {
        required_providers {
            azapi = {
                source = "Azure/azapi"
            }
        }
    }

    provider "azapi" {
        skip_provider_registration = false
    }

    variable "resource_name" {
        type    = string
        default = "acctest0001"
    }

    variable "location" {
        type    = string
        default = "westeurope"
    }

    resource "azapi_resource" "resourceGroup" {
        type     = "Microsoft.Resources/resourceGroups@2020-06-01"
        name     = var.resource_name
        location = var.location
    }

    resource "azapi_resource" "automationAccount" {
        type      = "Microsoft.Automation/automationAccounts@2021-06-22"
        parent_id = azapi_resource.resourceGroup.id
        name      = var.resource_name
        location  = var.location
        body = {
            properties = {
                encryption = {
                    keySource = "Microsoft.Automation"
                }
                publicNetworkAccess = true
                sku = {
                    name = "Basic"
                }
            }
        }
    }

    ```

    2. Resource label: The resource label should be the singular form of the resource type name. For example, for type `Microsoft.Storage/storageAccounts@2021-04-01`, the resource label should be `storageAccount`; for type `Microsoft.Compute/virtualMachines@2021-03-01`, the resource label should be `virtualMachine`. You can add suffixes to the resource label if there are multiple resources of the same type in the scenario. For example, `storageAccount1`, `storageAccount2`.

    3. Resource name: The resource name should be set to the `resource_name` variable or computed based on the `resource_name` variable for example `name = "${var.resource_name}-ip"`. One exception is when the resource name is required to be a specific value, for example, when creating a subnet for bastion host, the subnet name should be `AzureBastionSubnet`. During the acceptance test, the variable `resource_name` will be replaced with a unique name.

    4. Location: The location should be set to the `location` variable. The location should be a location that is supported by the resource type.

    5. Identity: If the resource type supports identity, you should set the identity block in the resource. For example, `identity { type = "SystemAssigned" }`.

    6. Dependencies: If the resource has dependencies on other resources, you should create the dependent resources in the same `main.tf` file and use `azapi` resources to represent the dependent resources. You could find examples of how to create dependent resources in the existing examples.
        1. `subscription_id`: `azapi_client_config` data source is used to get the subscription ID, for example, `subscription_id = data.azapi_client_config.current.subscription_id`.

    5. Other providers: If other providers are required to create the resource, you should add the required providers block in the `main.tf` file. For example, if the resource requires the `azurerm` provider, you should add the following block to the `main.tf` file:

    ```hcl
    terraform {
        required_providers {
                azapi = {
                    source = "Azure/azapi"  
                }
                azurerm = {
                    source = "hashicorp/azurerm"
                    version = "=4.20.0"
                }
        }
    }
    ```


4. (Optional) Create a `README.md` file in the scenario directory that provides additional information about the example if the example is complex or requires additional context.

5. Verify that the example works by running the acceptance tests. You can run the acceptance tests by running the following command:

    ```shell
    make example-test TARGET=AzureResourceType@APIVersion
    ```

    It will run acceptance tests for all scenarios in the specified example. 

    For example, to run the acceptance tests for the `basic` scenario of the `Microsoft.Storage_storageAccounts@2021-04-01` example, you can run the following command:

    ```shell
    make example-test TARGET=Microsoft_Storage_storageAccounts@2021-04-01
    ```

    For example, to run multiple examples, you can run the following command:

    ```shell
    make example-test TARGET=Microsoft_Storage_storageAccounts@2021-04-01,Microsoft_Compute_virtualMachines@2021-03-01
    ```

    Behind the scenes, it will run the unit test `TestAccExamples_Selected` with environment variable `ARM_TEST_EXAMPLES` set to the target example. The test will create the resources defined in the example and verify that the resources are created successfully, then it will destroy the resources.
    
    It's also acceptable to run the below terraform commands to test the examples:

    ```shell
    cd examples/AzureResourceType@APIVersion/scenario_name
    terraform init
    terraform apply
    terraform plan # it's expected to have no changes
    terraform destroy
    ```

6. Run `make docs` to generate the documentation for the examples. This will generate the documentation files in the `docs` directory. The documentation files will be generated based on the `main.tf` file in each example directory. 

7. Create a pull request with the changes. Please also include the output of the acceptance tests in the pull request description.


## Documentation Examples 

These examples are mostly used for documentation, but can also be run/tested manually via the Terraform CLI.

The document generation tool looks for files in the following locations by default. All other *.tf files besides the ones mentioned below are ignored by the documentation tool. This is useful for creating examples that can run and/or ar testable even if some parts are not relevant for the documentation.

* **provider/provider.tf** example file for the provider index page
* **data-sources/`full data source name`/data-source.tf** example file for the named data source page
* **resources/`full resource name`/resource.tf** example file for the named resource page
* **ephemeral-resources/`full resource name`/ephemeral-resource.tf** example file for the named ephemeral resource page

