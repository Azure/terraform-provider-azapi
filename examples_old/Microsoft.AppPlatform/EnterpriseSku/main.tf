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
  location = "west europe"
}

resource "azapi_resource" "spring" {
  type      = "Microsoft.AppPlatform/Spring@2022-01-01-preview"
  name      = "example"
  parent_id = azurerm_resource_group.test.id

  location = azurerm_resource_group.test.location
  body = jsonencode({
    sku = {
      tier = "Enterprise"
      name = "E0"
    }
  })
}

resource "azapi_resource" "app" {
  type      = "Microsoft.AppPlatform/Spring/apps@2022-01-01-preview"
  name      = "example"
  parent_id = azapi_resource.spring.id

  identity {
    type         = "SystemAssigned"
    identity_ids = []
  }
  body = jsonencode({
    properties = {
      addonConfigs = {
        applicationConfigurationService = {
          resourceId = azapi_resource.config.id
        }
        serviceRegistry = {
          resourceId = azapi_resource.sg.id
        }
      }
    }
  })
  response_export_values = ["properties.fqdn"]
}


resource "azapi_resource" "config" {
  type      = "Microsoft.AppPlatform/Spring/configurationServices@2022-01-01-preview"
  name      = "default"
  parent_id = azapi_resource.spring.id

  body = jsonencode({
    properties = {
      settings = {
        gitProperty = {
          repositories = [
            {
              label = "master"
              name  = "fake"
              patterns = [
                "app/dev"
              ]
              uri                   = "https://github.com/Azure-Samples/piggymetrics"
              hostKey               = "*" // "testkey"
              hostKeyAlgorithm      = "*" // "RSA"
              password              = "*" // "mypassword"
              privateKey            = "*" // "mykey"
              searchPaths           = ["dir1", "dir2"]
              strictHostKeyChecking = false
              username              = "*" // "admin"
            }
          ]
        }
      }
    }
  })
}

resource "azapi_resource" "sg" {
  type      = "Microsoft.AppPlatform/Spring/serviceRegistries@2022-01-01-preview"
  name      = "default"
  parent_id = azapi_resource.spring.id
}

resource "azapi_resource" "gateway" {
  type      = "Microsoft.AppPlatform/Spring/gateways@2022-01-01-preview"
  name      = "default"
  parent_id = azapi_resource.spring.id

  body = jsonencode({
    sku = {
      capacity = 2
      name     = "E0"
      tier     = "Enterprise"
    }
    properties = {
      apiMetadataProperties = {
        description   = "example description"
        documentation = "www.example.com/docs"
        serverUrl     = "example.com"
        title         = "test title"
        version       = "1.0"
      }
      corsProperties = {
        allowCredentials = true
        allowedHeaders   = ["*"]
        allowedMethods   = ["PUT"]
        allowedOrigins   = ["example.com"]
        exposedHeaders   = ["test"]
        maxAge           = 86400
      }
      httpsOnly = true
      public    = true
      resourceRequests = {
        cpu    = "1"
        memory = "1Gi"
      }
      ssoProperties = {
        clientId     = "*" // "00000000-0000-0000-0000-000000000000"
        clientSecret = "*" // "00000000-0000-0000-0000-000000000000"
        issuerUri    = "/issueToken"
        scope        = ["read"]
      }
    }
  })
}


resource "azapi_resource" "route" {
  type      = "Microsoft.AppPlatform/Spring/gateways/routeConfigs@2022-01-01-preview"
  name      = "example"
  parent_id = azapi_resource.gateway.id

  body = jsonencode({
    properties = {
      appResourceId = azapi_resource.app.id
      routes = [
        {
          description = "example description"
          filters     = ["StripPrefix=2", "RateLimit=1,1s"]
          order       = 1
          predicates  = ["Path=/api5/customer/**"]
          ssoEnabled  = false
          tags        = ["tag1", "tag2"]
          title       = "example route config"
          tokenRelay  = true
          uri         = "exampleuri"
        }
      ]
    }
  })

}

resource "azapi_resource" "builderService" {
  type      = "Microsoft.AppPlatform/Spring/buildServices/builders@2022-01-01-preview"
  name      = "example"
  parent_id = "${azapi_resource.spring.id}/buildServices/default"

  body = jsonencode({
    properties = {
      buildpackGroups = [
        {
          name = "mix"
          buildpacks = [
            {
              id = "tanzu-buildpacks/java-azure"
            }
          ]
        }
      ]
      stack = {
        id      = "io.buildpacks.stacks.bionic"
        version = "base"
      }
    }
  })
}

resource "azapi_resource" "binding" {
  type      = "Microsoft.AppPlatform/Spring/buildServices/builders/buildpackBindings@2022-01-01-preview"
  name      = "example"
  parent_id = azapi_resource.builderService.id

  body = jsonencode({
    properties = {
      bindingType = "ApplicationInsights"
      launchProperties = {
        properties = {
          abc = "def"
        }
        secrets = {
          connection-string = "*" // "XXXXXXXXXXXXXXXXX=XXXXXXXXXXXXX-XXXXXXXXXXXXXXXXXXX;"
        }
      }
    }
  })
}
