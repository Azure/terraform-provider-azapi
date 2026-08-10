package services_test

import (
	"fmt"
	"testing"

	"github.com/Azure/terraform-provider-azapi/internal/acceptance"
	"github.com/Azure/terraform-provider-azapi/internal/acceptance/check"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

type DataPlaneResourceDataSource struct{}

func TestAccDataPlaneResourceDataSource_keyVaultSecret(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azapi_data_plane_resource", "test")
	d := DataPlaneResourceDataSource{}
	data.DataSourceTest(t, []resource.TestStep{
		{
			Config: d.keyVaultSecret(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("name").HasValue(fmt.Sprintf("acctest%s", data.RandomString)),
				check.That(data.ResourceName).Key("type").HasValue("Microsoft.KeyVault/vaults/secrets@7.4"),
				check.That(data.ResourceName).Key("output.value").HasValue("secret-value"),
			),
		},
	})
}

func TestAccDataPlaneResourceDataSource_keyVaultCertificate(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azapi_data_plane_resource", "test")
	d := DataPlaneResourceDataSource{}
	data.DataSourceTest(t, []resource.TestStep{
		{
			Config: d.keyVaultCertificate(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("name").HasValue(fmt.Sprintf("acctest%s", data.RandomString)),
				check.That(data.ResourceName).Key("type").HasValue("Microsoft.KeyVault/vaults/certificates@7.4"),
				check.That(data.ResourceName).Key("output.policy.issuer.name").HasValue("Self"),
			),
		},
	})
}

func TestAccDataPlaneResourceDataSource_keyVaultCertificateContact(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azapi_data_plane_resource", "test")
	d := DataPlaneResourceDataSource{}
	data.DataSourceTest(t, []resource.TestStep{
		{
			Config: d.keyVaultCertificateContact(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("type").HasValue("Microsoft.KeyVault/vaults/certificates/contacts@7.4"),
				check.That(data.ResourceName).Key("output.contacts.0.email").HasValue("foo@contoso.com"),
			),
		},
	})
}

func TestAccDataPlaneResourceDataSource_timeouts(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azapi_data_plane_resource", "test")
	d := DataPlaneResourceDataSource{}

	data.DataSourceTest(t, []resource.TestStep{
		{
			Config:            d.timeouts(data),
			ExternalProviders: externalProvidersAzurerm(),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("name").HasValue("mykey"),
				check.That(data.ResourceName).Key("output.value").HasValue("myvalue"),
			),
		},
	})
}

func (d DataPlaneResourceDataSource) keyVaultSecret(data acceptance.TestData) string {
	return fmt.Sprintf(`
	%s

	data "azapi_data_plane_resource" "test"{
		name                   = azapi_data_plane_resource.test.name
  		parent_id              = azapi_data_plane_resource.test.parent_id
  		type                   = azapi_data_plane_resource.test.type
  		response_export_values = ["*"]
	}
	`, DataPlaneResource{}.keyVaultSecret(data))
}

func (d DataPlaneResourceDataSource) keyVaultCertificate(data acceptance.TestData) string {
	return fmt.Sprintf(`
	%s

	data "azapi_data_plane_resource" "test"{
		name                   = azapi_data_plane_resource.test.name
  		parent_id              = azapi_data_plane_resource.test.parent_id
  		type                   = azapi_data_plane_resource.test.type
  		response_export_values = ["*"]
	}
	`, DataPlaneResource{}.keyVaultCertificate(data))
}

func (d DataPlaneResourceDataSource) keyVaultCertificateContact(data acceptance.TestData) string {
	return fmt.Sprintf(`
	%s

	data "azapi_data_plane_resource" "test"{
  		parent_id              = azapi_data_plane_resource.contact.parent_id
  		type                   = azapi_data_plane_resource.contact.type
  		response_export_values = ["*"]
	}
	`, DataPlaneResource{}.keyVaultCertificateContact(data))
}

func (d DataPlaneResourceDataSource) timeouts(data acceptance.TestData) string {
	return fmt.Sprintf(`
	%s
	
	data "azapi_data_plane_resource" "test"{
		name                   = azapi_data_plane_resource.test.name
  		parent_id              = azapi_data_plane_resource.test.parent_id
  		type                   = azapi_data_plane_resource.test.type
  		response_export_values = {
  			value = "value"
  		}
		timeouts {
			read = "10m"
		}
	}
	`, DataPlaneResource{}.timeouts(data))
}
