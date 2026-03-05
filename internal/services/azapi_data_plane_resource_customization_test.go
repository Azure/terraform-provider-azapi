package services_test

import (
	"fmt"
	"os"
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
	host := os.Getenv("ARM_TEST_FOUNDRY_HOST")
	projectName := os.Getenv("ARM_TEST_FOUNDRY_PROJECT_NAME")
	modelName := os.Getenv("ARM_TEST_FOUNDRY_MODEL")
	if host == "" || projectName == "" || modelName == "" {
		t.Skip("Skipping as Foundry env vars ARM_TEST_FOUNDRY_HOST, ARM_TEST_FOUNDRY_PROJECT_NAME and ARM_TEST_FOUNDRY_MODEL must be specified")
	}

	data := acceptance.BuildTestData(t, "azapi_data_plane_resource", "test")
	r := DataPlaneResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.foundryAgent(data, host, projectName, modelName),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(DataPlaneResource{}),
			),
		},
		{
			Config: r.foundryAgentUpdate(data, host, projectName, modelName),
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

func (r DataPlaneResource) foundryAgent(data acceptance.TestData, host, projectName, modelName string) string {
	return fmt.Sprintf(`
resource "azapi_data_plane_resource" "test" {
  type      = "Microsoft.Foundry/agents@v1"
  parent_id = "%s/api/projects/%s"
  name      = "acctest-%s"

  body = {
    name         = "acctest-%s"
    definition = {
      kind         = "prompt"
      model        = "%s"
      instructions = "You are an acceptance test agent"
    }
  }
}
`, host, projectName, data.RandomString, data.RandomString, modelName)
}

func (r DataPlaneResource) foundryAgentUpdate(data acceptance.TestData, host, projectName, modelName string) string {
	return fmt.Sprintf(`
resource "azapi_data_plane_resource" "test" {
  type      = "Microsoft.Foundry/agents@v1"
  parent_id = "%s/api/projects/%s"
  name      = "acctest-%s"

  body = {
    definition = {
      kind         = "prompt"
      model        = "%s"
      instructions = "You are an updated acceptance test agent"
    }
  }
}
`, host, projectName, data.RandomString, modelName)
}
