package acceptance

import (
	"fmt"
	"os"
	"testing"

	"github.com/Azure/terraform-provider-azapi/internal/clients"
	"github.com/Azure/terraform-provider-azapi/internal/provider"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

const (
	// charSetAlphaNum is the alphanumeric character set for use with randStringFromCharSet
	charSetAlphaNum = "abcdefghijklmnopqrstuvwxyz012346789"
)

type TestData struct {
	// LocationPrimary is the Primary Azure Region which should be used for testing
	LocationPrimary string

	// LocationSecondary is the Secondary Azure Region which should be used for testing
	LocationSecondary string

	// LocationTernary is the Ternary Azure Region which should be used for testing
	LocationTernary string

	// RandomInteger is a random integer which is unique to this test case
	RandomInteger int

	// RandomString is a random 5 character string is unique to this test case
	RandomString string

	// ResourceName is the fully qualified resource name, comprising of the
	// resource type and then the resource label
	// e.g. `azurerm_resource_group.test`
	ResourceName string

	// ResourceType is the Terraform Resource Type - `azurerm_resource_group`
	ResourceType string

	// resourceLabel is the local used for the resource - generally "test""
	resourceLabel string
}

// BuildTestData generates some test data for the given resource
func BuildTestData(t *testing.T, resourceType string, resourceLabel string) TestData {
	return TestData{
		RandomInteger: RandTimeInt(),
		RandomString:  acctest.RandStringFromCharSet(5, charSetAlphaNum),
		ResourceName:  fmt.Sprintf("%s.%s", resourceType, resourceLabel),

		ResourceType:      resourceType,
		resourceLabel:     resourceLabel,
		LocationPrimary:   os.Getenv("ARM_TEST_LOCATION"),
		LocationSecondary: os.Getenv("ARM_TEST_LOCATION_ALT"),
		LocationTernary:   os.Getenv("ARM_TEST_LOCATION_ALT2"),
	}
}

// RandomIntOfLength is a random 8 to 18 digit integer which is unique to this test case
func (td *TestData) RandomInt() int {
	return RandTimeInt()
}

func (td *TestData) RandomStringOfLength(len int) string {
	return acctest.RandStringFromCharSet(len, charSetAlphaNum)
}

// lintignore:AT001
func (td TestData) DataSourceTest(t *testing.T, steps []resource.TestStep) {
	// DataSources don't need a check destroy - however since this is a wrapper function
	// and not matching the ignore pattern `XXX_data_source_test.go`, this needs to be explicitly opted out
	testCase := resource.TestCase{
		PreCheck: func() { PreCheck(t) },
		Steps:    steps,
	}
	td.runAcceptanceTest(t, testCase)
}

func (td TestData) ResourceTest(t *testing.T, testResource TestResource, steps []resource.TestStep) {
	testCase := resource.TestCase{
		PreCheck: func() { PreCheck(t) },
		CheckDestroy: func(s *terraform.State) error {
			client, err := BuildTestClient()
			if err != nil {
				return fmt.Errorf("building client: %+v", err)
			}
			return CheckDestroyedFunc(client, testResource, td.ResourceType, td.ResourceName)(s)
		},
		Steps: steps,
	}
	td.runAcceptanceTest(t, testCase)
}

func (td TestData) runAcceptanceTest(t *testing.T, testCase resource.TestCase) {
	testCase.ExternalProviders = td.externalProviders()
	testCase.ProviderFactories = td.providers()

	resource.ParallelTest(t, testCase)
}

func (td TestData) providers() map[string]func() (*schema.Provider, error) {
	return map[string]func() (*schema.Provider, error){
		"azapi": func() (*schema.Provider, error) { //nolint:unparam
			azapi := provider.AzureProvider()
			return azapi, nil
		},
	}
}

func (td TestData) externalProviders() map[string]resource.ExternalProvider {
	return map[string]resource.ExternalProvider{
		"azurerm": {
			Source:            "registry.terraform.io/hashicorp/azurerm",
			VersionConstraint: "= 3.41.0",
		},
	}
}

func PreCheck(t *testing.T) {
	variables := []string{
		"ARM_CLIENT_ID",
		"ARM_CLIENT_SECRET",
		"ARM_SUBSCRIPTION_ID",
		"ARM_TENANT_ID",
		"ARM_TEST_LOCATION",
		"ARM_TEST_LOCATION_ALT",
		"ARM_TEST_LOCATION_ALT2",
	}

	for _, variable := range variables {
		value := os.Getenv(variable)
		if value == "" {
			t.Fatalf("`%s` must be set for acceptance tests!", variable)
		}
	}
}

// CheckDestroyedFunc returns a TestCheckFunc which validates the resource no longer exists
func CheckDestroyedFunc(client *clients.Client, testResource TestResource, resourceType, resourceName string) func(state *terraform.State) error {
	return func(state *terraform.State) error {
		ctx := client.StopContext

		for label, resourceState := range state.RootModule().Resources {
			if resourceState.Type != resourceType {
				continue
			}
			if label != resourceName {
				continue
			}

			// Destroy is unconcerned with an error checking the status, since this is going to be "not found"
			result, err := testResource.Exists(ctx, client, resourceState.Primary)
			if result == nil && err == nil {
				return fmt.Errorf("should have either an error or a result when checking if %q has been destroyed", resourceName)
			}
			if result != nil && *result {
				return fmt.Errorf("%q still exists", resourceName)
			}
		}

		return nil
	}
}
