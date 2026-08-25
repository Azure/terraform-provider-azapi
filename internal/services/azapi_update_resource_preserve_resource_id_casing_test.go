package services_test

import (
	"fmt"
	"testing"

	"github.com/Azure/terraform-provider-azapi/internal/acceptance"
	"github.com/Azure/terraform-provider-azapi/internal/acceptance/check"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// TestAccGenericUpdateResource_preserveResourceIDCasing exercises the
// preserve_resource_id_casing provider feature flag end-to-end through
// azapi_update_resource. The source resource is declared with a lower-case
// leading character (webPubSub) while the azapi_update_resource type uses the
// canonical casing (WebPubSub), so the flag governs how the id / resource_id
// attributes are stored. Applying with the flag enabled must succeed and leave
// no perpetual diff (the implicit post-apply plan is asserted empty by the
// testing framework).
func TestAccGenericUpdateResource_preserveResourceIDCasing(t *testing.T) {
	data := acceptance.BuildTestData(t, "azapi_update_resource", "test")
	r := GenericUpdateResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.preserveResourceIDCasing(data, true),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("resource_id").Exists(),
			),
		},
	})
}

// TestAccAzapiUpdateResourceUpgrade_preserveResourceIDCasing reproduces the
// state-migration scenario from issue #1120. The resource is first deployed with
// an older provider version (which stores the resource_id with whatever casing it
// computed), then the provider is upgraded to the local build with
// preserve_resource_id_casing enabled. The plan step asserts that no perpetual
// diff is produced after the upgrade: the flag keeps the previously stored casing
// rather than re-normalising it, which is exactly the behaviour the flag adds for
// azapi_update_resource.resource_id.
func TestAccAzapiUpdateResourceUpgrade_preserveResourceIDCasing(t *testing.T) {
	data := acceptance.BuildTestData(t, "azapi_update_resource", "test")
	r := GenericUpdateResource{}

	data.UpgradeTest(t, r, []resource.TestStep{
		data.UpgradeTestDeployStep(resource.TestStep{
			Config: r.preserveResourceIDCasing(data, false),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		}, PreviousVersion),
		data.UpgradeTestPlanStep(resource.TestStep{
			Config: r.preserveResourceIDCasing(data, true),
		}),
		data.UpgradeTestApplyStep(resource.TestStep{
			Config: r.preserveResourceIDCasing(data, true),
		}),
	})
}

// preserveResourceIDCasing builds a config where an azapi_update_resource targets
// a resource whose id casing differs from the canonical type casing. When enabled
// is true an explicit provider block turns on preserve_resource_id_casing; when
// false the config omits the block so it can be deployed by older provider
// versions that do not know the attribute.
func (r GenericUpdateResource) preserveResourceIDCasing(data acceptance.TestData, enabled bool) string {
	providerBlock := ""
	if enabled {
		providerBlock = `
provider "azapi" {
  preserve_resource_id_casing = true
}
`
	}

	return fmt.Sprintf(`
%[1]s
%[2]s

resource "azapi_update_resource" "test" {
  type        = "Microsoft.SignalRService/WebPubSub@2024-04-01-preview"
  resource_id = azapi_resource.webPubSub.id

  body = {
    properties = {
      networkACLs = {
        defaultAction = "Deny"
        publicNetwork = {
          allow = ["ClientConnection"]
        }
        ipRules = [{
          value  = "0.0.0.0/0"
          action = "Allow"
        }]
      }
    }
  }
}
`, providerBlock, r.templateForIDCasing(data))
}
