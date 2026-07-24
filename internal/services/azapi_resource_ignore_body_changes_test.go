package services_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/Azure/terraform-provider-azapi/internal/acceptance"
	"github.com/Azure/terraform-provider-azapi/internal/acceptance/check"
	"github.com/Azure/terraform-provider-azapi/internal/clients"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

// TestAccGenericResource_ignoreBodyChangesSubnetsDrift demonstrates the problem tracked by
// https://github.com/Azure/terraform-provider-azapi/issues/583 and
// https://github.com/Azure/terraform-provider-azapi/issues/1137.
//
// A virtual network is created with a single subnet. An out-of-band controller (simulating
// Azure Policy, an autoscaler, or a DINE remediation) then adds a second subnet directly via
// the Azure API. When Terraform re-plans the unchanged (single-subnet) config it wants to
// remove the externally-added subnet, producing a perpetual non-empty plan.
//
// There is currently no dynamic, variable-driven way to tell azapi to ignore that body change:
// `lifecycle.ignore_changes` only accepts static literal paths (see the sibling test) and the
// former `ignore_body_changes` argument was removed. This test asserts the current behaviour
// (drift is reconciled -> non-empty plan). It is the baseline that a re-introduced dynamic
// body-ignore feature is expected to flip to an empty plan.
func TestAccGenericResource_ignoreBodyChangesSubnetsDrift(t *testing.T) {
	data := acceptance.BuildTestData(t, "azapi_resource", "test")
	r := GenericResource{}

	var vnetID string

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.ignoreBodyChangesSubnets(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				extractResourceID(data.ResourceName, &vnetID),
			),
		},
		{
			// An external controller adds a second subnet out-of-band before we re-plan.
			PreConfig: func() {
				r.addSubnetOutOfBand(t, vnetID)
			},
			// Re-applying the same single-subnet config reconciles the drift and removes
			// subnet2. Today there is no way to ignore this dynamically, so the plan is
			// non-empty. This is the gap that #583 / #1137 ask us to close.
			Config:             r.ignoreBodyChangesSubnets(data),
			PlanOnly:           true,
			ExpectNonEmptyPlan: true,
		},
	})
}

func (r GenericResource) ignoreBodyChangesSubnets(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azapi_resource" "test" {
  type      = "Microsoft.Network/virtualNetworks@2024-01-01"
  parent_id = azapi_resource.resourceGroup.id
  name      = "acctest-vnet-%[2]d"
  location  = azapi_resource.resourceGroup.location
  body = {
    properties = {
      addressSpace = {
        addressPrefixes = ["10.0.0.0/16"]
      }
      subnets = [
        {
          name = "subnet1"
          properties = {
            addressPrefixes = ["10.0.1.0/24"]
          }
        },
      ]
    }
  }
}
`, r.template(data), data.RandomInteger)
}

// addSubnetOutOfBand adds a subnet to the virtual network directly via the Azure API, bypassing
// Terraform, to simulate an external controller (e.g. Azure Policy) mutating the resource body.
func (GenericResource) addSubnetOutOfBand(t *testing.T, resourceID string) {
	client, err := acceptance.BuildTestClient()
	if err != nil {
		t.Fatalf("building test client: %+v", err)
	}

	ctx := context.Background()
	const apiVersion = "2024-01-01"

	existing, err := client.ResourceClient.Get(ctx, resourceID, apiVersion, clients.DefaultRequestOptions())
	if err != nil {
		t.Fatalf("reading %s: %+v", resourceID, err)
	}

	body, ok := existing.(map[string]interface{})
	if !ok {
		t.Fatalf("unexpected response type %T for %s", existing, resourceID)
	}
	properties, ok := body["properties"].(map[string]interface{})
	if !ok {
		properties = map[string]interface{}{}
		body["properties"] = properties
	}
	subnets, _ := properties["subnets"].([]interface{})
	subnets = append(subnets, map[string]interface{}{
		"name": "subnet2",
		"properties": map[string]interface{}{
			"addressPrefixes": []interface{}{"10.0.2.0/24"},
		},
	})
	properties["subnets"] = subnets

	if _, err := client.ResourceClient.CreateOrUpdate(ctx, resourceID, apiVersion, body, clients.DefaultRequestOptions()); err != nil {
		t.Fatalf("adding subnet out-of-band to %s: %+v", resourceID, err)
	}
}

// extractResourceID captures the Azure resource ID of the named resource from Terraform state
// so a later test step can reference it (e.g. to perform an out-of-band change).
func extractResourceID(resourceName string, out *string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("resource not found in state: %s", resourceName)
		}
		*out = rs.Primary.ID
		return nil
	}
}
