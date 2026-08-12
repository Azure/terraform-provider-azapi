package services_test

import (
	"fmt"
	"testing"

	"github.com/Azure/terraform-provider-azapi/internal/acceptance"
	"github.com/Azure/terraform-provider-azapi/internal/acceptance/check"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccGenericResource_sqlJobAgentPrivateEndpoint(t *testing.T) {
	data := acceptance.BuildTestData(t, "azapi_resource", "test")
	r := GenericResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.sqlJobAgentPrivateEndpoint(data),
			ExternalProviders: map[string]resource.ExternalProvider{
				"time": {
					Source:            "hashicorp/time",
					VersionConstraint: "0.12.0",
				},
			},
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("output.properties.privateEndpointId").Exists(),
				check.That("data.azapi_resource.test").Key("output.properties.privateLinkServiceConnectionState.status").HasValue("Approved"),
			),
		},
	})
}

func (r GenericResource) sqlJobAgentPrivateEndpoint(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

data "azapi_client_config" "current" {}

locals {
  maintenance_configuration_id = "/subscriptions/${data.azapi_client_config.current.subscription_id}/providers/Microsoft.Maintenance/publicMaintenanceConfigurations/SQL_Default"
  private_endpoint_name        = "acctestpep%[2]s"
}

resource "azapi_resource" "agentServer" {
  type      = "Microsoft.Sql/servers@2024-05-01-preview"
  name      = "acctestagent%[2]s"
  parent_id = azapi_resource.resourceGroup.id
  location  = azapi_resource.resourceGroup.location
  body = {
    properties = {
      administratorLogin            = "acctestadmin"
      administratorLoginPassword    = "P@ssw0rd!1234"
      minimalTlsVersion             = "1.2"
      publicNetworkAccess           = "Enabled"
      restrictOutboundNetworkAccess = "Disabled"
      version                       = "12.0"
    }
  }
}

resource "azapi_resource" "agentDatabase" {
  type      = "Microsoft.Sql/servers/databases@2024-05-01-preview"
  name      = "acctestdb%[2]s"
  parent_id = azapi_resource.agentServer.id
  location  = azapi_resource.resourceGroup.location

  lifecycle {
    ignore_changes = [body.sku]
  }

  body = {
    properties = {
      collation                  = "SQL_Latin1_General_CP1_CI_AS"
      createMode                 = "Default"
      maintenanceConfigurationId = local.maintenance_configuration_id
    }
    sku = {
      name = "S0"
    }
  }
}

resource "azapi_resource" "jobAgent" {
  type      = "Microsoft.Sql/servers/jobAgents@2024-05-01-preview"
  name      = "acctestagent%[2]s"
  parent_id = azapi_resource.agentServer.id
  location  = azapi_resource.resourceGroup.location
  body = {
    properties = {
      databaseId = azapi_resource.agentDatabase.id
    }
    sku = {
      name = "JA100"
    }
  }
}

resource "azapi_resource" "targetServer" {
  type      = "Microsoft.Sql/servers@2024-05-01-preview"
  name      = "acctesttarget%[2]s"
  parent_id = azapi_resource.resourceGroup.id
  location  = azapi_resource.resourceGroup.location
  body = {
    properties = {
      administratorLogin            = "acctestadmin"
      administratorLoginPassword    = "P@ssw0rd!1234"
      minimalTlsVersion             = "1.2"
      publicNetworkAccess           = "Enabled"
      restrictOutboundNetworkAccess = "Disabled"
      version                       = "12.0"
    }
  }
}

resource "azapi_resource" "test" {
  type                      = "Microsoft.Sql/servers/jobAgents/privateEndpoints@2024-05-01-preview"
  name                      = local.private_endpoint_name
  parent_id                 = azapi_resource.jobAgent.id
  schema_validation_enabled = true
  ignore_missing_property   = true
  body = {
    properties = {
      targetServerAzureResourceId = azapi_resource.targetServer.id
    }
  }
  response_export_values = ["properties.privateEndpointId"]
}

# The private endpoint resource polls until the connection is approved, so this
# delay must run in parallel rather than depend on azapi_resource.test.
resource "time_sleep" "wait_for_private_endpoint" {
  create_duration = "1m"

  depends_on = [
    azapi_resource.jobAgent,
    azapi_resource.targetServer,
  ]
}

data "azapi_resource_list" "private_endpoint_connections" {
  type      = "Microsoft.Sql/servers/privateEndpointConnections@2024-05-01-preview"
  parent_id = azapi_resource.targetServer.id
  response_export_values = {
    private_endpoint_connection_id = "value[0].id"
  }

  depends_on = [time_sleep.wait_for_private_endpoint]
}

action "azapi_resource_action" "approve_private_endpoint" {
  config {
    type        = "Microsoft.Sql/servers/privateEndpointConnections@2024-05-01-preview"
    resource_id = data.azapi_resource_list.private_endpoint_connections.output.private_endpoint_connection_id
    method      = "PUT"
    body = {
      properties = {
        privateLinkServiceConnectionState = {
          status      = "Approved"
          description = "Approved by AzAPI acceptance test"
        }
      }
    }
  }
}

resource "terraform_data" "approve_private_endpoint" {
  input = data.azapi_resource_list.private_endpoint_connections.output.private_endpoint_connection_id

  lifecycle {
    action_trigger {
      events  = [before_create]
      actions = [action.azapi_resource_action.approve_private_endpoint]
    }
  }
}

data "azapi_resource" "test" {
  type                   = "Microsoft.Sql/servers/privateEndpointConnections@2024-05-01-preview"
  resource_id            = data.azapi_resource_list.private_endpoint_connections.output.private_endpoint_connection_id
  response_export_values = ["properties.privateLinkServiceConnectionState.status"]

  depends_on = [terraform_data.approve_private_endpoint]
}
`, r.template(data), data.RandomString)
}
