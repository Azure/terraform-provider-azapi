package services_test

import (
	"strings"
	"testing"

	"github.com/Azure/terraform-provider-azapi/internal/acceptance"
	"github.com/Azure/terraform-provider-azapi/internal/acceptance/check"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAzapiDataPlaneResourceUpgrade_appConfigKeyValues(t *testing.T) {
	data := acceptance.BuildTestData(t, "azapi_data_plane_resource", "test")
	r := DataPlaneResource{}

	data.UpgradeTest(t, r, []resource.TestStep{
		data.UpgradeTestDeployStep(resource.TestStep{
			Config: r.appConfigKeyValues(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		}, PreviousVersion),
		data.UpgradeTestPlanStep(resource.TestStep{
			Config: r.appConfigKeyValues(data),
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
			Config: r.keyVaultIssuer(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		}, PreviousVersion),
		data.UpgradeTestPlanStep(resource.TestStep{
			Config: r.keyVaultIssuer(data),
		}),
	})
}

func TestAccAzapiDataPlaneResourceUpgrade_iotAppsUser(t *testing.T) {
	data := acceptance.BuildTestData(t, "azapi_data_plane_resource", "test")
	r := DataPlaneResource{}

	data.UpgradeTest(t, r, []resource.TestStep{
		data.UpgradeTestDeployStep(resource.TestStep{
			Config: r.iotAppsUser(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		}, PreviousVersion),
		data.UpgradeTestPlanStep(resource.TestStep{
			Config: r.iotAppsUser(data),
		}),
	})
}

func TestAccAzapiDataPlaneResourceUpgrade_timeouts(t *testing.T) {
	data := acceptance.BuildTestData(t, "azapi_data_plane_resource", "test")
	r := DataPlaneResource{}

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

func TestAccAzapiDataPlaneResourceUpgrade_timeouts_from_v1_13_1(t *testing.T) {
	data := acceptance.BuildTestData(t, "azapi_data_plane_resource", "test")
	r := DataPlaneResource{}

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
