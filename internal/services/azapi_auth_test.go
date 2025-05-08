package services_test

import (
	"os"
	"testing"

	"github.com/Azure/terraform-provider-azapi/internal/acceptance"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccAuth_clientCertificatePath(t *testing.T) {
	if ok := os.Getenv("ARM_CLIENT_CERTIFICATE_PATH"); ok == "" {
		t.Skip("Skipping as `ARM_CLIENT_CERTIFICATE_PATH` is not specified")
	}

	data := acceptance.BuildTestData(t, "data.azapi_resource_action", "test")
	r := ActionDataSource{}

	data.DataSourceTest(t, []resource.TestStep{
		{
			Config: r.providerAction(),
			Check:  resource.ComposeTestCheckFunc(),
		},
	})
}

func TestAccAuth_clientCertificate(t *testing.T) {
	if ok := os.Getenv("ARM_CLIENT_CERTIFICATE"); ok == "" {
		t.Skip("Skipping as `ARM_CLIENT_CERTIFICATE` is not specified")
	}

	data := acceptance.BuildTestData(t, "data.azapi_resource_action", "test")
	r := ActionDataSource{}

	data.DataSourceTest(t, []resource.TestStep{
		{
			Config: r.providerAction(),
			Check:  resource.ComposeTestCheckFunc(),
		},
	})
}

func TestAccAuth_oidc(t *testing.T) {
	if ok := os.Getenv("ARM_USE_OIDC"); ok == "" {
		t.Skip("Skipping as `ARM_USE_OIDC` is not specified")
	}

	data := acceptance.BuildTestData(t, "data.azapi_resource_action", "test")
	r := ActionDataSource{}

	data.DataSourceTest(t, []resource.TestStep{
		{
			Config: r.providerAction(),
			Check:  resource.ComposeTestCheckFunc(),
		},
	})
}

func TestAccAuth_azcli(t *testing.T) {
	if ok := os.Getenv("ARM_USE_CLI"); ok == "" {
		t.Skip("Skipping as `ARM_USE_CLI` is not specified")
	}

	data := acceptance.BuildTestData(t, "data.azapi_resource_action", "test")
	r := ActionDataSource{}

	data.DataSourceTest(t, []resource.TestStep{
		{
			Config: r.providerAction(),
			Check:  resource.ComposeTestCheckFunc(),
		},
	})
}
