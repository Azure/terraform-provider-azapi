package services_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/Azure/terraform-provider-azapi/internal/acceptance"
	"github.com/Azure/terraform-provider-azapi/internal/azure"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/querycheck"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

type AzapiResourceListResource struct{}

func TestAccAzapiResourceList_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azapi_resource", "test")
	r := AzapiResourceListResource{}

	resource.Test(t, resource.TestCase{
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_14_0),
		},
		ProtoV6ProviderFactories: data.Providers(),
		Steps: []resource.TestStep{
			{
				Config: r.basicList(data),
			}, {
				Query:  true,
				Config: r.basicQueryByResourceGroup(data),
				QueryResultChecks: []querycheck.QueryResultCheck{
					querycheck.ExpectLength("azapi_resource.test", 1),
					querycheck.ExpectIdentity("azapi_resource.test", map[string]knownvalue.Check{
						"type": knownvalue.StringExact("Microsoft.DataFactory/factories@2018-06-01"),
						"id":   knownvalue.StringExact(fmt.Sprintf(`/subscriptions/%s/resourceGroups/acctestrg-%d/providers/Microsoft.DataFactory/factories/acctestdf%d`, os.Getenv("ARM_SUBSCRIPTION_ID"), data.RandomInteger, data.RandomInteger)),
					}),
					querycheck.ContainsResourceWithName("azapi_resource.test", fmt.Sprintf("Microsoft.DataFactory/factories@2018-06-01 - acctestdf%d", data.RandomInteger)),
				},
			},
		},
	})
}

func (r AzapiResourceListResource) basicList(data acceptance.TestData) string {
	return fmt.Sprintf(`
resource "azapi_resource" "resourceGroup" {
  type     = "Microsoft.Resources/resourceGroups@2023-07-01"
  name     = "acctestrg-%[1]d"
  location = "%[2]s"
}

resource "azapi_resource" "factory" {
  type      = "Microsoft.DataFactory/factories@2018-06-01"
  parent_id = azapi_resource.resourceGroup.id
  name      = "acctestdf%[1]d"
  location  = azapi_resource.resourceGroup.location
  body = {
    properties = {
      publicNetworkAccess = "Enabled"
    }
  }
}
`, data.RandomInteger, data.LocationPrimary)
}

func (r AzapiResourceListResource) basicQueryByResourceGroup(data acceptance.TestData) string {
	subscriptionId := os.Getenv("ARM_SUBSCRIPTION_ID")
	return fmt.Sprintf(`
list "azapi_resource" "test" {
  provider = azapi
  config {
    type      = "Microsoft.DataFactory/factories@2018-06-01"
    parent_id = "/subscriptions/%[2]s/resourceGroups/acctestrg-%[1]d"
  }
}
`, data.RandomInteger, subscriptionId)
}

func TestAccAzapiResourceList_withoutType(t *testing.T) {
	data := acceptance.BuildTestData(t, "azapi_resource", "test")
	r := AzapiResourceListResource{}

	storageApiVersions := azure.GetApiVersions("Microsoft.Storage/storageAccounts")
	if len(storageApiVersions) == 0 {
		t.Skip("No API versions found for Microsoft.Storage/storageAccounts")
	}
	storageApiVersion := storageApiVersions[len(storageApiVersions)-1]
	storageResourceType := fmt.Sprintf("Microsoft.Storage/storageAccounts@%s", storageApiVersion)

	resource.Test(t, resource.TestCase{
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_14_0),
		},
		ProtoV6ProviderFactories: data.Providers(),
		Steps: []resource.TestStep{
			{
				Config: r.listWithoutTypeSetup(data),
			}, {
				Query:  true,
				Config: r.queryWithoutType(data),
				QueryResultChecks: []querycheck.QueryResultCheck{
					querycheck.ExpectLength("azapi_resource.test", 2),
					querycheck.ContainsResourceWithName("azapi_resource.test", fmt.Sprintf("%s - acctst%d", storageResourceType, data.RandomInteger)),
					querycheck.ContainsResourceWithName("azapi_resource.test", fmt.Sprintf("Microsoft.DataFactory/factories@2018-06-01 - acctestdf%d", data.RandomInteger)),
					querycheck.ExpectIdentity("azapi_resource.test", map[string]knownvalue.Check{
						"type": knownvalue.StringExact(storageResourceType),
						"id":   knownvalue.StringExact(fmt.Sprintf(`/subscriptions/%s/resourceGroups/acctestrg-%d/providers/Microsoft.Storage/storageAccounts/acctst%d`, os.Getenv("ARM_SUBSCRIPTION_ID"), data.RandomInteger, data.RandomInteger)),
					}),
					querycheck.ExpectIdentity("azapi_resource.test", map[string]knownvalue.Check{
						"type": knownvalue.StringExact("Microsoft.DataFactory/factories@2018-06-01"),
						"id":   knownvalue.StringExact(fmt.Sprintf(`/subscriptions/%s/resourceGroups/acctestrg-%d/providers/Microsoft.DataFactory/factories/acctestdf%d`, os.Getenv("ARM_SUBSCRIPTION_ID"), data.RandomInteger, data.RandomInteger)),
					}),
				},
			},
		},
	})
}

func (r AzapiResourceListResource) listWithoutTypeSetup(data acceptance.TestData) string {
	return fmt.Sprintf(`
resource "azapi_resource" "resourceGroup" {
  type     = "Microsoft.Resources/resourceGroups@2023-07-01"
  name     = "acctestrg-%[1]d"
  location = "%[2]s"
}

resource "azapi_resource" "storage" {
  type      = "Microsoft.Storage/storageAccounts@2023-05-01"
  parent_id = azapi_resource.resourceGroup.id
  name      = "acctst%[1]d"
  location  = azapi_resource.resourceGroup.location
  body = {
    kind = "StorageV2"
    sku = {
      name = "Standard_LRS"
    }
    properties = {
      accessTier = "Hot"
    }
  }
}

resource "azapi_resource" "factory" {
  type      = "Microsoft.DataFactory/factories@2018-06-01"
  parent_id = azapi_resource.resourceGroup.id
  name      = "acctestdf%[1]d"
  location  = azapi_resource.resourceGroup.location
  body = {
    properties = {
      publicNetworkAccess = "Enabled"
    }
  }
}
`, data.RandomInteger, data.LocationPrimary)
}

func (r AzapiResourceListResource) queryWithoutType(data acceptance.TestData) string {
	subscriptionId := os.Getenv("ARM_SUBSCRIPTION_ID")
	return fmt.Sprintf(`
list "azapi_resource" "test" {
  provider = azapi
  config {
    parent_id = "/subscriptions/%[2]s/resourceGroups/acctestrg-%[1]d"
  }
}
`, data.RandomInteger, subscriptionId)
}
