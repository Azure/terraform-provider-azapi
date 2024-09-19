package services_test

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/Azure/terraform-provider-azapi/internal/acceptance"
	"github.com/Azure/terraform-provider-azapi/internal/acceptance/check"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

const PreviousVersion = "1.14.0"

func TestAccAzapiResourceUpgrade_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azapi_resource", "test")
	r := GenericResource{}

	data.UpgradeTest(t, r, []resource.TestStep{
		data.UpgradeTestDeployStep(resource.TestStep{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		}, PreviousVersion),
		data.UpgradeTestPlanStep(resource.TestStep{
			Config: r.basic(data),
		}),
	})
}

func TestAccAzapiResourceUpgrade_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azapi_resource", "test")
	r := GenericResource{}

	data.UpgradeTest(t, r, []resource.TestStep{
		data.UpgradeTestDeployStep(resource.TestStep{
			Config: r.complete(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		}, PreviousVersion),
		data.UpgradeTestPlanStep(resource.TestStep{
			Config: r.complete(data),
		}),
	})
}

func TestAccAzapiResourceUpgrade_identityNone(t *testing.T) {
	data := acceptance.BuildTestData(t, "azapi_resource", "test")
	r := GenericResource{}

	data.UpgradeTest(t, r, []resource.TestStep{
		data.UpgradeTestDeployStep(resource.TestStep{
			Config: r.identityNone(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		}, PreviousVersion),
		data.UpgradeTestPlanStep(resource.TestStep{
			Config: r.identityNone(data),
		}),
	})
}

func TestAccAzapiResourceUpgrade_identitySystemAssigned(t *testing.T) {
	data := acceptance.BuildTestData(t, "azapi_resource", "test")
	r := GenericResource{}

	data.UpgradeTest(t, r, []resource.TestStep{
		data.UpgradeTestDeployStep(resource.TestStep{
			Config: r.identitySystemAssigned(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		}, PreviousVersion),
		data.UpgradeTestPlanStep(resource.TestStep{
			Config: r.identitySystemAssigned(data),
		}),
	})
}

func TestAccAzapiResourceUpgrade_identityUserAssigned(t *testing.T) {
	data := acceptance.BuildTestData(t, "azapi_resource", "test")
	r := GenericResource{}

	data.UpgradeTest(t, r, []resource.TestStep{
		data.UpgradeTestDeployStep(resource.TestStep{
			Config: r.identityUserAssigned(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		}, PreviousVersion),
		data.UpgradeTestPlanStep(resource.TestStep{
			Config: r.identityUserAssigned(data),
		}),
	})
}

func TestAccAzapiResourceUpgrade_completeBody(t *testing.T) {
	data := acceptance.BuildTestData(t, "azapi_resource", "test")
	r := GenericResource{}

	data.UpgradeTest(t, r, []resource.TestStep{
		data.UpgradeTestDeployStep(resource.TestStep{
			Config: r.completeJsonBody(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		}, PreviousVersion),
		data.UpgradeTestPlanStep(resource.TestStep{
			Config: r.completeBody(data),
		}),
	})
}

func TestAccAzapiResourceUpgrade_defaultTag(t *testing.T) {
	data := acceptance.BuildTestData(t, "azapi_resource", "test")
	r := GenericResource{}

	data.UpgradeTest(t, r, []resource.TestStep{
		data.UpgradeTestDeployStep(resource.TestStep{
			Config: r.defaultTag(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		}, PreviousVersion),
		data.UpgradeTestPlanStep(resource.TestStep{
			Config: r.defaultTag(data),
		}),
	})
}

func TestAccAzapiResourceUpgrade_defaultTagOverrideInBody(t *testing.T) {
	data := acceptance.BuildTestData(t, "azapi_resource", "test")
	r := GenericResource{}

	data.UpgradeTest(t, r, []resource.TestStep{
		data.UpgradeTestDeployStep(resource.TestStep{
			Config: r.defaultTagOverrideInBody(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		}, PreviousVersion),
		data.UpgradeTestPlanStep(resource.TestStep{
			Config: r.defaultTagOverrideInBody(data),
		}),
	})
}

func TestAccAzapiResourceUpgrade_defaultTagOverrideInHcl(t *testing.T) {
	data := acceptance.BuildTestData(t, "azapi_resource", "test")
	r := GenericResource{}

	data.UpgradeTest(t, r, []resource.TestStep{
		data.UpgradeTestDeployStep(resource.TestStep{
			Config: r.defaultTagOverrideInHcl(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		}, PreviousVersion),
		data.UpgradeTestPlanStep(resource.TestStep{
			Config: r.defaultTagOverrideInHcl(data),
		}),
	})
}

func TestAccAzapiResourceUpgrade_defaultLocation(t *testing.T) {
	data := acceptance.BuildTestData(t, "azapi_resource", "test")
	r := GenericResource{}

	data.UpgradeTest(t, r, []resource.TestStep{
		data.UpgradeTestDeployStep(resource.TestStep{
			Config: r.defaultLocation(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		}, PreviousVersion),
		data.UpgradeTestPlanStep(resource.TestStep{
			Config: r.defaultLocation(data),
		}),
	})
}

func TestAccAzapiResourceUpgrade_defaultLocationOverrideInHcl(t *testing.T) {
	data := acceptance.BuildTestData(t, "azapi_resource", "test")
	r := GenericResource{}

	data.UpgradeTest(t, r, []resource.TestStep{
		data.UpgradeTestDeployStep(resource.TestStep{
			Config: r.defaultLocationOverrideInHcl(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		}, PreviousVersion),
		data.UpgradeTestPlanStep(resource.TestStep{
			Config: r.defaultLocationOverrideInHcl(data),
		}),
	})
}

func TestAccAzapiResourceUpgrade_defaultParentId(t *testing.T) {
	data := acceptance.BuildTestData(t, "azapi_resource", "test")
	r := GenericResource{}

	updatedConfig := fmt.Sprintf(`
provider "azapi" {
}

resource "azapi_resource" "test" {
  type     = "Microsoft.Resources/resourceGroups@2023-07-01"
  name     = "acctest-%[2]d"
  location = "%[1]s"

  // add the below config to make the config match the migrated state
  body = {}
}
`, data.LocationPrimary, data.RandomInteger)

	data.UpgradeTest(t, r, []resource.TestStep{
		data.UpgradeTestDeployStep(resource.TestStep{
			Config: r.defaultParentId(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		}, PreviousVersion),
		data.UpgradeTestPlanStep(resource.TestStep{
			Config: updatedConfig,
		}),
	})
}

func TestAccAzapiResourceUpgrade_defaultNaming(t *testing.T) {
	data := acceptance.BuildTestData(t, "azapi_resource", "test")
	r := GenericResource{}

	data.UpgradeTest(t, r, []resource.TestStep{
		data.UpgradeTestDeployStep(resource.TestStep{
			Config: r.defaultNaming(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		}, PreviousVersion),
		data.UpgradeTestPlanStep(resource.TestStep{
			Config: r.defaultNaming(data),
		}),
	})
}

func TestAccAzapiResourceUpgrade_defaultNamingOverrideInHcl(t *testing.T) {
	data := acceptance.BuildTestData(t, "azapi_resource", "test")
	r := GenericResource{}

	data.UpgradeTest(t, r, []resource.TestStep{
		data.UpgradeTestDeployStep(resource.TestStep{
			Config: r.defaultNamingOverrideInHcl(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		}, PreviousVersion),
		data.UpgradeTestPlanStep(resource.TestStep{
			Config: r.defaultNamingOverrideInHcl(data),
		}),
	})
}

func TestAccAzapiResourceUpgrade_defaultsNotApplicable(t *testing.T) {
	data := acceptance.BuildTestData(t, "azapi_resource", "test")
	r := GenericResource{}

	data.UpgradeTest(t, r, []resource.TestStep{
		data.UpgradeTestDeployStep(resource.TestStep{
			Config: r.defaultsNotApplicable(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		}, PreviousVersion),
		data.UpgradeTestPlanStep(resource.TestStep{
			Config: r.defaultsNotApplicable(data),
		}),
	})
}

func TestAccAzapiResourceUpgrade_subscriptionScope(t *testing.T) {
	data := acceptance.BuildTestData(t, "azapi_resource", "test")
	r := GenericResource{}

	subscriptionId := os.Getenv("ARM_SUBSCRIPTION_ID")

	updatedConfig := fmt.Sprintf(`
resource "azapi_resource" "test" {
  type      = "Microsoft.Resources/resourceGroups@2023-07-01"
  name      = "acctestRG-%[1]d"
  parent_id = "/subscriptions/%[2]s"

  location = "%[3]s"

  // add the below config to make the config match the migrated state
  body = {}
}
`, data.RandomInteger, subscriptionId, data.LocationPrimary)

	data.UpgradeTest(t, r, []resource.TestStep{
		data.UpgradeTestDeployStep(resource.TestStep{
			Config: r.subscriptionScope(data, subscriptionId),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		}, PreviousVersion),
		data.UpgradeTestPlanStep(resource.TestStep{
			Config: updatedConfig,
		}),
	})
}

func TestAccAzapiResourceUpgrade_extensionScope(t *testing.T) {
	data := acceptance.BuildTestData(t, "azapi_resource", "test")
	r := GenericResource{}

	data.UpgradeTest(t, r, []resource.TestStep{
		data.UpgradeTestDeployStep(resource.TestStep{
			Config: r.extensionScope(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		}, PreviousVersion),
		data.UpgradeTestPlanStep(resource.TestStep{
			Config: r.extensionScope(data),
		}),
	})
}

func TestAccAzapiResourceUpgrade_ignoreMissingProperty(t *testing.T) {
	data := acceptance.BuildTestData(t, "azapi_resource", "test")
	r := GenericResource{}

	data.UpgradeTest(t, r, []resource.TestStep{
		data.UpgradeTestDeployStep(resource.TestStep{
			Config: r.ignoreMissingProperty(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		}, PreviousVersion),
		data.UpgradeTestPlanStep(resource.TestStep{
			Config: r.ignoreMissingProperty(data),
		}),
	})
}

func TestAccAzapiResourceUpgrade_ignoreCasing(t *testing.T) {
	data := acceptance.BuildTestData(t, "azapi_resource", "test")
	r := GenericResource{}

	data.UpgradeTest(t, r, []resource.TestStep{
		data.UpgradeTestDeployStep(resource.TestStep{
			Config: r.ignoreCasing(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		}, PreviousVersion),
		data.UpgradeTestPlanStep(resource.TestStep{
			Config: r.ignoreCasing(data),
		}),
	})
}

func TestAccAzapiResourceUpgrade_deleteLROEndsWithNotFoundError(t *testing.T) {
	data := acceptance.BuildTestData(t, "azapi_resource", "test")
	r := GenericResource{}

	data.UpgradeTest(t, r, []resource.TestStep{
		data.UpgradeTestDeployStep(resource.TestStep{
			Config: r.deleteLROEndsWithNotFoundError(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		}, PreviousVersion),
		data.UpgradeTestPlanStep(resource.TestStep{
			Config: r.deleteLROEndsWithNotFoundError(data),
		}),
	})
}

func TestAccAzapiResourceUpgrade_locks(t *testing.T) {
	data := acceptance.BuildTestData(t, "azapi_resource", "test")
	r := GenericResource{}

	data.UpgradeTest(t, r, []resource.TestStep{
		data.UpgradeTestDeployStep(resource.TestStep{
			Config: r.locks(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		}, PreviousVersion),
		data.UpgradeTestPlanStep(resource.TestStep{
			Config: r.locks(data),
		}),
	})
}

func TestAccAzapiResourceUpgrade_secretsInAsterisks(t *testing.T) {
	data := acceptance.BuildTestData(t, "azapi_resource", "test")
	r := GenericResource{}

	clientId := os.Getenv("ARM_CLIENT_ID")
	clientSecret := os.Getenv("ARM_CLIENT_SECRET")

	data.UpgradeTest(t, r, []resource.TestStep{
		data.UpgradeTestDeployStep(resource.TestStep{
			Config: r.secretsInAsterisks(data, clientId, clientSecret),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		}, PreviousVersion),
		data.UpgradeTestPlanStep(resource.TestStep{
			Config: r.secretsInAsterisks(data, clientId, clientSecret),
		}),
	})
}

func TestAccAzapiResourceUpgrade_nullLocation(t *testing.T) {
	data := acceptance.BuildTestData(t, "azapi_resource", "test")
	r := GenericResource{}

	data.UpgradeTest(t, r, []resource.TestStep{
		data.UpgradeTestDeployStep(resource.TestStep{
			Config:            r.nullLocation(data),
			ExternalProviders: externalProvidersAzurerm(),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		}, PreviousVersion),
		data.UpgradeTestPlanStep(resource.TestStep{
			Config:            r.nullLocation(data),
			ExternalProviders: externalProvidersAzurerm(),
		}),
	})
}

func TestAccAzapiResourceUpgrade_timeouts(t *testing.T) {
	data := acceptance.BuildTestData(t, "azapi_resource", "test")
	r := GenericResource{}

	data.UpgradeTest(t, r, []resource.TestStep{
		data.UpgradeTestDeployStep(resource.TestStep{
			Config: r.timeouts(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		}, PreviousVersion),
		data.UpgradeTestPlanStep(resource.TestStep{
			Config: r.timeouts(data),
		}),
	})
}

func TestAccAzapiResourceUpgrade_timeouts_from_v1_13_1(t *testing.T) {
	data := acceptance.BuildTestData(t, "azapi_resource", "test")
	r := GenericResource{}

	data.UpgradeTest(t, r, []resource.TestStep{
		data.UpgradeTestDeployStep(resource.TestStep{
			Config: strings.ReplaceAll(r.timeouts(data), `update = "10m"`, ""),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		}, "1.13.1"),
		data.UpgradeTestApplyStep(resource.TestStep{
			Config: r.timeouts(data),
		}),
		data.UpgradeTestPlanStep(resource.TestStep{
			Config: r.timeouts(data),
		}),
	})
}

func TestAccAzapiResourceUpgrade_completeBody_from_schema_v0(t *testing.T) {
	data := acceptance.BuildTestData(t, "azapi_resource", "test")
	r := GenericResource{}

	data.UpgradeTest(t, r, []resource.TestStep{
		data.UpgradeTestDeployStep(resource.TestStep{
			Config: r.completeJsonBody(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		}, "1.12.1"),
		data.UpgradeTestPlanStep(resource.TestStep{
			Config: r.completeBody(data),
			// It's a known breaking change that identity in the body will not be synced to the top-level identity block
			ExpectNonEmptyPlan: true,
		}),
		data.UpgradeTestApplyStep(resource.TestStep{
			Config: r.completeBody(data),
		}),
		data.UpgradeTestPlanStep(resource.TestStep{
			Config: r.completeBody(data),
		}),
	})
}

func TestAccAzapiResourceUpgrade_basic_from_schema_v0(t *testing.T) {
	data := acceptance.BuildTestData(t, "azapi_resource", "test")
	r := GenericResource{}

	updatedConfig := r.oldConfig(data)
	updatedConfig = strings.ReplaceAll(updatedConfig, "jsonencode({", "{")
	updatedConfig = strings.ReplaceAll(updatedConfig, "})", "}")

	data.UpgradeTest(t, r, []resource.TestStep{
		data.UpgradeTestDeployStep(resource.TestStep{
			Config: r.oldConfig(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		}, "1.12.1"),
		data.UpgradeTestPlanStep(resource.TestStep{
			Config: updatedConfig,
		}),
	})
}
