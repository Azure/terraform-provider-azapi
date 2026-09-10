package services_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/Azure/terraform-provider-azapi/internal/acceptance"
	"github.com/Azure/terraform-provider-azapi/internal/clients"
	"github.com/Azure/terraform-provider-azapi/internal/services/common"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccGenericResource_movePortalDashboardFromAzureRMUsesConfiguredApiVersion(t *testing.T) {
	acceptance.SkipIfCoreAcctestsOnly(t, "Edge case: latest spec version runtime not yet functional")

	data := acceptance.BuildTestData(t, "azapi_resource", "test")
	r := GenericResource{}

	defer func() {
		client, err := acceptance.BuildTestClient()
		if err != nil {
			t.Errorf("building test client for cleanup: %+v", err)
			return
		}

		resourceGroupID := fmt.Sprintf("/subscriptions/%s/resourceGroups/acctestRG-%d", client.Account.GetSubscriptionId(), data.RandomInteger)
		_, _ = client.ResourceClient.Delete(context.Background(), resourceGroupID, "2023-07-01", clients.DefaultRequestOptions())
	}()

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config:            r.portalDashboardAzureRM(data),
			ExternalProviders: common.ExternalProvidersAzurermVersionFour(),
		},
		{
			Config:             r.portalDashboardAzureRMMoved(data),
			ExternalProviders:  common.ExternalProvidersAzurermVersionFour(),
			ExpectNonEmptyPlan: false,
		},
	})
}

func (r GenericResource) portalDashboardAzureRM(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

provider "azurerm" {
  features {
    resource_group {
      prevent_deletion_if_contains_resources = false
    }
  }
}

resource "azurerm_portal_dashboard" "test" {
  name                 = "acctest%[2]s"
  resource_group_name  = azapi_resource.resourceGroup.name
  location             = azapi_resource.resourceGroup.location
  dashboard_properties = jsonencode({ lenses = {} })

  tags = {
    source = "acctest"
  }
}
`, r.template(data), data.RandomString)
}

func (r GenericResource) portalDashboardAzureRMMoved(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

provider "azurerm" {
  features {
    resource_group {
      prevent_deletion_if_contains_resources = false
    }
  }
}

moved {
  from = azurerm_portal_dashboard.test
  to   = azapi_resource.test
}

resource "azapi_resource" "test" {
  type      = "Microsoft.Portal/dashboards@2019-01-01-preview"
  name      = "acctest%[2]s"
  parent_id = azapi_resource.resourceGroup.id
  location  = azapi_resource.resourceGroup.location

  tags = {
    source = "acctest"
  }

  body = {
    properties = {
      lenses = {}
    }
  }
}
`, r.template(data), data.RandomString)
}
