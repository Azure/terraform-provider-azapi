terraform {
  required_providers {
    azapi = {
      source = "Azure/azapi"
    }
  }
}

provider "azurerm" {
  features {}
}

provider "azapi" {
}

resource "azurerm_resource_group" "test" {
  name     = "example-resource-group"
  location = "eastus"
}

resource "azurerm_log_analytics_workspace" "test" {
  name                = "example"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku                 = "PerGB2018"
  retention_in_days   = 30
}

resource "azurerm_container_registry" "test" {
  name                = "example"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  sku                 = "Premium"
  admin_enabled       = true
}

resource "azapi_resource" "run_acr_task" {
  type      = "Microsoft.ContainerRegistry/registries/taskRuns@2019-06-01-preview"
  name      = "example"
  parent_id = azurerm_container_registry.test.id
  location  = azurerm_resource_group.test.location
  body = jsonencode({
    properties = {
      runRequest = {
        type           = "DockerBuildRequest"
        sourceLocation = "https://github.com/Azure-Samples/aci-helloworld.git#master"
        dockerFilePath = "Dockerfile"
        platform = {
          os = "Linux"
        }
        imageNames = ["helloworld:{{.Run.ID}}", "helloworld:latest"]
      }
    }
  })
}

resource "azapi_resource" "container_app_environment" {
  type      = "Microsoft.App/managedEnvironments@2022-03-01"
  name      = "example"
  parent_id = azurerm_resource_group.test.id
  location  = azurerm_resource_group.test.location
  body = jsonencode({
    properties = {
      appLogsConfiguration = {
        destination = "log-analytics"
        logAnalyticsConfiguration = {
          customerId = azurerm_log_analytics_workspace.test.workspace_id
          sharedKey  = azurerm_log_analytics_workspace.test.primary_shared_key
        }
      }
    }
  })

  // properties/appLogsConfiguration/logAnalyticsConfiguration/sharedKey contains credential which will not be returned,
  // using this property to suppress plan-diff
  ignore_missing_property = true
}

resource "azapi_resource" "container_app" {
  type      = "Microsoft.App/containerApps@2022-01-01-preview"
  name      = "example"
  parent_id = azurerm_resource_group.test.id
  location  = azurerm_resource_group.test.location
  body = jsonencode({
    properties = {
      managedEnvironmentId = azapi_resource.container_app_environment.id
      configuration = {
        ingress = {
          targetPort = 3333
          external   = true
        }
        /*
        It supports authentication with service principle, by replacing the `admin_username` and `admin_password`
        with client id and client secret.
        */
        secrets = [
          {
            name  = "registry-password"
            value = azurerm_container_registry.test.admin_password
          }
        ]
        registries = [
          {
            passwordSecretRef = "registry-password"
            server            = azurerm_container_registry.test.login_server
            username          = azurerm_container_registry.test.admin_username
          }
        ]
      }
      template = {
        containers = [
          {
            image = "${azurerm_container_registry.test.login_server}/helloworld:latest",
            name  = "helloworld"
          }
        ]
      }
    }
  })

  // properties/configuration/secrets/value contains credential which will not be returned,
  // using this property to suppress plan-diff
  ignore_missing_property = true
  depends_on = [
    azapi_resource.run_acr_task
  ]
}
