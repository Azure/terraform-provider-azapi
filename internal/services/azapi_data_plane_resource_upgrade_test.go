package services_test

import (
	"strings"
	"testing"

	"github.com/Azure/terraform-provider-azapi/internal/acceptance"
	"github.com/Azure/terraform-provider-azapi/internal/acceptance/check"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func externalProvidersAzurerm() map[string]resource.ExternalProvider {
	return map[string]resource.ExternalProvider{
		"azurerm": {
			VersionConstraint: "3.106.0",
			Source:            "hashicorp/azurerm",
		},
	}
}

func TestAccAzapiDataPlaneResourceUpgrade_appConfigKeyValues(t *testing.T) {
	data := acceptance.BuildTestData(t, "azapi_data_plane_resource", "test")
	r := DataPlaneResource{}

	data.UpgradeTest(t, r, []resource.TestStep{
		data.UpgradeTestDeployStep(resource.TestStep{
			Config:            r.appConfigKeyValues(data),
			ExternalProviders: externalProvidersAzurerm(),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		}, PreviousVersion),
		data.UpgradeTestPlanStep(resource.TestStep{
			Config:            r.appConfigKeyValues(data),
			ExternalProviders: externalProvidersAzurerm(),
		}),
	})
}

func TestAccAzapiDataPlaneResourceUpgrade_purviewClassification(t *testing.T) {
	data := acceptance.BuildTestData(t, "azapi_data_plane_resource", "test")
	r := DataPlaneResource{}

	data.UpgradeTest(t, r, []resource.TestStep{
		data.UpgradeTestDeployStep(resource.TestStep{
			Config: r.purviewClassification(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		}, PreviousVersion),
		data.UpgradeTestPlanStep(resource.TestStep{
			Config: r.purviewClassification(data),
		}),
	})
}

func TestAccAzapiDataPlaneResourceUpgrade_purviewCollection(t *testing.T) {
	data := acceptance.BuildTestData(t, "azapi_data_plane_resource", "test")
	r := DataPlaneResource{}

	data.UpgradeTest(t, r, []resource.TestStep{
		data.UpgradeTestDeployStep(resource.TestStep{
			Config: r.purviewCollection(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		}, PreviousVersion),
		data.UpgradeTestPlanStep(resource.TestStep{
			Config: r.purviewCollection(data),
		}),
	})
}

func TestAccAzapiDataPlaneResourceUpgrade_keyVaultIssuer(t *testing.T) {
	data := acceptance.BuildTestData(t, "azapi_data_plane_resource", "test")
	r := DataPlaneResource{}

	data.UpgradeTest(t, r, []resource.TestStep{
		data.UpgradeTestDeployStep(resource.TestStep{
			Config:            r.keyVaultIssuer(data),
			ExternalProviders: externalProvidersAzurerm(),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		}, PreviousVersion),
		data.UpgradeTestPlanStep(resource.TestStep{
			Config:            r.keyVaultIssuer(data),
			ExternalProviders: externalProvidersAzurerm(),
		}),
	})
}

func TestAccAzapiDataPlaneResourceUpgrade_iotAppsUser(t *testing.T) {
	data := acceptance.BuildTestData(t, "azapi_data_plane_resource", "test")
	r := DataPlaneResource{}

	data.UpgradeTest(t, r, []resource.TestStep{
		data.UpgradeTestDeployStep(resource.TestStep{
			Config:            r.iotAppsUser(data),
			ExternalProviders: externalProvidersAzurerm(),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		}, PreviousVersion),
		data.UpgradeTestPlanStep(resource.TestStep{
			Config:            r.iotAppsUser(data),
			ExternalProviders: externalProvidersAzurerm(),
		}),
	})
}

func TestAccAzapiDataPlaneResourceUpgrade_timeouts(t *testing.T) {
	data := acceptance.BuildTestData(t, "azapi_data_plane_resource", "test")
	r := DataPlaneResource{}

	data.UpgradeTest(t, r, []resource.TestStep{
		data.UpgradeTestDeployStep(resource.TestStep{
			Config:            r.timeouts(data),
			ExternalProviders: externalProvidersAzurerm(),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		}, PreviousVersion),
		data.UpgradeTestPlanStep(resource.TestStep{
			Config:            r.timeouts(data),
			ExternalProviders: externalProvidersAzurerm(),
		}),
	})
}

func TestAccAzapiDataPlaneResourceUpgrade_timeouts_from_v1_13_1(t *testing.T) {
	data := acceptance.BuildTestData(t, "azapi_data_plane_resource", "test")
	r := DataPlaneResource{}

	data.UpgradeTest(t, r, []resource.TestStep{
		data.UpgradeTestDeployStep(resource.TestStep{
			Config:            strings.ReplaceAll(r.timeouts(data), `update = "10m"`, ""),
			ExternalProviders: externalProvidersAzurerm(),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		}, "1.13.1"),
		data.UpgradeTestApplyStep(resource.TestStep{
			ExternalProviders: externalProvidersAzurerm(),
			Config:            r.timeouts(data),
		}),
		data.UpgradeTestPlanStep(resource.TestStep{
			ExternalProviders: externalProvidersAzurerm(),
			Config:            r.timeouts(data),
		}),
	})
}

func TestAccAzapiDataPlaneResourceUpgrade_basic_from_schema_v0(t *testing.T) {
	data := acceptance.BuildTestData(t, "azapi_data_plane_resource", "test")
	r := DataPlaneResource{}

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
