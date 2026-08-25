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
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
)

func TestAccGenericResource_ignoreBodyChangesSubnets(t *testing.T) {
	data := acceptance.BuildTestData(t, "azapi_resource", "test")
	r := GenericResource{}
	var resourceID string
	var resourceType string

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: ignoreBodyChangesSubnetsConfig(data, true),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				captureResourceState(data.ResourceName, &resourceID, &resourceType),
			),
		},
		{
			PreConfig: func() {
				if err := addSubnetOutsideTerraform(resourceID, resourceType); err != nil {
					t.Fatalf("adding subnet outside Terraform: %+v", err)
				}
			},
			Config:   ignoreBodyChangesSubnetsConfig(data, true),
			PlanOnly: true,
		},
		{
			Config: ignoreBodyChangesSubnetsConfig(data, false),
			Check: func(*terraform.State) error {
				return checkSubnetsInAzure(resourceID, resourceType, []string{"managed"})
			},
		},
		{
			Config: ignoreBodyChangesSubnetsConfig(data, true),
		},
	})
}

func TestAccGenericResource_ignoreBodyChangesNonExistentPath(t *testing.T) {
	data := acceptance.BuildTestData(t, "azapi_resource", "test")
	r := GenericResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: ignoreBodyChangesNonExistentPathConfig(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
	})
}

func TestAccGenericResource_ignoreBodyChangesTagsDuringSubnetUpdate(t *testing.T) {
	data := acceptance.BuildTestData(t, "azapi_resource", "test")
	r := GenericResource{}
	var resourceID string
	var resourceType string

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: ignoreBodyChangesTagsDuringSubnetUpdateConfig(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				captureResourceState(data.ResourceName, &resourceID, &resourceType),
			),
		},
		{
			PreConfig: func() {
				if err := addSubnetOutsideTerraform(resourceID, resourceType); err != nil {
					t.Fatalf("adding subnet outside Terraform: %+v", err)
				}
				if err := addTagOutsideTerraform(resourceID, resourceType); err != nil {
					t.Fatalf("adding tag outside Terraform: %+v", err)
				}
			},
			Config: ignoreBodyChangesTagsDuringSubnetUpdateConfig(data),
			ConfigPlanChecks: resource.ConfigPlanChecks{
				PreApply: []plancheck.PlanCheck{
					plancheck.ExpectResourceAction(data.ResourceName, plancheck.ResourceActionUpdate),
					plancheck.ExpectKnownValue(
						data.ResourceName,
						tfjsonpath.New("body").AtMapKey("tags"),
						knownvalue.MapExact(map[string]knownvalue.Check{
							"managed-by": knownvalue.StringExact("terraform"),
						}),
					),
					plancheck.ExpectKnownValue(
						data.ResourceName,
						tfjsonpath.New("body").AtMapKey("properties").AtMapKey("subnets"),
						knownvalue.ListSizeExact(1),
					),
				},
			},
			Check: func(*terraform.State) error {
				return checkTagInAzure(resourceID, resourceType, "out-of-band", "true")
			},
		},
	})
}

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

func ignoreBodyChangesSubnetsConfig(data acceptance.TestData, ignoreSubnets bool) string {
	ignoreBodyChanges := ""
	if ignoreSubnets {
		ignoreBodyChanges = `ignore_body_changes = ["properties.subnets"]`
	}

	config := fmt.Sprintf(`
%s

resource "azapi_resource" "test" {
  type      = "Microsoft.Network/virtualNetworks@2024-05-01"
  name      = "acctest-vnet-%s"
  parent_id = azapi_resource.resourceGroup.id
  location  = azapi_resource.resourceGroup.location

  body = {
    properties = {
      addressSpace = {
        addressPrefixes = ["10.0.0.0/16"]
      }
      subnets = [
        {
          name = "managed"
          properties = {
            addressPrefix = "10.0.1.0/24"
          }
        }
      ]
    }
  }
  # ignore_body_changes_placeholder
}
`, (GenericResource{}).template(data), data.RandomString)
	return strings.Replace(config, "# ignore_body_changes_placeholder", ignoreBodyChanges, 1)
}

func ignoreBodyChangesNonExistentPathConfig(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azapi_resource" "test" {
  type      = "Microsoft.Network/virtualNetworks@2024-05-01"
  name      = "acctest-vnet-%s"
  parent_id = azapi_resource.resourceGroup.id
  location  = azapi_resource.resourceGroup.location

  ignore_body_changes = ["does_not_exist"]

  body = {
    properties = {
      addressSpace = {
        addressPrefixes = ["10.0.0.0/16"]
      }
    }
  }
}
`, (GenericResource{}).template(data), data.RandomString)
}

func ignoreBodyChangesTagsDuringSubnetUpdateConfig(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azapi_resource" "test" {
  type      = "Microsoft.Network/virtualNetworks@2024-05-01"
  name      = "acctest-vnet-%s"
  parent_id = azapi_resource.resourceGroup.id
  location  = azapi_resource.resourceGroup.location

  ignore_body_changes = ["tags"]

  body = {
    tags = {
      "managed-by" = "terraform"
    }
    properties = {
      addressSpace = {
        addressPrefixes = ["10.0.0.0/16"]
      }
      subnets = [
        {
          name = "managed"
          properties = {
            addressPrefix = "10.0.1.0/24"
          }
        },
      ]
    }
  }
}
`, (GenericResource{}).template(data), data.RandomString)
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

func addSubnetOutsideTerraform(resourceID, resourceType string) error {
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
	properties, ok := body["properties"].(map[string]interface{})
	if !ok {
		return fmt.Errorf("virtual network response has unexpected properties type %T", body["properties"])
	}
	subnets, ok := properties["subnets"].([]interface{})
	if !ok {
		return fmt.Errorf("virtual network response has unexpected subnets type %T", properties["subnets"])
	}
	properties["subnets"] = append(subnets, map[string]interface{}{
		"name": "out-of-band",
		"properties": map[string]interface{}{
			"addressPrefix": "10.0.2.0/24",
		},
	})

	if id.ResourceDef != nil {
		writableBody := (*id.ResourceDef).GetWriteOnly(utils.NormalizeObject(body))
		body, ok = writableBody.(map[string]interface{})
		if !ok {
			return fmt.Errorf("unexpected writable virtual network body type %T", writableBody)
		}
	}

	if _, err := client.ResourceClient.CreateOrUpdate(ctx, id.AzureResourceId, id.ApiVersion, body, clients.DefaultRequestOptions()); err != nil {
		return fmt.Errorf("updating virtual network subnets: %+v", err)
	}
	return nil
}

func checkSubnetsInAzure(resourceID, resourceType string, expectedNames []string) error {
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
	properties, ok := bodyMap["properties"].(map[string]interface{})
	if !ok {
		return fmt.Errorf("virtual network response has unexpected properties type %T", bodyMap["properties"])
	}
	subnets, ok := properties["subnets"].([]interface{})
	if !ok {
		return fmt.Errorf("virtual network response has unexpected subnets type %T", properties["subnets"])
	}
	if len(subnets) != len(expectedNames) {
		return fmt.Errorf("expected %d Azure subnets, got %d", len(expectedNames), len(subnets))
	}
	for i, expectedName := range expectedNames {
		subnet, ok := subnets[i].(map[string]interface{})
		if !ok {
			return fmt.Errorf("virtual network subnet %d has unexpected type %T", i, subnets[i])
		}
		if actualName := subnet["name"]; actualName != expectedName {
			return fmt.Errorf("expected Azure subnet %d to be named %q, got %q", i, expectedName, actualName)
		}
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
