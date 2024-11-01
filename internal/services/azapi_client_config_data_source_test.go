package services_test

import (
	"os"
	"regexp"
	"testing"

	"github.com/Azure/terraform-provider-azapi/internal/acceptance"
	"github.com/Azure/terraform-provider-azapi/internal/acceptance/check"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

type ClientConfigDataSource struct{}

var idRegex *regexp.Regexp = regexp.MustCompile("^[A-Fa-f0-9]{8}-[A-Fa-f0-9]{4}-[A-Fa-f0-9]{4}-[A-Fa-f0-9]{4}-[A-Fa-f0-9]{12}$")

func TestAccClientConfigDataSource_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azapi_client_config", "test")
	r := ClientConfigDataSource{}

	tenantId := os.Getenv("ARM_TENANT_ID")
	subscriptionId := os.Getenv("ARM_SUBSCRIPTION_ID")

	data.DataSourceTest(t, []resource.TestStep{
		{
			Config: r.basic(),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("tenant_id").HasValue(tenantId),
				check.That(data.ResourceName).Key("subscription_id").HasValue(subscriptionId),
				check.That(data.ResourceName).Key("object_id").MatchesRegex(idRegex),
			),
		},
	})
}

func TestAccClientConfigDataSource_azcli(t *testing.T) {
	if ok := os.Getenv("ARM_USE_CLI"); ok == "" {
		t.Skip("Skipping as `ARM_USE_CLI` is not specified")
	}

	data := acceptance.BuildTestData(t, "data.azapi_client_config", "test")
	r := ClientConfigDataSource{}

	data.DataSourceTest(t, []resource.TestStep{
		{
			Config: r.basic(),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("tenant_id").MatchesRegex(idRegex),
				check.That(data.ResourceName).Key("subscription_id").MatchesRegex(idRegex),
				check.That(data.ResourceName).Key("object_id").MatchesRegex(idRegex),
			),
		},
	})
}

func (r ClientConfigDataSource) basic() string {
	return `
data "azapi_client_config" "test" {}
`
}
