package services_test

import (
	"testing"

	"github.com/Azure/terraform-provider-azapi/internal/acceptance"
	"github.com/Azure/terraform-provider-azapi/internal/acceptance/check"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAzapiUpdateResourceUpgrade_automationAccount(t *testing.T) {
	data := acceptance.BuildTestData(t, "azapi_update_resource", "test")
	r := GenericUpdateResource{}

	data.UpgradeTest(t, r, []resource.TestStep{
		data.UpgradeTestDeployStep(resource.TestStep{
			Config: r.automationAccount(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		}, PreviousVersion),
		data.UpgradeTestPlanStep(resource.TestStep{
			Config: r.automationAccount(data),
		}),
	})
}

func TestAccAzapiUpdateResourceUpgrade_automationAccountWithNameParentId(t *testing.T) {
	data := acceptance.BuildTestData(t, "azapi_update_resource", "test")
	r := GenericUpdateResource{}

	data.UpgradeTest(t, r, []resource.TestStep{
		data.UpgradeTestDeployStep(resource.TestStep{
			Config: r.automationAccountWithNameParentId(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		}, PreviousVersion),
		data.UpgradeTestPlanStep(resource.TestStep{
			Config: r.automationAccountWithNameParentId(data),
		}),
	})
}

func TestAccAzapiUpdateResourceUpgrade_siteConfigSlotConfigNames(t *testing.T) {
	data := acceptance.BuildTestData(t, "azapi_update_resource", "test")
	r := GenericUpdateResource{}

	data.UpgradeTest(t, r, []resource.TestStep{
		data.UpgradeTestDeployStep(resource.TestStep{
			Config: r.siteConfigSlotConfigNames(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		}, PreviousVersion),
		data.UpgradeTestPlanStep(resource.TestStep{
			Config: r.siteConfigSlotConfigNames(data),
		}),
	})
}

func TestAccAzapiUpdateResourceUpgrade_locks(t *testing.T) {
	data := acceptance.BuildTestData(t, "azapi_update_resource", "test")
	r := GenericUpdateResource{}

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

func TestAccAzapiUpdateResourceUpgrade_ignoreChanges(t *testing.T) {
	data := acceptance.BuildTestData(t, "azapi_update_resource", "test")
	r := GenericUpdateResource{}

	data.UpgradeTest(t, r, []resource.TestStep{
		data.UpgradeTestDeployStep(resource.TestStep{
			Config: r.ignoreChanges(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		}, PreviousVersion),
		data.UpgradeTestPlanStep(resource.TestStep{
			Config: r.ignoreChanges(data),
		}),
	})
}

func TestAccAzapiUpdateResourceUpgrade_ignoreChangesArray(t *testing.T) {
	data := acceptance.BuildTestData(t, "azapi_update_resource", "test")
	r := GenericUpdateResource{}

	data.UpgradeTest(t, r, []resource.TestStep{
		data.UpgradeTestDeployStep(resource.TestStep{
			Config: r.ignoreChangesArray(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		}, PreviousVersion),
		data.UpgradeTestPlanStep(resource.TestStep{
			Config: r.ignoreChangesArray(data),
		}),
	})
}
