package services_test

import (
	"fmt"
	"testing"

	"github.com/Azure/terraform-provider-azapi/internal/acceptance"
	"github.com/Azure/terraform-provider-azapi/internal/acceptance/check"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccDataPlaneResource_keyVaultKey(t *testing.T) {
	data := acceptance.BuildTestData(t, "azapi_data_plane_resource", "test")
	r := DataPlaneResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.keyVaultKey(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(DataPlaneResource{}),
			),
		},
		{
			Config: r.keyVaultKeyUpdate(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(DataPlaneResource{}),
			),
		},
	})
}

func TestAccDataPlaneResource_foundryAgent(t *testing.T) {
	acceptance.SkipIfCoreAcctestsOnly(t, "Acctest subscription has no quota to run this test")
	data := acceptance.BuildTestData(t, "azapi_data_plane_resource", "test")
	r := DataPlaneResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.foundryAgent(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(DataPlaneResource{}),
			),
		},
		{
			Config: r.foundryAgentUpdate(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(DataPlaneResource{}),
			),
		},
	})
}

func (r DataPlaneResource) keyVaultKey(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azapi_data_plane_resource" "test" {
  type      = "Microsoft.KeyVault/vaults/keys@7.4"
  parent_id = trimsuffix(trimprefix(azapi_resource.vault.output.vaultUri, "https://"), "/")
  name      = "acctest%[2]s"
  body = {
    kty      = "RSA"
    key_size = 2048
    key_ops  = ["encrypt", "decrypt", "sign", "verify", "wrapKey", "unwrapKey"]
  }

  depends_on = [
    azapi_resource_action.add_accesspolicy
  ]
}`, r.keyVaultKeyTemplate(data), data.RandomString)
}

func (r DataPlaneResource) keyVaultKeyUpdate(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azapi_data_plane_resource" "test" {
  type      = "Microsoft.KeyVault/vaults/keys@7.4"
  parent_id = trimsuffix(trimprefix(azapi_resource.vault.output.vaultUri, "https://"), "/")
  name      = "acctest%[2]s"
  body = {
    kty      = "RSA"
    key_size = 2048
    key_ops  = ["encrypt", "decrypt", "sign", "verify", "wrapKey", "unwrapKey"]
    attributes = {
      enabled = true
    }
  }

  depends_on = [
    azapi_resource_action.add_accesspolicy
  ]
}`, r.keyVaultKeyTemplate(data), data.RandomString)
}

func (r DataPlaneResource) keyVaultKeyTemplate(data acceptance.TestData) string {
	return fmt.Sprintf(`
data "azapi_client_config" "current" {}

resource "azapi_resource" "resourceGroup" {
  type     = "Microsoft.Resources/resourceGroups@2020-06-01"
  name     = "acctest%[2]s"
  location = "%[1]s"
}

resource "azapi_resource" "vault" {
  type      = "Microsoft.KeyVault/vaults@2023-02-01"
  parent_id = azapi_resource.resourceGroup.id
  name      = "acctest%[2]s"
  location  = azapi_resource.resourceGroup.location
  body = {
    properties = {
      sku = {
        family = "A"
        name   = "standard"
      }
      tenantId                  = data.azapi_client_config.current.tenant_id
      enableSoftDelete          = true
      softDeleteRetentionInDays = 7
      enablePurgeProtection     = true
      accessPolicies            = []
    }
  }
  schema_validation_enabled = false
  lifecycle {
    ignore_changes = [body.properties.accessPolicies]
  }
  response_export_values = {
    "vaultUri" = "properties.vaultUri"
  }
}

resource "azapi_resource_action" "add_accesspolicy" {
  type        = "Microsoft.KeyVault/vaults/accessPolicies@2023-02-01"
  resource_id = "${azapi_resource.vault.id}/accessPolicies/add"
  method      = "PUT"
  body = {
    properties = {
      accessPolicies = [{
        tenantId = data.azapi_client_config.current.tenant_id
        objectId = data.azapi_client_config.current.object_id
        permissions = {
          keys = [
            "Get", "Create", "Delete", "List", "Update", "Restore", "Recover",
            "UnwrapKey", "WrapKey", "Purge", "Encrypt", "Decrypt", "Sign", "Verify"
          ]
        }
      }]
    }
  }
}


`, data.LocationPrimary, data.RandomString)
}

func (r DataPlaneResource) foundryAgent(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azapi_data_plane_resource" "test" {
  type      = "Microsoft.Foundry/agents@v1"
  parent_id = "acctest%s.services.ai.azure.com/api/projects/project%s"
  name      = "acctest-%s"

  depends_on = [
    azapi_resource.foundry_deployment,
    azapi_resource.foundry_project_user,
  ]

  retry = {
    error_message_regex  = ["PermissionDenied", "Unauthorized", "authorization", "context deadline exceeded"]
    interval_seconds     = 30
    max_interval_seconds = 180
  }

  body = {
    name = "acctest-%s"
    definition = {
      kind         = "prompt"
      model        = "gpt-4o"
      instructions = "You are an acceptance test agent"
    }
  }
}
`, r.foundryAgentTemplate(data), data.RandomString, data.RandomString, data.RandomString, data.RandomString)
}

func (r DataPlaneResource) foundryAgentUpdate(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azapi_data_plane_resource" "test" {
  type      = "Microsoft.Foundry/agents@v1"
  parent_id = "acctest%s.services.ai.azure.com/api/projects/project%s"
  name      = "acctest-%s"

  depends_on = [
    azapi_resource.foundry_deployment,
    azapi_resource.foundry_project_user,
  ]

  retry = {
    error_message_regex  = ["PermissionDenied", "Unauthorized", "authorization", "context deadline exceeded"]
    interval_seconds     = 30
    max_interval_seconds = 180
  }

  body = {
    definition = {
      kind         = "prompt"
      model        = "gpt-4o"
      instructions = "You are an updated acceptance test agent"
    }
  }
}
`, r.foundryAgentTemplate(data), data.RandomString, data.RandomString, data.RandomString)
}

func (r DataPlaneResource) foundryAgentTemplate(data acceptance.TestData) string {
	location := data.LocationPrimary
	if location == "" {
		location = "westus3"
	}

	return fmt.Sprintf(`
data "azapi_client_config" "current" {}

locals {
  foundry_project_user_role_definition_id = "/subscriptions/${data.azapi_client_config.current.subscription_id}/providers/Microsoft.Authorization/roleDefinitions/53ca6127-db72-4b80-b1b0-d745d6d5456d"
  foundry_project_user_role_assignment_id = format("%%s-%%s-%%s-%%s-%%s",
    substr(md5("foundry-project-user-%[1]s"), 0, 8),
    substr(md5("foundry-project-user-%[1]s"), 8, 4),
    substr(md5("foundry-project-user-%[1]s"), 12, 4),
    substr(md5("foundry-project-user-%[1]s"), 16, 4),
    substr(md5("foundry-project-user-%[1]s"), 20, 12)
  )
}

resource "azapi_resource" "resource_group" {
  type     = "Microsoft.Resources/resourceGroups@2021-04-01"
  name     = "acctest%[1]s"
  location = "%[2]s"
}

resource "azapi_resource" "foundry" {
  type                      = "Microsoft.CognitiveServices/accounts@2025-06-01"
  name                      = "acctest%[1]s"
  parent_id                 = azapi_resource.resource_group.id
  location                  = azapi_resource.resource_group.location
  schema_validation_enabled = false

  body = {
    kind = "AIServices"
    sku = {
      name = "S0"
    }
    identity = {
      type = "SystemAssigned"
    }
    properties = {
      disableLocalAuth       = false
      allowProjectManagement = true
      customSubDomainName    = "acctest%[1]s"
    }
  }
}

resource "azapi_resource" "foundry_deployment" {
  type      = "Microsoft.CognitiveServices/accounts/deployments@2023-05-01"
  name      = "gpt-5-mini"
  parent_id = azapi_resource.foundry.id

  body = {
    sku = {
      name     = "DataZoneStandard"
      capacity = 1
    }
    properties = {
      model = {
        format  = "OpenAI"
        name    = "gpt-5-mini"
        version = "2025-08-07"
      }
    }
  }
}

resource "azapi_resource" "foundry_project" {
  type                      = "Microsoft.CognitiveServices/accounts/projects@2025-06-01"
  name                      = "project%[1]s"
  parent_id                 = azapi_resource.foundry.id
  location                  = azapi_resource.foundry.location
  schema_validation_enabled = false

  body = {
    sku = {
      name = "S0"
    }
    identity = {
      type = "SystemAssigned"
    }
    properties = {
      displayName = "project"
      description = "Foundry project for acceptance test"
    }
  }
}

resource "azapi_resource" "foundry_project_user" {
  type      = "Microsoft.Authorization/roleAssignments@2022-04-01"
  name      = local.foundry_project_user_role_assignment_id
  parent_id = azapi_resource.foundry_project.id

  body = {
    properties = {
      principalId      = data.azapi_client_config.current.object_id
      roleDefinitionId = local.foundry_project_user_role_definition_id
    }
  }
}
`,
		data.RandomString,
		location,
	)
}
