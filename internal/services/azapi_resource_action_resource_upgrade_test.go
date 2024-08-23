package services_test

import (
	"os"
	"strings"
	"testing"

	"github.com/Azure/terraform-provider-azapi/internal/acceptance"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAzapiActionResourceUpgrade_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azapi_resource_action", "test")
	r := ActionResource{}

	data.UpgradeTest(t, r, []resource.TestStep{
		data.UpgradeTestDeployStep(resource.TestStep{
			Config: r.basic(data),
			Check:  resource.ComposeTestCheckFunc(),
		}, PreviousVersion),
		data.UpgradeTestPlanStep(resource.TestStep{
			Config: r.basic(data),
		}),
	})
}

func TestAccAzapiActionResourceUpgrade_basicWhenDestroy(t *testing.T) {
	data := acceptance.BuildTestData(t, "azapi_resource_action", "test")
	r := ActionResource{}

	data.UpgradeTest(t, r, []resource.TestStep{
		data.UpgradeTestDeployStep(resource.TestStep{
			Config: r.basicWhenDestroy(data),
			Check:  resource.ComposeTestCheckFunc(),
		}, PreviousVersion),
		data.UpgradeTestPlanStep(resource.TestStep{
			Config: r.basicWhenDestroy(data),
		}),
	})
}

func TestAccAzapiActionResourceUpgrade_registerResourceProvider(t *testing.T) {
	data := acceptance.BuildTestData(t, "azapi_resource_action", "test")
	r := ActionResource{}

	subscriptionId := os.Getenv("ARM_SUBSCRIPTION_ID")
	data.UpgradeTest(t, r, []resource.TestStep{
		data.UpgradeTestDeployStep(resource.TestStep{
			Config: r.registerResourceProvider(subscriptionId),
			Check:  resource.ComposeTestCheckFunc(),
		}, PreviousVersion),
		data.UpgradeTestPlanStep(resource.TestStep{
			Config: r.registerResourceProvider(subscriptionId),
		}),
	})
}

func TestAccAzapiActionResourceUpgrade_upgradeFromVeryOldVersion(t *testing.T) {
	data := acceptance.BuildTestData(t, "azapi_resource_action", "test")
	r := ActionResource{}

	subscriptionId := os.Getenv("ARM_SUBSCRIPTION_ID")
	data.UpgradeTest(t, r, []resource.TestStep{
		data.UpgradeTestDeployStep(resource.TestStep{
			Config: r.registerResourceProvider(subscriptionId),
			Check:  resource.ComposeTestCheckFunc(),
		}, "1.8.0"),
		data.UpgradeTestPlanStep(resource.TestStep{
			Config: r.registerResourceProvider(subscriptionId),
		}),
	})
}

func TestAccAzapiActionResourceUpgrade_providerAction(t *testing.T) {
	data := acceptance.BuildTestData(t, "azapi_resource_action", "test")
	r := ActionResource{}

	data.UpgradeTest(t, r, []resource.TestStep{
		data.UpgradeTestDeployStep(resource.TestStep{
			Config: r.providerAction(data),
			Check:  resource.ComposeTestCheckFunc(),
		}, PreviousVersion),
		data.UpgradeTestPlanStep(resource.TestStep{
			Config: r.providerAction(data),
		}),
	})
}

func TestAccAzapiActionResourceUpgrade_nonstandardLRO(t *testing.T) {
	data := acceptance.BuildTestData(t, "azapi_resource_action", "test")
	r := ActionResource{}

	data.UpgradeTest(t, r, []resource.TestStep{
		data.UpgradeTestDeployStep(resource.TestStep{
			Config:            r.nonstandardLRO(data),
			ExternalProviders: externalProvidersAzurerm(),
			Check:             resource.ComposeTestCheckFunc(),
		}, PreviousVersion),
		data.UpgradeTestPlanStep(resource.TestStep{
			ExternalProviders: externalProvidersAzurerm(),
			Config:            r.nonstandardLRO(data),
		}),
	})
}

func TestAccAzapiActionResourceUpgrade_timeouts(t *testing.T) {
	data := acceptance.BuildTestData(t, "azapi_resource_action", "test")
	r := ActionResource{}

	data.UpgradeTest(t, r, []resource.TestStep{
		data.UpgradeTestDeployStep(resource.TestStep{
			Config: r.timeouts(data),
			Check:  resource.ComposeTestCheckFunc(),
		}, PreviousVersion),
		data.UpgradeTestPlanStep(resource.TestStep{
			Config: r.timeouts(data),
		}),
	})
}

func TestAccAzapiActionResourceUpgrade_timeouts_from_v1_13_1(t *testing.T) {
	data := acceptance.BuildTestData(t, "azapi_resource_action", "test")
	r := ActionResource{}

	data.UpgradeTest(t, r, []resource.TestStep{
		data.UpgradeTestDeployStep(resource.TestStep{
			Config: strings.ReplaceAll(r.timeouts(data), `update = "10m"`, ""),
			Check:  resource.ComposeTestCheckFunc(),
		}, "1.13.1"),
		data.UpgradeTestApplyStep(resource.TestStep{
			Config: r.timeouts(data),
		}),
		data.UpgradeTestPlanStep(resource.TestStep{
			Config: r.timeouts(data),
		}),
	})
}

func TestAccAzapiActionResourceUpgrade_basic_from_schema_v0(t *testing.T) {
	data := acceptance.BuildTestData(t, "azapi_resource_action", "test")
	r := ActionResource{}

	subscriptionId := os.Getenv("ARM_SUBSCRIPTION_ID")
	updatedConfig := r.oldConfig(data, subscriptionId)
	updatedConfig = strings.ReplaceAll(updatedConfig, "jsonencode({", "{")
	updatedConfig = strings.ReplaceAll(updatedConfig, "})", "}")

	data.UpgradeTest(t, r, []resource.TestStep{
		data.UpgradeTestDeployStep(resource.TestStep{
			Config: r.oldConfig(data, subscriptionId),
			Check:  resource.ComposeTestCheckFunc(),
		}, "1.12.1"),
		data.UpgradeTestPlanStep(resource.TestStep{
			Config: updatedConfig,
		}),
	})
}
