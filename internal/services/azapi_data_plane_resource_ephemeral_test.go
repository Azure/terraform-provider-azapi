package services_test

import (
	"fmt"
	"testing"

	"github.com/Azure/terraform-provider-azapi/internal/acceptance"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

type DataPlaneResourceEphemeral struct{}

func TestAccDataPlaneResourceEphemeral_keyVaultSecret(t *testing.T) {
	data := acceptance.BuildTestData(t, "ephemeral.azapi_data_plane_resource", "test")
	r := DataPlaneResourceEphemeral{}

	data.DataSourceTest(t, []resource.TestStep{
		{
			Config: r.keyVaultSecret(data),
			Check:  resource.ComposeTestCheckFunc(),
		},
	})
}

func TestAccDataPlaneResourceEphemeral_keyVaultCertificate(t *testing.T) {
	data := acceptance.BuildTestData(t, "ephemeral.azapi_data_plane_resource", "test")
	r := DataPlaneResourceEphemeral{}

	data.DataSourceTest(t, []resource.TestStep{
		{
			Config: r.keyVaultCertificate(data),
			Check:  resource.ComposeTestCheckFunc(),
		},
	})
}

func TestAccDataPlaneResourceEphemeral_timeouts(t *testing.T) {
	data := acceptance.BuildTestData(t, "ephemeral.azapi_data_plane_resource", "test")
	r := DataPlaneResourceEphemeral{}

	data.DataSourceTest(t, []resource.TestStep{
		{
			Config: r.timeouts(data),
			Check:  resource.ComposeTestCheckFunc(),
		},
	})
}

func (r DataPlaneResourceEphemeral) keyVaultSecret(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

ephemeral "azapi_data_plane_resource" "test" {
  name                   = azapi_data_plane_resource.test.name
  parent_id              = azapi_data_plane_resource.test.parent_id
  type                   = azapi_data_plane_resource.test.type
  response_export_values = ["*"]
}
`, DataPlaneResource{}.keyVaultSecret(data))
}

func (r DataPlaneResourceEphemeral) keyVaultCertificate(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

ephemeral "azapi_data_plane_resource" "test" {
  name                   = azapi_data_plane_resource.test.name
  parent_id              = azapi_data_plane_resource.test.parent_id
  type                   = azapi_data_plane_resource.test.type
  response_export_values = ["*"]
}
`, DataPlaneResource{}.keyVaultCertificate(data))
}

func (r DataPlaneResourceEphemeral) timeouts(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

ephemeral "azapi_data_plane_resource" "test" {
  name                   = azapi_data_plane_resource.test.name
  parent_id              = azapi_data_plane_resource.test.parent_id
  type                   = azapi_data_plane_resource.test.type
  response_export_values = ["*"]

  timeouts {
    open = "10m"
  }
}
`, DataPlaneResource{}.keyVaultSecret(data))
}
