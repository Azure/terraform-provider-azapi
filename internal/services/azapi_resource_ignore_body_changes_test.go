package services_test

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/Azure/terraform-provider-azapi/internal/acceptance"
	"github.com/Azure/terraform-provider-azapi/internal/acceptance/check"
	"github.com/Azure/terraform-provider-azapi/internal/clients"
	"github.com/Azure/terraform-provider-azapi/internal/services/parse"
	"github.com/Azure/terraform-provider-azapi/utils"
	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccGenericResource_ignoreBodyChangesTags(t *testing.T) {
	data := acceptance.BuildTestData(t, "azapi_resource", "test")
	r := GenericResource{}
	var resourceID string
	var resourceType string

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: ignoreBodyChangesTagsConfig(data, false),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				captureResourceState(data.ResourceName, &resourceID, &resourceType),
			),
		},
		{
			PreConfig: func() {
				if err := addTagOutsideTerraform(resourceID, resourceType); err != nil {
					t.Fatalf("adding tag outside Terraform: %+v", err)
				}
			},
			Config:             ignoreBodyChangesTagsConfig(data, false),
			PlanOnly:           true,
			ExpectNonEmptyPlan: true,
		},
		{
			Config: ignoreBodyChangesTagsConfig(data, true),
			Check: func(*terraform.State) error {
				return checkTagInAzure(resourceID, resourceType, "out-of-band", "true")
			},
		},
		{
			Config:   ignoreBodyChangesTagsConfig(data, true),
			PlanOnly: true,
		},
	})
}

func TestAccGenericResource_ignoreBodyChangesTagsFromVariable(t *testing.T) {
	data := acceptance.BuildTestData(t, "azapi_resource", "test")
	r := GenericResource{}
	var resourceID string
	var resourceType string

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: ignoreBodyChangesTagsFromVariableConfig(data),
			ConfigVariables: config.Variables{
				"ignore_tags": config.BoolVariable(false),
			},
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				captureResourceState(data.ResourceName, &resourceID, &resourceType),
			),
		},
		{
			PreConfig: func() {
				if err := addTagOutsideTerraform(resourceID, resourceType); err != nil {
					t.Fatalf("adding tag outside Terraform: %+v", err)
				}
			},
			Config: ignoreBodyChangesTagsFromVariableConfig(data),
			ConfigVariables: config.Variables{
				"ignore_tags": config.BoolVariable(false),
			},
			PlanOnly:           true,
			ExpectNonEmptyPlan: true,
		},
		{
			Config: ignoreBodyChangesTagsFromVariableConfig(data),
			ConfigVariables: config.Variables{
				"ignore_tags": config.BoolVariable(true),
			},
			Check: func(*terraform.State) error {
				return checkTagInAzure(resourceID, resourceType, "out-of-band", "true")
			},
		},
		{
			Config: ignoreBodyChangesTagsFromVariableConfig(data),
			ConfigVariables: config.Variables{
				"ignore_tags": config.BoolVariable(true),
			},
			PlanOnly: true,
		},
	})
}

func ignoreBodyChangesTagsConfig(data acceptance.TestData, ignoreTags bool) string {
	ignoreBodyChanges := ""
	if ignoreTags {
		ignoreBodyChanges = `ignore_body_changes = ["tags"]`
	}

	config := fmt.Sprintf(`
%s

resource "azapi_resource" "test" {
  type      = "Microsoft.Network/virtualNetworks@2024-05-01"
  name      = "acctest-vnet-%s"
  parent_id = azapi_resource.resourceGroup.id
  location  = azapi_resource.resourceGroup.location

  body = {
    tags = {
      "managed-by" = "terraform"
    }
    properties = {
      addressSpace = {
        addressPrefixes = ["10.0.0.0/16"]
      }
    }
  }
  # ignore_body_changes_placeholder
}
`, (GenericResource{}).template(data), data.RandomString)
	return strings.Replace(config, "# ignore_body_changes_placeholder", ignoreBodyChanges, 1)
}

func ignoreBodyChangesTagsFromVariableConfig(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

variable "ignore_tags" {
  type = bool
}

resource "azapi_resource" "test" {
  type      = "Microsoft.Network/virtualNetworks@2024-05-01"
  name      = "acctest-vnet-%s"
  parent_id = azapi_resource.resourceGroup.id
  location  = azapi_resource.resourceGroup.location

  ignore_body_changes = var.ignore_tags ? ["tags"] : []

  body = {
    tags = {
      "managed-by" = "terraform"
    }
    properties = {
      addressSpace = {
        addressPrefixes = ["10.0.0.0/16"]
      }
    }
  }
}
`, (GenericResource{}).template(data), data.RandomString)
}

func addTagOutsideTerraform(resourceID, resourceType string) error {
	id, err := parse.ResourceIDWithResourceType(resourceID, resourceType)
	if err != nil {
		return fmt.Errorf("parsing resource ID: %+v", err)
	}

	client, err := acceptance.BuildTestClient()
	if err != nil {
		return fmt.Errorf("building test client: %+v", err)
	}

	ctx := context.Background()
	existing, err := client.ResourceClient.Get(ctx, id.AzureResourceId, id.ApiVersion, clients.DefaultRequestOptions())
	if err != nil {
		return fmt.Errorf("retrieving virtual network: %+v", err)
	}

	body, ok := existing.(map[string]interface{})
	if !ok {
		return fmt.Errorf("unexpected virtual network response type %T", existing)
	}
	tags, _ := body["tags"].(map[string]interface{})
	if tags == nil {
		tags = make(map[string]interface{})
	}
	tags["out-of-band"] = "true"
	body["tags"] = tags

	if id.ResourceDef != nil {
		writableBody := (*id.ResourceDef).GetWriteOnly(utils.NormalizeObject(body))
		body, ok = writableBody.(map[string]interface{})
		if !ok {
			return fmt.Errorf("unexpected writable virtual network body type %T", writableBody)
		}
	}

	if _, err := client.ResourceClient.CreateOrUpdate(ctx, id.AzureResourceId, id.ApiVersion, body, clients.DefaultRequestOptions()); err != nil {
		return fmt.Errorf("updating virtual network tags: %+v", err)
	}
	return nil
}

func checkTagInAzure(resourceID, resourceType, key, expectedValue string) error {
	id, err := parse.ResourceIDWithResourceType(resourceID, resourceType)
	if err != nil {
		return fmt.Errorf("parsing resource ID: %+v", err)
	}

	client, err := acceptance.BuildTestClient()
	if err != nil {
		return fmt.Errorf("building test client: %+v", err)
	}

	body, err := client.ResourceClient.Get(context.Background(), id.AzureResourceId, id.ApiVersion, clients.DefaultRequestOptions())
	if err != nil {
		return fmt.Errorf("retrieving virtual network: %+v", err)
	}
	bodyMap, ok := body.(map[string]interface{})
	if !ok {
		return fmt.Errorf("unexpected virtual network response type %T", body)
	}
	tags, ok := bodyMap["tags"].(map[string]interface{})
	if !ok {
		return fmt.Errorf("virtual network response has unexpected tags type %T", bodyMap["tags"])
	}
	if actualValue := tags[key]; actualValue != expectedValue {
		return fmt.Errorf("expected Azure tag %q to be %q, got %q", key, expectedValue, actualValue)
	}
	return nil
}
