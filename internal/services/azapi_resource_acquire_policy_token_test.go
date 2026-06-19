package services_test

import (
	"context"
	"fmt"
	"regexp"
	"testing"

	"github.com/Azure/terraform-provider-azapi/internal/acceptance"
	"github.com/Azure/terraform-provider-azapi/internal/acceptance/check"
	"github.com/Azure/terraform-provider-azapi/internal/clients"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

type ResourceWithAcquirePolicyToken struct{}

func TestAccResourceWithAcquirePolicyToken_allowedByPolicy(t *testing.T) {
	data := acceptance.BuildTestData(t, "azapi_resource", "storage")
	r := ResourceWithAcquirePolicyToken{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basic(data, true),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
	})
}

func TestAccResourceWithAcquirePolicyToken_disallowedByPolicy(t *testing.T) {
	data := acceptance.BuildTestData(t, "azapi_resource", "storage")
	r := ResourceWithAcquirePolicyToken{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config:      r.basic(data, false),
			ExpectError: regexp.MustCompile("RequestDisallowedByPolicy"),
		},
	})
}

func TestAccResourceWithAcquirePolicyToken_alwaysAcquire(t *testing.T) {
	data := acceptance.BuildTestData(t, "azapi_resource", "storage")
	r := ResourceWithAcquirePolicyToken{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.alwaysAcquire(data),
			Check: resource.ComposeTestCheckFunc(
				check.That("azapi_resource.storage").ExistsInAzure(r),
				check.That("azapi_resource.virtual_network").ExistsInAzure(r),
			),
		},
	})
}

func (r ResourceWithAcquirePolicyToken) Exists(ctx context.Context, client *clients.Client, state *terraform.InstanceState) (*bool, error) {
	out := true
	return &out, nil
}

func (r ResourceWithAcquirePolicyToken) basic(data acceptance.TestData, policyAlwaysPasses bool) string {
	// The CoinFlip demo endpoint always returns "testString" for the claims().string value.
	// Requiring that value makes the policy evaluation always pass; requiring any other value
	// makes it always fail.
	claimValue := "testString"
	if !policyAlwaysPasses {
		claimValue = "notTestString"
	}

	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

provider "azapi" {
}

data "azurerm_client_config" "current" {}

resource "azapi_resource" "policy_definition" {
  type      = "Microsoft.Authorization/policyDefinitions@2025-03-01"
  parent_id = "/subscriptions/${data.azurerm_client_config.current.subscription_id}"
  name      = "acctest-policy-%[1]d"

  body = {
    properties = {
      displayName = "acctest-policy-%[1]d"
      mode        = "All"
      policyType  = "Custom"
      versions    = ["1.0.0"]
      policyRule = {
        if = {
          allOf = [
            {
              field = "type"
              in    = ["Microsoft.Storage/storageAccounts"]
            },
            {
              value     = "[claims().string]"
              notEquals = "%[4]s"
            }
          ]
        }
        then = {
          effect = "deny"
        }
      }
      externalEvaluationEnforcementSettings = {
        roleDefinitionIds = ["/providers/Microsoft.Authorization/roleDefinitions/b24988ac-6180-42a0-ab88-20f7382dd24c"]
        endpointSettings = {
          kind = "CoinFlip"
          details = {
            successProbability = 1
          }
        }
      }
    }
  }

  schema_validation_enabled = false
}

resource "azurerm_resource_group" "test" {
  name     = "acctest-rg-%[1]d"
  location = "%[2]s"
}

resource "azurerm_resource_group_policy_assignment" "test" {
  name                 = "acctest-assignment-%[1]d"
  resource_group_id    = azurerm_resource_group.test.id
  policy_definition_id = azapi_resource.policy_definition.id
  location             = azurerm_resource_group.test.location

  identity {
    type = "SystemAssigned"
  }
}

resource "azapi_resource" "storage" {
  type      = "Microsoft.Storage/storageAccounts@2023-05-01"
  name      = "acctestsa%[3]s"
  parent_id = azurerm_resource_group.test.id
  location  = azurerm_resource_group.test.location

  body = {
    kind = "StorageV2"
    sku = {
      name = "Standard_LRS"
    }
    properties = {
      accessTier                   = "Hot"
      allowBlobPublicAccess        = false
      allowCrossTenantReplication  = true
      allowSharedKeyAccess         = true
      defaultToOAuthAuthentication = false
      isHnsEnabled                 = false
      isNfsV3Enabled               = false
      isSftpEnabled                = false
      minimumTlsVersion            = "TLS1_2"
      networkAcls = {
        defaultAction = "Allow"
      }
      publicNetworkAccess      = "Enabled"
      supportsHttpsTrafficOnly = true
    }
  }

  depends_on = [azurerm_resource_group_policy_assignment.test]
}
`, data.RandomInteger, data.LocationPrimary, data.RandomString, claimValue)
}

func (r ResourceWithAcquirePolicyToken) alwaysAcquire(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

provider "azapi" {
  always_acquire_policy_token = true
}

data "azurerm_client_config" "current" {}

resource "azapi_resource" "policy_definition" {
  type      = "Microsoft.Authorization/policyDefinitions@2025-03-01"
  parent_id = "/subscriptions/${data.azurerm_client_config.current.subscription_id}"
  name      = "acctest-policy-%[1]d"

  body = {
    properties = {
      displayName = "acctest-policy-%[1]d"
      mode        = "All"
      policyType  = "Custom"
      versions    = ["1.0.0"]
      policyRule = {
        if = {
          allOf = [
            {
              field = "type"
              in    = ["Microsoft.Storage/storageAccounts"]
            },
            {
              value     = "[claims().string]"
              notEquals = "testString"
            }
          ]
        }
        then = {
          effect = "deny"
        }
      }
      externalEvaluationEnforcementSettings = {
        roleDefinitionIds = ["/providers/Microsoft.Authorization/roleDefinitions/b24988ac-6180-42a0-ab88-20f7382dd24c"]
        endpointSettings = {
          kind = "CoinFlip"
          details = {
            successProbability = 1
          }
        }
      }
    }
  }

  schema_validation_enabled = false
}

resource "azurerm_resource_group" "test" {
  name     = "acctest-rg-%[1]d"
  location = "%[2]s"
}

resource "azurerm_resource_group_policy_assignment" "test" {
  name                 = "acctest-assignment-%[1]d"
  resource_group_id    = azurerm_resource_group.test.id
  policy_definition_id = azapi_resource.policy_definition.id
  location             = azurerm_resource_group.test.location

  identity {
    type = "SystemAssigned"
  }
}

resource "azapi_resource" "storage" {
  type      = "Microsoft.Storage/storageAccounts@2023-05-01"
  name      = "acctestsa%[3]s"
  parent_id = azurerm_resource_group.test.id
  location  = azurerm_resource_group.test.location

  body = {
    kind = "StorageV2"
    sku = {
      name = "Standard_LRS"
    }
    properties = {
      accessTier                   = "Hot"
      allowBlobPublicAccess        = false
      allowCrossTenantReplication  = true
      allowSharedKeyAccess         = true
      defaultToOAuthAuthentication = false
      isHnsEnabled                 = false
      isNfsV3Enabled               = false
      isSftpEnabled                = false
      minimumTlsVersion            = "TLS1_2"
      networkAcls = {
        defaultAction = "Allow"
      }
      publicNetworkAccess      = "Enabled"
      supportsHttpsTrafficOnly = true
    }
  }

  depends_on = [azurerm_resource_group_policy_assignment.test]
}

resource "azapi_resource" "virtual_network" {
  type      = "Microsoft.Network/virtualNetworks@2023-11-01"
  name      = "acctest-vnet-%[1]d"
  parent_id = azurerm_resource_group.test.id
  location  = azurerm_resource_group.test.location

  body = {
    properties = {
      addressSpace = {
        addressPrefixes = ["10.0.0.0/16"]
      }
    }
  }

  depends_on = [azurerm_resource_group_policy_assignment.test]
}
`, data.RandomInteger, data.LocationPrimary, data.RandomString)
}
